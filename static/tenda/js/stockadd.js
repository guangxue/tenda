import {formDataCollect, createTable, endTransact} from "./helper.js"

let stockAddBtn = document.querySelector("#stockAddBtn")

stockAddBtn.addEventListener('click', function(e) {
    e.preventDefault();
	let fdata = formDataCollect("#stockadd-form")
	fetch("/api/stock", {
		method: "POST",
		body: fdata,
	})
	.then(resp => { return resp.json() })
	.then(data => {
        if(data && data[0].SID) {
            data[0].confirm = `<a href="#" id="txcm" data-txname="StockAdd">Confirm</a> <a href="#" id="txrb" data-txname="StockAdd">Discard</a>`
            let titles = ['SID','location','model','unit','cartons','boxes','total','kind', 'confirm']; 
            let table = createTable(titles,data,titles);
            let tblformwrp = document.querySelector(".table-form-wrapper");
            tblformwrp.appendChild(table);
            return table
        }
	})
	.then((tbl)=>{
		endTransact("#txcm", "#txrb", tbl)
	})
});
