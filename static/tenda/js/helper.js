function getCurrentDateTime(numonly) {
	let date = new Date();
	let dates = date.toString().split(" ");
	let currTime = dates[4];
	var currentDateTime = ""
	var month = date.getMonth()+1;
	if(month < 10) {
		month = `0${month}`
	}
	let day = date.getDate();
	if (day < 10) {
		day = `0${day}`
	}
	if (numonly) {
		currentDateTime = `${date.getFullYear()}${month}${day}`;
	}
	else {
		currentDateTime = `${date.getFullYear()}-${month}-${day} ${currTime}`;
	}
	return currentDateTime;
}

function fetchDataList(url, idName, afterElem) {
	fetch(url)
	.then( res => { res.json() })
	.then(data => {
		// console.log("data->",data);
		var dataList = [];
		for(let [key, value] of Object.entries(data)) {
			dataList.push(value);
		}
		return dataList;
	})
	.then(dlist => {
		let datalistElem = document.createElement('datalist');
		datalistElem.id = idName;
		let parentElement = document.querySelector(afterElem);
		parentElement.insertAdjacentElement('afterend', datalistElem);
		let optionFragement = new DocumentFragment();
		dlist.forEach( dt => {
			let currentOpt = '<option value="'+dt+'">';
			let opt = document.createElement('option');
			opt.value = dt;
			optionFragement.appendChild(opt);
		});
		datalistElem.appendChild(optionFragement);
	})
	.catch( err => {
		console.log("datalist fetching error:", err);
	})
}

function WhenClick(elem, listener) {
	let theEl = document.querySelector(elem);
	if(theEl) {
		theEl.addEventListener('click', listener);
	}
}

function createTable(tableTitles, data, objNames) {
	let table = document.createElement("table");
	let theader = table.createTHead();
	let tbody = table.createTBody();
	let trow = theader.insertRow(0);

	tableTitles.forEach( (tht,idx) =>{
		let cell = trow.insertCell(idx);
		cell.innerHTML = tht;
	});

	let cells = tableTitles.length
	data.forEach( (d, i)=> {
		let tbrow = tbody.insertRow(i);
		for(let i = 0; i < cells; i++) {
			let cell = tbrow.insertCell(i);
			console.log("[helper.js] set cell value:",d[objNames[i]])
			cell.innerHTML = d[objNames[i]]
			cell.setAttribute("data-label", tableTitles[i])
		}
	});
	return table;
}

function rebuild_dbtable() {
	let dbtable = document.querySelector("#dbtable")
	let dbtable_wrapper = document.querySelector("#dbtable_wrapper");
	// dbtable width
	let dbtable_width = dbtable.offsetWidth;

	// set pagination width = dbtable width
	let paginate = document.querySelector("div#dbtable_paginate");
	// paginate.style.width = `${dbtable_width}px`;
	paginate.style.width = `520px`;

	// wrap dt-buttons and dt-filter with `div.dbtable_header`
	let dbtable_header = document.createElement("div");
	dbtable_header.id = "dbtable_header";
	let dtbtns = document.querySelector("#dbtable_wrapper .dt-buttons");
	let dtfilter = document.querySelector("#dbtable_filter");
	dbtable_header.appendChild(dtbtns);
	dbtable_header.appendChild(dtfilter);
	// set dbtable_header = dbtable width
	dbtable_header.style.width = `${dbtable_width}px`;
	
	// append dbtable_header to dbtable_wrapper
	dbtable_wrapper.insertAdjacentElement('afterbegin', dbtable_header)

	let dtBtn = document.querySelector("#dbtable_wrapper .dt-button");
	dtBtn.classList.add("btn", "btn-table")

	// set searchBar placeholder and styling
	let searchBar = document.querySelector("#dbtable_filter input[type='search']");
	searchBar.setAttribute("placeholder", "Search...")
	searchBar.style.marginLeft = "2em"
	return dbtable_width
}

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
	return lastSaturday;
}

function lastSunTS() {
	const t = new Date().getDate() + (6 - new Date().getDay() - 1) - 5 ;
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
	return lastSaturday;
}

function formDataCollect(formElement) {
	let form = document.querySelector(formElement)
	let formData = new FormData(form);
	let data = new URLSearchParams(formData);
	let timenow = getCurrentDateTime();
	for(const pair of formData) {
		if(!pair[1]) {
			let inputName = document.querySelector(`input[name=${pair[0]}]`)
			inputName.style.border="1px solid red"
			return 
		}
		console.log(`${pair[0]} is ${pair[1]}`)
		data.append(pair[0], pair[1])
	}
	return data
}

function fadeOut(elem, time=1500) {
	elem.style.opacity  = 1;
	setTimeout(()=>{
		elem.style.opacity  = 0;
	}, time);
}

function endTransact(txCname, txRname, cfmTbl, feedbackElem) {
	let txcm = document.querySelector(txCname);
	let txrb = document.querySelector(txRname);
	let txname = txcm.getAttribute("data-txname");
	let cmurl = `https://gzhang.dev/tenda/api/txcm?cmn=${txname}`
	let rburl = `https://gzhang.dev/tenda/api/txrb?rbn=${txname}`
	txcm.addEventListener("click",function(e) {
		e.preventDefault();
		fetch(cmurl)
		.then(res=>{return res.json()})
		.then(data=> {
			if(data.err === "") {
				cfmTbl.remove();
				feedbackElem.classList.add("alert-success");
				feedbackElem.innerHTML = `update successfully`;
				fadeOut(feedbackElem)
			}
			if(data.err !== "") {
				cfmTbl.remove();
				feedbackElem.classList.add("alert-danger");
				feedbackElem.innerHTML = `Error updating: ${err}`;
				fadeOut(feedbackElem)
			}
		});
	});
	txrb.addEventListener("click", function(e) {
		e.preventDefault();
		fetch(rburl)
		.then(res=>{return res.json()})
		.then(data=> {
			if(data.err === "") {
				cfmTbl.remove();
				feedbackElem.innerHTML = `discard successfully`;
				feedbackElem.classList.add("alert-secondary");
				fadeOut(feedbackElem)
			}
			if(data.err !== "") {
				cfmTbl.remove();
				feedbackElem.innerHTML = `Error updating: ${err}`;
				feedbackElem.classList.add("alert-danger");
				fadeOut(feedbackElem)
			}
		});
	});
}

const Tenda = {
	"API": {
		"models": "https://gzhang.dev/tenda/api/model",
		"locations":"https://gzhang.dev/tenda/api/locations",
	}
};

export {
	getCurrentDateTime,
	fetchDataList,
	WhenClick,
	createTable,
	rebuild_dbtable,
	fadeOut,
	lastSaturdayTS,
	lastSunTS,
	formDataCollect,
	endTransact,
	Tenda,
};