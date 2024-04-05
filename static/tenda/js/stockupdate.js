import { formDataCollect, createTable,endTransact } from './helper.js';


const updateStockBtn = document.querySelector("#updateStockBtn");
const deleteStockBtn = document.querySelector("#delStockBtn")
const inputCartons = document.querySelector("input[name=cartons]")
const inputBoxes = document.querySelector("input[name=boxes]")
const inputTotal = document.querySelector("input[name=total]")
// inputTotal.disabled = true;


inputCartons.addEventListener("change", (e)=> {
	const unit = document.querySelector("input[name=unit]").value;
	let cartons = inputCartons.value;
	let boxes = inputBoxes.value;
	console.log("unit:", unit)
	console.log("cartons:", cartons)
	console.log("boxes:", boxes)
	let total = parseInt(unit) * parseInt(cartons) + parseInt(boxes);
	console.log("total:", total)
	inputTotal.value = total;
});

inputBoxes.addEventListener("change", (e)=> {
	const unit = document.querySelector("input[name=unit]").value; 
	let cartons = inputCartons.value;
	let boxes = inputBoxes.value;
	console.log("unit:", unit)
	console.log("cartons:", cartons)
	console.log("boxes:", boxes)
	let total = parseInt(unit) * parseInt(cartons) + parseInt(boxes);
	console.log("total:", total)
	inputTotal.value = total;	
});

updateStockBtn.addEventListener('click', function(e) {
	e.preventDefault();
    let SID = document.querySelector("input[name=SID]").value
    let updateCode = document.querySelector("input[name=updateCode]").value
	let data = formDataCollect("#table-form");
	let fetch_url = `/api/stock/SID/${SID}`;
	fetch(fetch_url, {
		method: "PUT",
		body: data,
	})
	.then( resp => {
		return resp.json();
	})
	.then( resdata => {
		if(resdata) {
			let fbelm = document.querySelector(".update-fd");
			let titles = ["SID", "location", "model", "unit", "cartons", "boxes", "total", "kind","update_comments", "Actions"];
			let data = resdata;
			let currURL = window.location.pathname + window.location.search;
            let pathname = window.location.pathname;
			data[0].Actions = `<a href="#" id="txcm" data-txname="StockUpdate">Confirm</a> <a href="#" data-txname="StockUpdate" id="txrb">Discard</a>`
			let table = createTable(titles, data, titles);
			table.classList.add("single-row")
			let contentwrapper = document.querySelector(".content-wrapper");
			contentwrapper.appendChild(table);
			endTransact("#txcm", "#txrb", table, fbelm)
		}
	})
});

deleteStockBtn.addEventListener('click', function(e) {
	e.preventDefault();
	let SID = document.querySelector("input[name=SID]").value;
	let fetch_url = `/api/stock/SID/${SID}`;
	console.log("Deleting SID:", SID)
	fetch(fetch_url, {
		method: "DELETE",
		body:null,
	})
});

