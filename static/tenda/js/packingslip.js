import {getCurrentDateTime, createTable} from "./helper.js"

/** INSERT INTO picked **/
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
const pickButton = document.querySelector("#pickbtn")
let currentLN = 1;

const inputPNO = document.querySelector("#panum")
const tbody = document.querySelector("#insertFB table tbody")

inputPNO.addEventListener("input", function() {
	if(inputPNO.value && !tbody) {
		let fetch_url = `https://gzhang.dev/tenda/api/picklist?PNO=${inputPNO.value}`;
		fetch(fetch_url)
		.then(resp => {
			return resp.json();
		})
		.then(data => {
			if(data[0]) {
				let titles = ["PID", "PNO", "model", "qty", "customer", "location", "status", "created_at", "Action"];
				data.forEach( d=> {
					d.Action = `<a href="/tenda/update/picklist?PID=${d.PID}">Modify</a>`;
				});
				let table = createTable(titles, data, titles);
				let insertFB = document.querySelector("#insertFB")
				insertFB.innerHTML = table.outerHTML;
				insertFB.style.display= "block"
			}
		})
	}
	else {
		insertFB.innerHTML = ""
	}
});

pickButton.addEventListener('click', function(e) {
	e.preventDefault();
	console.log("Pick clicked")

	let formEl = document.querySelector('#pickform')
	let formData = new FormData(formEl);
	let data = new URLSearchParams(formData);

	for(const pair of formData) {
		if(!pair[1]) {
			console.log(`${pair[0]} is ${pair[1]} <empty>, then return`);
			let inputName = document.querySelector(`input[name=${pair[0]}]`)
			inputName.style.outline="1px solid red"
			inputName.style.border="1px solid red"
			return
		}
		data.append(pair[0], pair[1])
	}

	fetch("https://gzhang.dev/tenda/api/picklist", {
		method: "POST",
		body: data,
	})
	.then(resp => { 
		return resp.json();
	})
	.then(data => {
		console.log("data:",data);
		console.log(data[0].lastId)
		if(data[0].lastId) {
			let ifb = document.querySelector("#insertFB")
			ifb.style.display="block"
			let table = document.querySelector("#insertFB table");

			let tbody = table.createTBody();
			let newRow = tbody.insertRow(0);
			let cell1 = newRow.insertCell(0);
			let cell2 = newRow.insertCell(1);
			let cell3 = newRow.insertCell(2);
			let cell4 = newRow.insertCell(3);
			let cell5 = newRow.insertCell(4);
			let cell6 = newRow.insertCell(5);
			let cell7 = newRow.insertCell(6);
			let cell8 = newRow.insertCell(7);
			let cell9 = newRow.insertCell(8);

			let fetch_url = `/tenda/api/picklist?PID=${data[0].lastId}&status=Pending`;
			fetch(fetch_url)
			.then(resp => {
				return resp.json();
			})
			.then(data =>{
				data.forEach( p=> {
					cell1.innerHTML = p.PID;
					cell2.innerHTML = p.PNO;
					cell3.innerHTML = p.model;
					cell4.innerHTML = p.qty;
					cell5.innerHTML = p.customer;
					cell6.innerHTML = p.location;
					cell7.innerHTML = p.status;
					cell8.innerHTML = p.created_at;
					cell9.innerHTML = `<a href='/tenda/update/picklist?PID=${p.PID}'>Modify</a>`;
				})
			})
		}
	})
});


