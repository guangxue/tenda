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
const confirmButton = document.querySelector("#cfmbtn")
// const tbody = document.querySelector("#insertFB table tbody")
// let fbTable = document.createElement("table");
// let fbThead = fbTable.createTHead();
// let fbThrow = fbThead.insertRow(0);
// let cell0 = fbThrow.insertCell(0);
// let cell1 = fbThrow.insertCell(1);
// let cell2 = fbThrow.insertCell(2);
// let cell3 = fbThrow.insertCell(3);
// let cell4 = fbThrow.insertCell(4);
// let cell5 = fbThrow.insertCell(5);
// let cell6 = fbThrow.insertCell(6);
// let cell7 = fbThrow.insertCell(7);
// let cell8 = fbThrow.insertCell(8);
// let cell9 = fbThrow.insertCell(9);

// cell0.textContent = "PID";
// cell1.textContent = "PNO";
// cell2.textContent = "Sales Manager";
// cell3.textContent = "model";
// cell4.textContent = "qty";
// cell5.textContent = "customer";
// cell6.textContent = "location";
// cell7.textContent = "status";
// cell8.textContent = "created_at";
// cell9.textContent = "Action";


/**
 * input[PO Number]
 * WHEN input changes, get data from picklist
 * WHEN PNO number == input[PO Number] 
 **/
 /******
 * 
 *  <th>PID</th>
	<th>Sales</th>
	<th>PNO</th>
	<th>model</th>
	<th>qty</th>
	<th>customer</th>
	<th>location</th>
	<th>status</th>
	<th>created_at</th>
	<th>Action</th>
 * 
 *  let names = ["PID", ];
 * */
inputPNO.addEventListener("input", function() {
	if(inputPNO.value) {
		let fetch_url = `https://gzhang.dev/tenda/api/picklist/PNO/${inputPNO.value}`;
		fetch(fetch_url, {
			method: "GET"
		})
		.then(resp => {
			return resp.json();
		})
		.then(data => {
			if(data[0] && data) {
				console.log(data)
				let names = ["PID", "PNO","sales_mgr","customer", "model", "qty", "location", "status", "created_at", "Action"];
				let titles = ["PID", "PNO", "Sales", "Customer", "Model", "Quantity", "Location", "Status", "Created_at", "Action"]
				data.forEach( d=> {
					d.Action = `<a href="/tenda/picklist/update?PID=${d.PID}">Modify</a>`;
				});
				let table = createTable(titles, data, names);
				let insertFB = document.querySelector("#insertFB")
				insertFB.innerHTML = table.outerHTML;
				insertFB.style.display = "block"
			} else {
				// insertFB.innerHTML = ""
				// insertFB.appendChild(fbTable)
				// insertFB.style.display = "none"
				let lastTbody = document.querySelector("#insertFB table tbody")
				if(lastTbody) {
					lastTbody.innerHTML = ""
				}
			}
		})
	}
	else {
		insertFB.innerHTML = ""
		// insertFB.appendChild(fbTable);
		// fbTable.style.display="none";
	}
});


/**
 * WHEN pick button pressed
 **/
pickButton.addEventListener('click', function(e) {
	e.preventDefault();

	let formEl = document.querySelector('#pickform')
	let formData = new FormData(formEl);
	let data = new URLSearchParams(formData);

	for(const pair of formData) {
		if(!pair[1]) {
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
		console.log("[last insert ID]:", data[0].lastId)
		if(data[0].lastId) {
			let ifb = document.querySelector("#insertFB")
			ifb.style.display="block"
			let table = document.querySelector("#insertFB table");
			table.style.display = "block"
			let tbody =  document.querySelector("#insertFB tbody");
			console.log("tbody:",tbody)
			// let tbody = table.createTBody();
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
			let cell10 = newRow.insertCell(9);

			let fetch_url = `/tenda/api/picklist?PID=${data[0].lastId}&status=Pending`;
			fetch(fetch_url)
			.then(resp => {
				return resp.json();
			})
			.then(data =>{
				console.log("2.[packingslip.js] data from picklist table:", data)
				data.forEach( p=> {
					cell1.innerHTML = p.PID;
					cell2.innerHTML = p.PNO;
					cell3.innerHTML = p.sales_mgr;
					cell4.innerHTML = p.customer;
					cell5.innerHTML = p.model;
					cell6.innerHTML = p.qty;
					cell7.innerHTML = p.location;
					cell8.innerHTML = p.status;
					cell9.innerHTML = p.created_at;
					cell10.innerHTML = `<a href='/tenda/picklist/update?PID=${p.PID}' target="_blank">Modify</a>`;
				})
			})
		}
	})
});


/**
 * WHEN confirm button pressed
 **/

// confirmButton.addEventListener("click", function(e) {});