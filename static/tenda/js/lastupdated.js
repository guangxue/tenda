import { formDataCollect } from './helper.js';

let lidUpdateBtn = document.querySelector("#updatePickedButton");

lidUpdateBtn.addEventListener("click", function(e) {
	e.preventDefault();
	console.log("updated LID button");
	let data = formDataCollect("#updateLIDForm");
	console.log(data);
	
});