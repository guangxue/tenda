import { getCurrentDateTime } from './helper.js';

let inputElems = document.querySelectorAll("input")
inputElems.forEach( elm => {
	elm.oninput = function() {
		console.log("elem on input:", elm);
		document.querySelector("#inputchanged").innerHTML = elm.value
	};
})
let updatePickedBtn = document.querySelector("#updatePickedButton");
updatePickedBtn.addEventListener('click', function(e) {
	e.preventDefault();
	let form = document.querySelector('#updatePickedForm')
	let formData = new FormData(form);
	let data = new URLSearchParams(formData);
	let timenow = getCurrentDateTime();
	console.log("timenow:", timenow);
	for(const pair of formData) {
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
		console.log("data:",data);
		console.log(data[0].lastId)
	})
})


