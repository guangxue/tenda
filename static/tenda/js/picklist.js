import { createTable, rebuild_dbtable, getCurrentDateTime,fadeOut,lastSaturdayTS, lastSunTS } from "./helper.js"


const selectButton= document.querySelector(".picklist-select-btn");
const pickStatusOpt = document.querySelector("#pick_status");
const pickDate = document.querySelector("#pick_date");
const modelName = document.querySelector("input[name=model]");
const PNOnumber = document.querySelector("input[name=PNO]");

modelName.addEventListener("input", function() {
	PNOnumber.value = ""
})

PNOnumber.addEventListener("input", function() {
	modelName.value = ""
})

if(!pickStatusOpt.value.includes("weekly")) {
	pickDate.removeAttribute("min")
	pickDate.removeAttribute("step")
}

pickStatusOpt.addEventListener("change", function() {
	// weekly picked
	if(pickStatusOpt.value.includes("weekly")) {
		let lastSun = lastSunTS();
		// pickDate.setAttribute("min", lastSun);
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
	let pickModel = document.querySelector("input[name=model]").value;
	let searchPNO = document.querySelector("input[name=PNO]").value;


	var fetch_url = `https://gzhang.dev/tenda/api/picklist?date=${pickDate}&status=${pickStatus}`
	// SELECT * FROM picklist where model='pickmodel' AND created_at > 'pickDate'
	if(pickModel && pickStatus == "from") {
		fetch_url = `https://gzhang.dev/tenda/api/picklist/model/${pickModel}?date=${pickDate}&status=${pickStatus}`
	}
	// SELECT * FROM picklist where model='pickmodel'
	if(pickModel && pickStatus != "from") {
		fetch_url = `https://gzhang.dev/tenda/api/picklist/model/${pickModel}`
	}
	// SELECT * FROM picklist where PNO = 'searchPNO'
	if(searchPNO) {
		fetch_url = `https://gzhang.dev/tenda/api/picklist/search/PNO/${searchPNO}`
	}
	let dbtable_container = document.querySelector("#dbtable_container");
	if (dbtable_container.innerHTML !== "") {
		dbtable_container.innerHTML = ""
	}
	if((!pickDate && pickStatus == "weeklycompleted") || (!pickDate && pickStatus == "weeklypicked")) {
		dbtable_container.innerHTML = `<div class="alert alert-danger">error: pick date is empty</div>`;
		let alertDanger = document.querySelector("#dbtable_container div.alert-danger")
		alertDanger.style.width  = "500px";
		fadeOut(alertDanger)
		return;
	}
	let cbtn = document.querySelector("#completeBtn")
	if(cbtn) {
		cbtn.remove();
	}
	if(pickStatus) {
		fetch(fetch_url)
		.then( resp => {
			return resp.json();
		})
		.then( data  => {
			if(!data[0]) {
				const err = new Error("no rows return from database table: `picklist`")
				err.name = "Empty set"
				throw err;
			}
			if(data[0].PID) {
				let titles = ["PNO","Sales Manager","customer", "model", "quantity", "status","location", "created_at", "Action"];
				let objNames =   ["PNO","sales_mgr", "customer","model", "qty", "status","location", "created_at", "update"];
				let tbData = data;
				tbData.forEach( d => {
					d.update = `<a href="/tenda/picklist/update?PID=${d.PID}" target="_blank">update</a>`;
				})
				return {
					"titles": titles,
					"data": tbData,
					"names": objNames,
				}
			}
			if(data[0].LID) {
				let titles = ["Location", "Model","Unit","Last Total","Total Picks", "Cartons", "Boxes", "completed_at", "Action"];
				let objNames =   ["location", "model","unit","old_total","total_picks","cartons", "boxes", "completed_at", "update"];
				let tbData = data;
				tbData.forEach( d => {
					d.update = `<a href="/tenda/lastupdated?LID=${d.LID}" target="_blank">update</a>`;
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
			dbtable_container.appendChild(newtable);
			$("#dbtable").DataTable({
				dom: 'Bfrtip',
				buttons: ["excel"],
				// order: [5, "des"],
			});
			let table_width = rebuild_dbtable();
			return table_width;
		})
		.then((tw)=>{
			let cbtn1 = document.querySelector("#completeBtn")
			if(!cbtn1 && pickStatus === "Pending") {
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

				data.append("pickDate", pickDate)
				data.append("pickStatus", pickStatus)
				data.append("lastSaturday", lastSaturdayTS());

				let tableBody = document.querySelector('table tbody')
				let trows = tableBody.querySelectorAll('tr');
				if (tableBody && trows.length) {
					fetch("https://gzhang.dev/tenda/api/picklist/complete", {
						method: "POST",
						body: data,
					})
					.then(resp => { 
						return resp.json();
					})
					.then(data => {
						if(data) {
							const modal = document.querySelector(".modal");
							modal.classList.toggle("show-modal");
						}
						let model = document.querySelector(".model");
						let location = document.querySelector(".location");
						let sqlinfo = document.querySelector(".sqlinfo");
						let oldTotal = document.querySelector(".oldTotal");
						let pickQty = document.querySelector(".pickQty");
						let unit = document.querySelector(".unit");
						let newCartons = document.querySelector(".newCartons");
						let newBoxes = document.querySelector(".newBoxes");
						let newTotal = document.querySelector(".newTotal");

						data.forEach( (d,i) => {
							// let completeInfoTitle = ['Location',"Unit",'OLD total','Picked','NEW cartons', 'NEW boxes', 'NEW total'];
							// let completeInfoOrders = ['location','unit','oldTotal','pickQty','newCartons', 'newBoxes', 'newTotal']

							let completeInfoTitle = ['Location','Picked','NEW cartons', 'NEW boxes', 'NEW total'];
							let completeInfoOrders = ['location','pickQty','newCartons', 'newBoxes', 'newTotal']
							let completeData = []
							d.rowTitle = `${d.sqlinfo.split(' ')[0]} last_updated`
							completeData.push(d);

							let completeInfoTable = createTable(completeInfoTitle,completeData,completeInfoOrders)
							
							let calcTotal = d.oldTotal - d.pickQty
							let calcCartons = Math.trunc(calcTotal/d.unit)
							let calcBoxes = ((calcTotal/d.unit) % 1).toFixed(2) * d.unit
							calcBoxes = parseInt(calcBoxes)
							let ciTbody = completeInfoTable.tBodies[0]
							let ciRow = ciTbody.insertRow(1);

							let ciCell0 = ciRow.insertCell(0);
							ciCell0.textContent = d.location;

							let ciCell1 = ciRow.insertCell(1);
							ciCell1.textContent = "";

							let ciCell2 = ciRow.insertCell(2);
							ciCell2.textContent = "";


							let ciCell3 = ciRow.insertCell(3);
							ciCell3.textContent = "";

							let ciCell4 = ciRow.insertCell(4);
							ciCell4.textContent = calcCartons;

							let ciCell5 = ciRow.insertCell(5);
							ciCell5.textContent = calcBoxes;

							let ciCell6 = ciRow.insertCell(6);
							ciCell6.textContent = calcTotal;

							let completefb = document.createElement("div");
							let completefb_title = document.createElement("h3");
							completefb_title.textContent = `Model: ${d.model}`;
							completefb.appendChild(completefb_title);
							completefb.classList.add("complete-fd");
							completefb.appendChild(completeInfoTable);
							let cmpinfo = document.querySelector('#complete-info');
							cmpinfo.appendChild(completefb);
						})
					})
				}
			})
		})
		.catch(err => {
			dbtable_container.innerHTML = `<div class="alert alert-info">${err.name} : ${err.message}`;
			let alertInfo = document.querySelector("#dbtable_container div.alert-info")
			alertInfo.style.width  = "500px";
			fadeOut(alertInfo)
		})
	}
});

let closeBtn = document.querySelector(".btn-cancel");
if(closeBtn) {
	closeBtn.addEventListener('click', function() {
		let rburl = `https://gzhang.dev/tenda/api/txrb?rbn=CompletePickList`
		fetch(rburl)
		.then(resp => {return resp.json()})
		.then(data =>{
			console.log("resText[err]:", data);
			const modal = document.querySelector(".modal");
			modal.classList.toggle("show-modal");
		})
		.then(()=>{
			let cmpinfo = document.querySelector('#complete-info');
			cmpinfo.innerHTML = ""
		})
	});
}

let commitBtn = document.querySelector(".btn-commit");
if(commitBtn) {
	commitBtn.addEventListener('click', function() {
		//CompletePickList
		let cmurl = `https://gzhang.dev/tenda/api/txcm?cmn=CompletePickList`
		fetch(cmurl)
		.then(resp => {return resp.json()})
		.then(data =>{
			console.log("resText[err]:", data);
			const modal = document.querySelector(".modal");
			modal.classList.toggle("show-modal");
		})
		.then(()=>{
			let cmpinfo = document.querySelector('#complete-info');
			cmpinfo.innerHTML = ""
		})
	})
}



