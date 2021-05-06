/** INSERT INTO picked **/
const pickButton = document.querySelector("#pickbtn")
let currentLN = 1;
pickButton.addEventListener('click', function(e) {
	e.preventDefault();
	console.log("Pick clicked")

	let formEl = document.querySelector('#pickform')
	let formData = new FormData(formEl);
	let data = new URLSearchParams(formData);

	for(const pair of formData) {
		data.append(pair[0], pair[1])
	}

	fetch("https://gzhang.dev/tenda/api/picklist", {
		method: "POST",
		body: data,
	})
	.then(resp => { 
		return resp.json();
	})
	.then(data => {
		console.log("data:",data);
		console.log(data[0].lastId)
		if(data[0].lastId) {
			let ifb = document.querySelector("#insertFB")
			ifb.style.display="block"
			let table = document.querySelector("#insertFB table");

			let tbody = table.createTBody();
			let newRow = tbody.insertRow(0);
			let cell1 = newRow.insertCell(0);
			let cell2 = newRow.insertCell(1);
			let cell3 = newRow.insertCell(2);
			let cell4 = newRow.insertCell(3);
			let cell5 = newRow.insertCell(4);
			let cell6 = newRow.insertCell(5);
			let cell7 = newRow.insertCell(6);
			let cell8 = newRow.insertCell(7);

			let fetch_url = `/tenda/api/picklist?PID=${data[0].lastId}&status=Pending`;
			fetch(fetch_url)
			.then(resp => {
				return resp.json();
			})
			.then(data =>{
				data.forEach( p=> {
					cell1.innerHTML = p.PID;
					cell2.innerHTML = p.PNO;
					cell3.innerHTML = p.model;
					cell4.innerHTML = p.qty;
					cell5.innerHTML = p.customer;
					cell6.innerHTML = p.location;
					cell7.innerHTML = p.status;
					cell8.innerHTML = p.created_at;
				})
			})
		}
	})
});


