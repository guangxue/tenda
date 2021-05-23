import { formDataCollect, createTable,endTransact } from './helper.js';


const updateStockBtn = document.querySelector("#updateStockBtn");

updateStockBtn.addEventListener('click', function(e) {
	e.preventDefault();
    let SID = document.querySelector("input[name=SID]").value
	let data = formDataCollect("#table-form");
	console.log("table-form,data:", data);
	let fetch_url = `https://gzhang.dev/tenda/api/stock/SID/${SID}`;
	fetch(fetch_url, {
		method: "PUT",
		body: data,
	})
	.then( resp => {
		return resp.json();
	})
	.then( resdata => {
		console.log("resdata for update stock:", resdata);
		if(resdata) {
			let titles = ["SID", "location", "model", "unit", "cartons", "boxes", "total", "update_comments", "Actions"];
			let data = resdata;
			let currURL = window.location.pathname + window.location.search;
            let pathname = window.location.pathname;
			console.log("currURL:", currURL);
			data[0].Actions = `<a href="#" id="txcm" data-txname="StockUpdate">Confirm</a> <a href="#" data-txname="StockUpdate" id="txrb">Discard</a>`
			let table = createTable(titles, data, titles);
			table.classList.add("single-row")
			let contentwrapper = document.querySelector(".content-wrapper");
			contentwrapper.appendChild(table);
			endTransact("#txcm", "#txrb", table)
		}
	})
});
