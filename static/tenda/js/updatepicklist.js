import { getCurrentDateTime, createTable } from './helper.js';

let updatePickedBtn = document.querySelector("#updatePickedButton");
let updateFD = document.querySelector('.update-fd')

updatePickedBtn.addEventListener('click', function(e) {
	e.preventDefault();

	let form = document.querySelector('#updatePickedForm')
	let formData = new FormData(form);
	let data = new URLSearchParams(formData);
	let timenow = getCurrentDateTime();
	console.log("timenow:", timenow);
	for(const pair of formData) {
		if(!pair[1]) {
			console.log(`${pair[0]} is ${pair[1]} <empty>, then return`);
			let inputName = document.querySelector(`input[name=${pair[0]}]`)
			inputName.style.border="1px solid red"
			return
		}
		console.log(`${pair[0]} is ${pair[1]}`)
		data.append(pair[0], pair[1])
	}
	data.append("timenow", timenow)

	fetch("https://gzhang.dev/tenda/update/picklist", {
		method: "POST",
		body: data,
	})
	.then(resp => { 
		return resp.json();
	})
	.then(data => {
		console.log("Updated picked:", data);
		if(data[0].PNO) {
			// updateFD.innerHTML = `PID: ${data[0].PNO} update successfully`;
			// updateFD.style.opacity  = 1
			let pid = document.querySelector("input[name=PID]").value
			let titles = ['PNO', 'customer', 'model', 'qty', 'location', 'status', 'confirm'];
			data.forEach( d=> {
				d.confirm = `<a href="/tenda/api/txcm?cmname=PickList&UPID=${pid}">Confirm</a> <a href="/tenda/api/txrb?rbname=PickList&UPID=${pid}">Discard</a>`
			})
			let table = createTable(titles, data, titles);
			let updated = document.querySelector("#updated-picked")
			updated.appendChild(table);
		}
	})
});


let delBtn = document.querySelector('#delBtn');
delBtn.addEventListener("click", function(e) {
	e.preventDefault();
	let form = document.querySelector('#updatePickedForm')
	let formData = new FormData(form);
	let data = new URLSearchParams(formData);
	console.log("data:",data.toString());

	console.log("[clicked] delBtn");
	fetch("https://gzhang.dev/tenda/api/picked/delete", {
		method: "POST",
		body: data,
	});
});
