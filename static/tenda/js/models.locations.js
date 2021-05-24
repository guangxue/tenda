import { Tenda } from "./helper.js"

let modelinput = document.querySelector('input[name=model]')
let modelNameinput = document.querySelector('input[name=modelName]')
if(modelinput || modelNameinput) {
	fetch("https://gzhang.dev/tenda/api/model")
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


const modelInput = document.querySelector("input[name=modelName]");
// get locations for a model
if(modelInput) {
	modelInput.addEventListener('input', function(e) {
		let model = document.querySelector("input[name=modelName]").value;
		if(model) {
			fetch("https://gzhang.dev/tenda/api/locations?model="+model)
			.then( resp => {
				return resp.json();
			})
			.then( data => {
				let optionFragement = new DocumentFragment();
				let selectButton = document.querySelector('#pick-location');
				selectButton.innerHTML = ""
				data.forEach( loc => {
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