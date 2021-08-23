import { createTable, rebuild_dbtable } from "./helper.js"

/*
let fetch_url = `https://gzhang.dev/tenda/api/stock`

fetch(fetch_url)
.then( resp => {
	return resp.json();
})
.then( data => {
	let theadTitles = ["位置", "产品", "每箱数量", "箱数", "散货数", "数量总计", "品种", "备注"];
	let orders = ["location", "model", "unit", "cartons", "boxes", "total", "kind", "notes"];
	// data.forEach( d => {
	// 	d.update = `<a href="/tenda/stock/update?SID=${d.SID}">update</a>`
	// })
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
*/

console.log("-- soh.js --")
let fetch_url = `https://gzhang.dev/tenda/api/soh`

fetch(fetch_url)
.then( resp => {
	return resp.json();
})
.then(data => {
	let theadTitles = ["产品", "库存总计"];
	let orders = ["model", "totals"];
	// data.forEach( d => {
	// 	d.update = `<a href="/tenda/stock/update?SID=${d.SID}">update</a>`
	// })
	let newtable = createTable(theadTitles, data, orders);
	let insDom = document.querySelector("#dbtable_container")
	if (insDom.innerHTML !== "") {
		insDom.innerHTML = ""
	}
	newtable.id = "dbtable"
	insDom.appendChild(newtable);

	let table = $('#dbtable').DataTable({
		dom: 'Bfrtip',
		paging: true,
		buttons: ['excel'],
		order: [1, "des"],
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

