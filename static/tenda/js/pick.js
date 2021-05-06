/** INSERT INTO picked **/
const pickButton = document.querySelector("#pickbtn")
pickButton.addEventListener('click', function(e) {
	e.preventDefault();
	console.log("Pick clicked")
	let pickTable = document.querySelector("#pick-table");
	let newRow = pickTable.insertRow(-1);
	
	// let formEl = document.querySelector('#pickform')
	// let formData = new FormData(formEl);
	// let data = new URLSearchParams(formData);

	// for(const pair of formData) {
	// 	data.append(pair[0], pair[1])
	// }

	// fetch("https://gzhang.dev/tenda/api/picklist", {
	// 	method: "POST",
	// 	body: data,
	// })
	// .then(resp => { 
	// 	return resp.json();
	// })
	// .then(data => {
	// 	console.log("data:",data);
	// 	console.log(data[0].lastId)
	// })
});


