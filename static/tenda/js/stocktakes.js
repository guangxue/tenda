document.querySelector("#tbvalue").addEventListener('change', function() {
	window.location.reload(false)
})
let ST_Btn = document.querySelector(".select-stocktakes")
ST_Btn.addEventListener("click", function() {
	const tbvalue = document.querySelector("#tbvalue").value;
	let fetch_url = `https://gzhang.dev/tenda/stocktakes?tbname=${tbvalue}`
	console.log("fetch_url:", fetch_url)
	fetch(fetch_url)
	.then( resp => {
		return resp.json();
	})
	.then( data => {
		console.log(data)
		let tableBody = document.querySelector('tbody')
		tableBody.innerHTML = ""
		let tplrow = document.querySelector('#htmpl_pick');
		
		data.forEach(p=> {
			var row = tplrow.content.cloneNode(true);
			var td = row.querySelectorAll('td');
			td.item(0).textContent = p.location;
			td.item(1).textContent = p.model;
			td.item(2).textContent = p.unit;
			td.item(3).textContent = p.cartons;
			td.item(4).textContent = p.boxes;
			td.item(5).textContent = p.total;
			td.item(6).textContent = p.kind;
			td.item(7).textContent = p.notes;
			tableBody.appendChild(row);
		})
	})
	.then(()=>{
		
		$('#stock_tb').DataTable({
			dom: 'Bfrtip',
			"ordering": false,
			buttons: [
				'excel'
			]
		});
	})
})