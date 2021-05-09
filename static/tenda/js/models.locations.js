
// function getCurrentDateTime(numonly) {
// 	let date = new Date();
// 	let dates = date.toString().split(" ");
// 	let currTime = dates[4];
// 	var currentDateTime = ""
// 	var month = date.getMonth()+1;
// 	if(month < 10) {
// 		month = `0${month}`
// 	}
// 	let day = date.getDate();
// 	if (day < 10) {
// 		day = `0${day}`
// 	}
// 	if (numonly) {
// 		currentDateTime = `${date.getFullYear()}${month}${day}`;
// 	}
// 	else {
// 		currentDateTime = `${date.getFullYear()}-${month}-${day} ${currTime}`;
// 	}
// 	return currentDateTime;
// }

// let today = '<strong style="font-weight:bold;">Date:</strong> '+ getCurrentDateTime();
// let dateElm = document.querySelector('#today');
// if(dateElm) {
// 	dateElm.innerHTML = today;
// }

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
}






// let addBtn = document.querySelector('.addBtn');
// if(addBtn) {
// 	addBtn.addEventListener('click', function(e) {
// 		e.preventDefault();

// 		let currentDateTime = getCurrentDateTime()

// 		let model = document.querySelector('#mod').value;
// 		let qty = document.querySelector('input[name="qty"]').value;
// 		let customer = document.querySelector('input[name="customer"]').value;
// 		let tableElem = document.querySelector('table > tbody');
// 		if(tableElem) {
// 			let tplrow = document.querySelector('#htmpl_order');
// 			var row = tplrow.content.cloneNode(true);
// 			var td = row.querySelectorAll('td');
// 			td.item(0).textContent = model;
// 			td.item(1).textContent = qty;
// 			td.item(2).textContent = customer;
// 			tableElem.appendChild(row);
// 		}
// 		else {
// 			let table = "<table><tr><th>Model</th><th>Quantity</th><th>Cusotmer</th></tr>";
// 			var row = `<tr><td>${model}</td><td>${qty}</td><td>${customer}</td></tr>`;
// 			table = table + row + "</table>";
// 			document.querySelector("#po").innerHTML = table;
// 		}
		
// 	});
// }






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
