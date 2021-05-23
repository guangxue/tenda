let queryButton = document.querySelector('#querybutton');

queryButton.addEventListener('click', function(e) {
	e.preventDefault();
	let model = document.querySelector('#mod').value;
	let location = document.querySelector('#loc').value;
	if(model) {
		let url = "https://gzhang.dev/tenda/api/models?model="+model;
		fetch(url).then(response => { return response.json()})
		.then(data => {
			let sum_total = 0;
			let table = "<table><thead><tr><th>Location</th><th>Unit</th><th>Cartons</th><th>Boxes</th><th>Total</th></tr></thead><tbody>";
			data.forEach(m => {
				let row = `<tr><td>${m.location}</td><td>${m.unit}</td><td>${m.cartons}</td><td>${m.boxes}</td><td>${m.total}</td></tr>`
				table += row
				sum_total += parseFloat(m.total);
			});
			table += `<tr><td></td><td></td><td></td><td></td><td>${sum_total}</td></tr></tbody></table>`
			document.querySelector("#mfb").innerHTML = table;
		})
		.then(()=>{
			let tableTitle = document.querySelector(".model-title")
			tableTitle.innerHTML = `<div class="table-title">Model: ${model}</div>`
		})
		.catch( err => {
			console.log("model query FAILED:", err);
		});
	}
	if(location) {
		console.log("location->", location);
		let url = "https://gzhang.dev/tenda/api/models?location="+location;
		console.log("fetch url:", url);
		fetch(url).then(response => {
			return response.json()
		})
		.then(data => {
			console.log("location json:", data);
			let table = "<table><thead><tr><th>model</th><th>Unit</th><th>Cartons</th><th>Loose</th></tr></thead></tbody>";
			data.forEach( m=> {
				let row = `<tr><td>${m.model}</td><td>${m.unit}</td><td>${m.cartons}</td><td>${m.boxes}</td></tr>`
				table += row
			})
			document.querySelector("#lfb").innerHTML = table;
		})
		.then(()=> {
			let tableTitle = document.querySelector(".location-title")
			tableTitle.innerHTML = `<div class="table-title">Model: ${location}</div>`
		})
		.catch( err => {
			console.log("location query FAILED:", err);
		});
	}
});