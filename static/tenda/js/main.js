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
	let day = date.getDate();
	if (day < 10) {
		day = `0${day}`
	}
	if (numonly) {
		currentDateTime = `${date.getFullYear()}${month}${day}`;
	}
	else {
		currentDateTime = `${date.getFullYear()}-${month}-${day} ${currTime}`;
	}
	return currentDateTime;
}

// let today = 'Today: '+ getCurrentDateTime();
// document.querySelector('#today').innerHTML = today;

let modelinput = document.querySelector('input[name="model"]')
let modelNameinput = document.querySelector('input[name="modelName"]')
if(modelinput || modelNameinput) {
	fetch("https://gzhang.dev/tenda/api/models")
	.then(response => {
		return response.json()
	})
	.then(data=> {
		var modelist = [];
		for(let [key, value] of Object.entries(data)) {
			modelist.push(value);
		}
		return modelist;
	})
	.then(models => {
		let datalistElem = document.createElement('datalist')
		datalistElem.id='modelist'

		let container = document.querySelector('.container');
		container.insertAdjacentElement('afterend', datalistElem);
		let optionFragement = new DocumentFragment();
		models.forEach( m => {
			let currentOpt = '<option value="'+m.model+'">';
			let opt = document.createElement('option');
			opt.value = m.model;
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
				let table = "<table><thead><tr><th>Model</th><th>Location</th><th>Unit</th><th>Cartons</th><th>Boxes</th><th>Total</th></tr></thead><tbody>";
				data.forEach(m => {
					console.log("m.total ->", m.total);
					let row = `<tr><td>${m.model}</td><td>${m.location}</td><td>${m.unit}</td><td>${m.cartons}</td><td>${m.boxes}</td><td>${m.total}</td></tr>`
					table += row
					sum_total += parseFloat(m.total);
				});
				table += `<tr><td></td><td></td><td></td><td></td><td></td><td>${sum_total}</td></tr></tbody></table>`
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
					let row = `<tr><td>${m.location}</td><td>${m.model}</td><td>${m.unit}</td><td>${m.cartons}</td><td>${m.boxes}</td><td>${m.total}</td></tr>`
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
let RENumber = `RE${getCurrentDateTime(true)}`;

var PO_option = document.createElement('option');
PO_option.textContent = PONumber;
var AE_option = document.createElement('option');
AE_option.textContent = AENumber;
var RE_option = document.createElement('option');
RE_option.textContent = RENumber;

datalistElem.appendChild(PO_option);
datalistElem.appendChild(AE_option);
datalistElem.appendChild(RE_option);

let formElem = document.querySelector('form');
if(formElem) {
	formElem.insertAdjacentElement('afterend', datalistElem);
}



const modelInput = document.querySelector("input[name=modelName]");

if(modelInput) {
	modelInput.addEventListener('input', function(e) {
		let model = document.querySelector("input[name=modelName]").value;
		if(model) {
			console.log("model no:", model)
			fetch("https://gzhang.dev/tenda/api/locations?model="+model)
			.then( resp => {
				return resp.json();
			})
			.then( data => {
				console.log("locations data:", data);
				let optionFragement = new DocumentFragment();
				let selectButton = document.querySelector('#pick-location');
				selectButton.innerHTML = ""
				data.forEach( loc => {
					console.log("loc:", loc.location);
					// let currentOpt = '<option value="'+loc.location+'">'+loc.location+"</option>";
					let opt = document.createElement('option');
					opt.value = loc.location;
					opt.textContent = loc.location;
					optionFragement.appendChild(opt);
				});
				selectButton.appendChild(optionFragement)
			})
		}
	})
}
