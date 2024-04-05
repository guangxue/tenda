import { Tenda } from "./helper.js"

let inputModel = document.querySelector('input[name=model]')
if(inputModel) {
  console.log(inputModel)
	fetch("/api/model")
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

if(inputModel) {
	inputModel.addEventListener('input', function(e) {
		let model = document.querySelector("input[name=model]").value;
		if(model) {
			fetch("/api/locations?model="+model)
			.then( resp => {
				return resp.json();
			})
			.then( data => {
				let optionFragement = new DocumentFragment();
				let selectButton = document.querySelector('#pick-location');
				if(selectButton) { selectButton.innerHTML = ""; }
				
				data.forEach( loc => {
					// let currentOpt = '<option value="'+loc.location+'">'+loc.location+"</option>";
					let opt = document.createElement('option');
					opt.value = loc.location;
					opt.textContent = loc.location;
					optionFragement.appendChild(opt);
				});
				if(selectButton) { selectButton.appendChild(optionFragement) }
			})
		}
	})
}