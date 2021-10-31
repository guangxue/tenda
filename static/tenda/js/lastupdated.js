import { formDataCollect } from './helper.js';

let lidUpdateBtn = document.querySelector("#updatePickedButton");

lidUpdateBtn.addEventListener("click", function(e) {
	e.preventDefault();
	console.log("updated LID button");
	let data = formDataCollect("#updateLIDForm");
	let lid = data.get("LID");
	fetch("https://gzhang.dev/tenda/api/lastupdated/LID/"+lid).then( resp => {
		return resp.json();
	})
	.then( data => {
		console.log("lastupdated:=>", data)
	})
});