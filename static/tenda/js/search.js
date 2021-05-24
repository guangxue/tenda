console.log("[search.js]")
let searchBtn = document.querySelector('#search-btn');
let inputModel = document.querySelector("#mod");
let inputLocation = document.querySelector("#loc");
const getModel = "https://gzhang.dev/tenda/api/model/";
const getTheModel = "https://gzhang.dev/tenda/api/model";
const getModelsWhereLocation = "https://gzhang.dev/tenda/api/model?location=";

searchBtn.addEventListener('click', function(e) {
	e.preventDefault();
	let model = inputModel.value;
	let location = inputLocation.value;

	if(model && location) {
		let url = `${getTheModel}/model?${model}&location=${location}`
		console.log("url->", url)
		fetch(url).then(resp=>{resp.json()})
		.then(data=> {
			console.log("data the model", data)
		})
	}

	if(model) {
		console.log("[search.js] searchModel url:", getModel+model)
		fetch(getModel+model).then(response => { return response.json()})
		.then(data => {
			let sum_total = 0;
			let table = "<table><thead><tr><th>Location</th><th>Unit</th><th>Cartons</th><th>Boxes</th><th>Total</th></tr></thead><tbody>";
			data.forEach(m => {
				let row = `<tr><td>${m.location}</td><td>${m.unit}</td><td>${m.cartons}</td><td>${m.boxes}</td><td>${m.total}</td></tr>`
				table += row
				sum_total += parseFloat(m.total);
			});
			table += `<tr><td></td><td></td><td></td><td></td><td>${sum_total}</td></tr></tbody></table>`
			document.querySelector("#sfb").innerHTML = table;
		})
		.then(()=>{
			let tableTitle = document.querySelector(".model-title")
			tableTitle.innerHTML = `<div class="table-title">Model: ${model}</div>`
		})
		.catch( err => {
			console.log("model query FAILED:", err);
		});
	}
	if(location) {
		let url = getModelsWhereLocation+location;
		console.log("[search.js] searchLocation url:", url)
		fetch(url)
		.then(resp => {
			return resp.json();
		})
		.then(data => {
			console.log("data: json:", data)
			let table = "<table><thead><tr><th>model</th><th>Unit</th><th>Cartons</th><th>Loose</th></tr></thead></tbody>";
			data.forEach( m=> {
				let row = `<tr><td>${m.model}</td><td>${m.unit}</td><td>${m.cartons}</td><td>${m.boxes}</td></tr>`
				table += row
			})
			document.querySelector("#sfb").innerHTML = table;
		})
		.then(()=> {
			let tableTitle = document.querySelector(".location-title")
			tableTitle.innerHTML = `<div class="table-title">Model: ${location}</div>`
		})
		
	}
});