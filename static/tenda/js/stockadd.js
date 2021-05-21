import {formDataCollect, createTable} from "./helper.js"

let stockAddBtn = document.querySelector("#stockAddBtn")

stockAddBtn.addEventListener('click', function(e) {
    e.preventDefault();
	let fdata = formDataCollect("#stockadd-form")
    console.log("fdata:", fdata)
	fetch("https://gzhang.dev/tenda/api/stock", {
		method: "POST",
		body: fdata,
	})
	.then(resp => { return resp.json() })
	.then(data => {
		console.log(data);
        if(data && data[0].SID) {
            data[0].confirm = `<a href="/tenda/api/txcm?cmname=StockAdd">Confirm</a> <a href="/tenda/api/txrb?rbname=StockAdd">Discard</a>`
            data[0].SID=data[0].lastId;
            let titles = ['SID','location','model','unit','cartons','boxes','total', 'confirm']; 
            let table = createTable(titles,data,titles);
            let tblformwrp = document.querySelector(".table-form-wrapper");
            tblformwrp.appendChild(table);
        }
	})
});
