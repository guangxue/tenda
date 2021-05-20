import { formDataCollect, createTable } from './helper.js';


const updateStockBtn = document.querySelector("#updateStockBtn");

updateStockBtn.addEventListener('click', function(e) {
	e.preventDefault();
	let data = formDataCollect("#table-form");
	console.log("table-form,data:", data);
	let fetch_url = "https://gzhang.dev/tenda/api/stock/update";
	fetch(fetch_url, {
		method: "POST",
		body: data,
	})
	.then( resp => {
		return resp.json();
	})
	.then( resdata => {
		if(resdata) {
			let titles = ["SID", "location", "model", "unit", "cartons", "boxes", "total", "update_comments", "Actions"];
			let data = resdata;
			let currURL = window.location.pathname + window.location.search;
            let pathname = window.location.pathname;
			console.log("currURL:", currURL);
			data[0].Actions = `<a href="/tenda/api/txcm?cmname=StockUpdate&urlname=${currURL}&SID=${data[0].SID}">Confirm</a> <a href="/tenda/api/txrb?rbname=StockUpdate&urlname=${currURL}&SID=${data[0].SID}">Discard</a>`
			let table = createTable(titles, data, titles);
			let contentwrapper = document.querySelector(".content-wrapper");
			contentwrapper.appendChild(table);
		}
	})
});
