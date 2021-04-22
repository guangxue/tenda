import {WhenClick } from './helper.js';

function getCurrentDateTime(numonly) {
	let date = new Date();
	let dates = date.toString().split(" ");
	let currTime = dates[4];
	var currentDateTime = ""
	var month = date.getMonth()+1;
	if(month < 10) {
		month = `0${month}`
	}
	if (numonly) {
		currentDateTime = `${date.getFullYear()}${month}${date.getDate()}`;
	}
	else {
		currentDateTime = `${date.getFullYear()}-${month}-${date.getDate()} ${currTime}`;
	}
	console.log(currentDateTime)
	return currentDateTime;
}

let today = 'Today: '+ getCurrentDateTime();
document.querySelector('#today').innerHTML = today;

let modelinput = document.querySelector('input[name="model"]')
if(modelinput) {
	fetch("https://gzhang.dev/tenda/api/models?model=all")
	.then(response => {
		return response.json()
	})
	.then(data=> {
		// console.log("data->",data);
		var modelist = [];
		for(let [key, value] of Object.entries(data)) {
			modelist.push(value);
		}
		return modelist;
	})
	.then(models => {
		let datalistElem = document.createElement('datalist')
		datalistElem.id='modelist'

		let formElem = document.querySelector('form');
		formElem.insertAdjacentElement('afterend', datalistElem);
		let optionFragement = new DocumentFragment();
		models.forEach( model => {
			let currentOpt = '<option value="'+model+'">';
			let opt = document.createElement('option');
			opt.value = model;
			optionFragement.appendChild(opt);
		});
		datalistElem.appendChild(optionFragement)
	})
	.catch( err => {
		console.log("modelallerr:", err);
	})
}




let queryButton = document.querySelector('#querybutton');
if(queryButton) {
	queryButton.addEventListener('click', function(e) {
		e.preventDefault();
		let model = document.querySelector('#mod').value;
		let location = document.querySelector('#loc').value;
		if(model) {
			let url = "https://gzhang.dev/tenda/api/models?model="+model;
			fetch(url).then(response => { return response.json()})
			.then(data => {
				let sum_total = 0;
				let table = "<table><tr><th>Model</th><th>Location</th><th>Unit</th><th>Cartons</th><th>Loose</th><th>Total</th></tr>";
				data.forEach(m => {
					// console.log("model ->", m);
					let row = `<tr><td>${m.model}</td><td>${m.location}</td><td>${m.unit}</td><td>${m.cartons}</td><td>${m.loose}</td><td>${m.total}</td></tr>`
					table += row
					sum_total += m.total;
				});
				table += `<tr><td></td><td></td><td></td><td></td><td></td><td>${sum_total}</td></tr></table>`
				document.querySelector("#mfb").innerHTML = table;
			})
			.catch( err => {
				console.log("model query FAILED:", err);
			});
		}
		if(location) {
			console.log("location->", location);
			let url = "https://gzhang.dev/tenda/api/models?location="+location;
			console.log("fetch url:", url);
			fetch(url).then(response => {
				return response.json()
			})
			.then(data => {
				console.log("location json:", data);
				let table = "<table><tr><th>Location</th><th>model</th><th>Unit</th><th>Cartons</th><th>Loose</th><th>Total</th></tr>";
				data.forEach( m=> {
					let row = `<tr><td>${m.location}</td><td>${m.model}</td><td>${m.unit}</td><td>${m.cartons}</td><td>${m.loose}</td><td>${m.total}</td></tr>`
					table += row
				})
				document.querySelector("#lfb").innerHTML = table;
			})
			.catch( err => {
				console.log("location query FAILED:", err);
			});
		}
	});
}

let addBtn = document.querySelector('.addBtn');
if(addBtn) {
	addBtn.addEventListener('click', function(e) {
		e.preventDefault();

		let currentDateTime = getCurrentDateTime()

		let model = document.querySelector('#mod').value;
		let qty = document.querySelector('input[name="qty"]').value;
		let customer = document.querySelector('input[name="customer"]').value;
		let tableElem = document.querySelector('table > tbody');
		if(tableElem) {
			let tplrow = document.querySelector('#htmpl_order');
			var row = tplrow.content.cloneNode(true);
			var td = row.querySelectorAll('td');
			td.item(0).textContent = model;
			td.item(1).textContent = qty;
			td.item(2).textContent = customer;
			tableElem.appendChild(row);
		}
		else {
			let table = "<table><tr><th>Model</th><th>Quantity</th><th>Cusotmer</th></tr>";
			var row = `<tr><td>${model}</td><td>${qty}</td><td>${customer}</td></tr>`;
			table = table + row + "</table>";
			document.querySelector("#po").innerHTML = table;
		}
		
	});
}


let datalistElem = document.createElement('datalist');
datalistElem.id='POAENumber';
let PONumber = `PO${getCurrentDateTime(true)}`;
let AENumber = `AE${getCurrentDateTime(true)}`;
var PO_option = document.createElement('option');
PO_option.textContent = PONumber;
var AE_option = document.createElement('option');
AE_option.textContent = AENumber;
datalistElem.appendChild(PO_option);
datalistElem.appendChild(AE_option);
let formElem = document.querySelector('form');
if(formElem) {
	formElem.insertAdjacentElement('afterend', datalistElem);
}

WhenClick('#pickbtn', function(e) {
	e.preventDefault();
	console.log("Pick clicked")
	let currentDateTime = getCurrentDateTime();
	let formEl = document.querySelector('#pickform')
	let formData = new FormData(formEl);
	let PNO = document.querySelector('input[name="PNO"').value;
	let model = document.querySelector('input[name="model"]').value;
	let qty = document.querySelector('input[name="qty"]').value;
	let customer = document.querySelector('input[name="customer"]').value;
	let data = new URLSearchParams(formData);

	for(const pair of formData) {
		data.append(pair[0], pair[1])
	}

	data.append('now', currentDateTime)
	fetch("https://gzhang.dev/tenda/api/picked", {
		method: "POST",
		body: data,
	})
	.then(resp => { 
		return resp.json()
	})
	.then(data => {
		console.log(data.lastId)
	})
});

let picktable = document.querySelector("#picktable");
if(picktable) {
	let todayDateTime = getCurrentDateTime();
	let today = todayDateTime.split(' ')[0]
	console.log("today::", today);
	let url = "https://gzhang.dev/tenda/api/picked?date=today";
	fetch(url).then(response => { return response.json()})
	.then(data => {
		console.log("data from /picked/today:", data);
		// let tplrow = document.querySelector('#htmpl_pack');
		// var row = tplrow.content.cloneNode(true);
		// var td = row.querySelectorAll('td');
		// td.item(0).textContent = model;
		// td.item(1).textContent = qty;
		// td.item(2).textContent = customer;
		// tableElem.appendChild(row);
		// data.forEach(m => {
		// 	// console.log("model ->", m);
		// 	let row = `<tr><td>${m.PID}</td><td>${m.PNO}</td><td>${m.model}</td><td>${m.qty}</td><td>${m.customer}</td><td>${m.updated}</td></tr>`
		// });
	})
	// .catch( err => {
	// 	console.log("today packed  FAILED:", err);
	// });
}
