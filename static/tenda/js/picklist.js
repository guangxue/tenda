import { createTable, rebuild_dbtable, getCurrentDateTime } from "./helper.js"


const selectButton= document.querySelector(".picklist-select-btn");
const pickStatusOpt = document.querySelector("#pick_status");


pickStatusOpt.addEventListener("change", function() {
	let pickDate = document.querySelector("#pick_date");
	if(pickStatusOpt.value.includes("_at")) {
		console.log("option:", pickStatusOpt.value)
		
		pickDate.setAttribute("min", "2021-04-04");
		pickDate.setAttribute("step", "7");
	} else {
		pickDate.removeAttribute("min")
		pickDate.removeAttribute("step")
	}
});

selectButton.addEventListener("click", function() {
	let pickDate = document.querySelector("#pick_date").value
	let pickStatus = document.querySelector("#pick_status").value;

	console.log("pick date is:", pickDate)
	console.log("pick status is:", pickStatus)

	var fetch_url = `https://gzhang.dev/tenda/api/picklist?date=${pickDate}&status=${pickStatus}`
	console.log("fetch_url:", fetch_url)
	let dbtable_container = document.querySelector("#dbtable_container");
	if (dbtable_container.innerHTML !== "") {
		dbtable_container.innerHTML = ""
	}
	if(!pickDate) {
		dbtable_container.innerHTML = `<div class="alert alert-danger">error: pick date is empty</div>`;
		let alertDanger = document.querySelector("#dbtable_container div.alert-danger")
		alertDanger.style.width  = "500px";
		alertDanger.style.opacity  = 1;
		return;
	}
	let cmpbtn = document.querySelector("#completeBtn")
	if(cmpbtn) {
		cmpbtn.remove();
	}
	if(pickStatus) {
		console.log("pick status is:", pickStatus)
		fetch(fetch_url)
		.then( resp => {
			return resp.json();
		})
		.then( data  => {
			if(!data[0]) {
				console.log("!data")
				const err = new Error("no rows return from database table: `picklist`")
				err.name = "Empty set"
				throw err;
			}
			if(data[0].PID) {
				let titles = ["PNO","customer", "model", "quantity", "status","location", "created_at", "Action"];
				let objNames =   ["PNO","customer", "model", "qty", "status","location", "created_at", "update"];
				let tbData = data;
				tbData.forEach( d => {
					d.update = `<a href="/tenda/update/picklist?PID=${d.PID}">update</a>`;
				})
				return {
					"titles": titles,
					"data": tbData,
					"names": objNames,
				}
			}
			if(data[0].LID) {
				let titles = ["LID","location", "model", "cartons", "boxes", "completed_at", "Action"];
				let objNames =   ["LID","location", "model", "cartons", "boxes", "completed_at", "update"];
				let tbData = data;
				tbData.forEach( d => {
					d.update = `<a href="/tenda/update/picklist?LID=${d.LID}">update</a>`;
				});
				return {
					"titles": titles,
					"data": tbData,
					"names": objNames,
				}
			}
			if(data[0].model) {
				let titles = ["item", "model", "total"];
				let tbData = data;
				let i = 0;
				tbData.forEach( d=> {
					d.item = i + 1;
					i++;
				})
				let objNames = ["item", "model", "total"];
				return {
					"titles":titles,
					"data": tbData,
					"names": objNames,
				}
			}
			
		})
		.then( tableObj => {
			let newtable = createTable(tableObj.titles, tableObj.data, tableObj.names)
			let dbtable_container = document.querySelector("#dbtable_container");
			if (dbtable_container.innerHTML !== "") {
				dbtable_container.innerHTML = ""
			}
			newtable.id = "dbtable"
			console.log("newtable:", newtable)
			dbtable_container.appendChild(newtable);
			$("#dbtable").DataTable({
				dom: 'Bfrtip',
				buttons: ['print'],
				// order: [5, "des"],
			});

			let table_width = rebuild_dbtable();
			return table_width;
		})
		.then((tw)=>{
			console.log("=> PENDING pickStatus:", pickStatus)
			let cmpbtn1 = document.querySelector("#completeBtn")
			if(!cmpbtn1 && pickStatus === "Pending") {
				let completeButton = document.createElement("button")
				completeButton.id = "completeBtn"
				completeButton.classList.add("btn", "btn-table")
				completeButton.textContent = "Complete"
				let dbtable_info = document.querySelector("#dbtable_info");
				let actionbtn_wrapper = document.createElement("div");
				actionbtn_wrapper.id="actionbtn_wrapper"
				actionbtn_wrapper.style.width=`${tw}px`
				actionbtn_wrapper.appendChild(completeButton)
				dbtable_info.insertAdjacentElement('afterend', actionbtn_wrapper)
			}
		})
		.then(()=>{
			// Add Event Listenser to appened `complete` Button
			const completePickBtn = document.querySelector('#completeBtn')
			if(!completePickBtn) {
				return
			}
			completePickBtn.addEventListener('click', function() {
				let formData = new FormData();
				let data = new URLSearchParams(formData);
				let pickStatus = document.querySelector('#pick_status').value
				let pickDate   = document.querySelector("#pick_date").value
				console.log("pick status to Complete:", pickStatus);
				console.log("pick date to complete", pickDate);
				console.log("lastSaturday", lastSaturdayTS());

				data.append("pickDate", pickDate)
				data.append("pickStatus", pickStatus)
				data.append("lastSaturday", lastSaturdayTS());

				let tableBody = document.querySelector('table tbody')
				let trows = tableBody.querySelectorAll('tr');
				// console.log("Complete Button Disabled");
				if (tableBody && trows.length) {
					fetch("https://gzhang.dev/tenda/api/complete/picklist", {
						method: "POST",
						body: data,
					})
					.then(resp => { 
						return resp.json();
					})
					.then(data => {
						console.log("data:",data);
					})
				}
			})
		})
		.catch(err => {
			console.log("err:", err)
			dbtable_container.innerHTML = `<div class="alert alert-info">${err.name} : ${err.message}`;
			let alertInfo = document.querySelector("#dbtable_container div.alert-info")
			alertInfo.style.width  = "500px";
			alertInfo.style.opacity  = 1;
		})
	}
});

function lastSaturdayTS() {
	const t = new Date().getDate() + (6 - new Date().getDay() - 1) - 6 ;
	const lfri = new Date();
	lfri.setDate(t);
	var month = lfri.getMonth()+1;
	if(month < 10) {
		month = `0${month}`
	}
	let day = lfri.getDate();
	if (day < 10) {
		day = `0${day}`
	}
	let lastSaturday = `${lfri.getFullYear()}-${month}-${day}`;
	console.log(lastSaturday);
	return lastSaturday;
}


