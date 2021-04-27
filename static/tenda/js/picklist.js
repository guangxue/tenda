fetch("https://gzhang.dev/tenda/api/picked?date=today")
.then( resp => {
	return resp.json();
})
.then( data => {
	console.log(data)
	let tableBody = document.querySelector('#picklist_table tbody')
	tableBody.innerHTML = ""
	let tplrow = document.querySelector('#htmpl_pick');
	
	data.forEach(p=> {
		var row = tplrow.content.cloneNode(true);
		var td = row.querySelectorAll('td');
		td.item(0).textContent = p.PID;
		td.item(1).textContent = p.PNO;
		td.item(2).textContent = p.model;
		td.item(3).textContent = p.qty;
		td.item(4).textContent = p.customer;
		td.item(5).textContent = p.location;
		td.item(6).textContent = p.status;
		td.item(7).textContent = p.last_updated;
		td.item(8).innerHTML = `<a href="/tenda/update/picked?PID=${p.PID}">update</a>`;
		tableBody.appendChild(row);
	})
	
})