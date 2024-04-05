import {createTable} from "./helper.js"

console.log("[search.js]")

let searchBtn = document.querySelector('#search-btn');
let inputModel = document.querySelector("#mod");
let inputLocation = document.querySelector("#loc");
let sfb = document.querySelector("#sfb")

const getModel = "/api/model/";
const getTheModel = "/api/model";
const getModelsWhereLocation = "/api/model?location=";

searchBtn.addEventListener('click', function(e) {
	e.preventDefault();
	let model = inputModel.value;
	let location = inputLocation.value;

	if(model && location) {
		console.log("model & location query")
		let url = `${getTheModel}/${model}?location=${location}`
		fetch(url).then(resp=>{ return resp.json()})
		.then(data=> {
            data[0].update =`<a href="/stock/update?SID=${data[0].sid}">Update</a>`;;
			console.log("data the model", data)
            let titles = ['sid','location', 'unit', 'cartons', 'boxes', 'total','update'];
            sfb.innerHTML = "" 
            let tbl = createTable(titles, data, titles);
            sfb.appendChild(tbl);
            
		})
	}
	if(model && location === "") {
		console.log("model & NO location query")
		console.log("[search.js] searchModel url:", getModel+model)
		fetch(getModel+model)
    .then(res => {
      return res.json()
    })
		.then(data => {
      console.log(data)
			let sum_total = 0;
			let table = "<table><thead><tr><th>Location</th><th>Unit</th><th>Cartons</th><th>Boxes</th><th>Total</th><th>Modify</th></tr></thead><tbody>";
			data.forEach(m => {
				let row = `<tr><td>${m.location}</td><td>${m.unit}</td><td>${m.cartons}</td><td>${m.boxes}</td><td>${m.total}</td><td><a href="/stock/update?SID=${m.SID}">Update</a></td></tr>`
				table += row
				sum_total += parseFloat(m.total);
			});
			table += `<tr><td></td><td></td><td></td><td></td><td>${sum_total}</td><td></td></tr></tbody></table>`
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
	if(location && model === "") {
		let url = getModelsWhereLocation+location;
		console.log("location & NO model query")
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
