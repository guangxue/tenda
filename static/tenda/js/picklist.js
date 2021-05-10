import { createTable, rebuild_dbtable, getCurrentDateTime } from "./helper.js"


const selectButton= document.querySelector(".picklist-select-btn")

selectButton.addEventListener("click", function() {
	let pickDate = document.querySelector("#pick_date").value
	let pickStatus = document.querySelector("#pick_status").value;

	console.log("pick date is:", pickDate)
	console.log("pick status is:", pickStatus)

	var fetch_url = `https://gzhang.dev/tenda/api/picklist?date=${pickDate}&status=${pickStatus}`
	console.log("fetch_url:", fetch_url)
	if(pickStatus == 'completed_at' && !pickDate) {
		console.log("return: pickDate is empty")
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
		.then( data => {
			console.log(data[0]);
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
		})
		.then( tableObj => {
			let newtable = createTable(tableObj.titles, tableObj.data, tableObj.names)
			let dbtable_container = document.querySelector("#dbtable_container");
			let insDom = document.querySelector("#dbtable_container")
			if (insDom.innerHTML !== "") {
				insDom.innerHTML = ""
			}
			newtable.id = "dbtable"
			insDom.appendChild(newtable);
			$("#dbtable").DataTable({
				dom: 'Bfrtip',
				buttons: ['print'],
				order: [5, "des"],
			});

			let table_width = rebuild_dbtable();
			return table_width;
		})
		.then((tw)=>{
			console.log("=> PENDING pickStatus:", pickStatus)
			let cmpbtn1 = document.querySelector("#completeBtn")
			console.log("=> CMB :", cmpbtn1)
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
				// if (tableBody && trows.length) {
				// 	fetch("https://gzhang.dev/tenda/api/complete/picklist", {
				// 		method: "POST",
				// 		body: data,
				// 	})
				// 	.then(resp => { 
				// 		return resp.json();
				// 	})
				// 	.then(data => {
				// 		console.log("data:",data);
				// 	})
				// }
			})
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


