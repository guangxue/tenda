import { createTable, rebuild_dbtable } from "./helper.js"

let fetch_url = `https://gzhang.dev/tenda/api/stock`
console.log("fetch_url:", fetch_url)

fetch(fetch_url)
.then( resp => {
	return resp.json();
})
.then( data => {
	console.log(data)
	let theadTitles = ["位置", "产品", "每箱数量", "箱数", "散货数", "数量总计", "品种", "备注", "Update"];
	let orders = ["location", "model", "unit", "cartons", "boxes", "total", "kind", "notes", "update"];
	data.forEach( d => {
		d.update = `<a href="/tenda/stock/update?SID=${d.SID}">update</a>`
	})
	let newtable = createTable(theadTitles, data, orders);
	let insDom = document.querySelector("#dbtable_container")
	if (insDom.innerHTML !== "") {
		insDom.innerHTML = ""
	}
	newtable.id = "dbtable"
	insDom.appendChild(newtable);

	let table = $('#dbtable').DataTable({
		dom: 'Bfrtip',
		ordering: false,
		buttons: ['excel'],
	});
	// table.column(columnNo).search('^\\s' + this.value +'\\s*$', true, false).draw();
	rebuild_dbtable();
	return table;
})
.then((tbl)=>{
	$('input[type=search]').on( 'input', function () {
	    tbl.search( '^'+this.value,true ).draw();
	});
});
// let ST_Btn = document.querySelector(".select-stocktakes")
// ST_Btn.addEventListener("click", function() {
// 	const tbvalue = document.querySelector("#tbvalue").value;
// 	let fetch_url = `https://gzhang.dev/tenda/stock?tbname=${tbvalue}`
// 	console.log("fetch_url:", fetch_url)
	
// 	fetch(fetch_url)
// 	.then( resp => {
// 		return resp.json();
// 	})
// 	.then( data => {
// 		console.log(data)
// 		let theadTitles = ["位置", "产品", "每箱数量", "箱数", "散货数", "数量总计", "品种", "备注", "Update"];
// 		let orders = ["location", "model", "unit", "cartons", "boxes", "total", "kind", "notes", "update"];
// 		data.forEach( d => {
// 			d.update = `<a href="/tenda/stock/update?tbname=${tbvalue}&SID=${d.SID}">update</a>`
// 		})
// 		let newtable = createTable(theadTitles, data, orders);
// 		let insDom = document.querySelector("#dbtable_container")
// 		if (insDom.innerHTML !== "") {
// 			insDom.innerHTML = ""
// 		}
// 		newtable.id = "dbtable"
// 		insDom.appendChild(newtable);

// 		let table = $('#dbtable').DataTable({
// 			dom: 'Bfrtip',
// 			ordering: false,
// 			buttons: ['excel']
// 		});
// 		rebuild_dbtable();
// 	});
// })
