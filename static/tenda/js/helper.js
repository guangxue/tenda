function getCurrentDateTime(numonly) {
	let date = new Date();
	let dates = date.toString().split(" ");
	let currTime = dates[4];
	var currentDateTime = ""
	var month = date.getMonth()+1;
	if(month < 10) {
		month = `0${month}`
	}
	if (numonly) {
		currentDateTime = `${date.getFullYear()}${month}${date.getDate()}`;
	}
	else {
		currentDateTime = `${date.getFullYear()}-${month}-${date.getDate()} ${currTime}`;
	}
	console.log(currentDateTime)
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

function fetchModelData(argument) {
	// body...
}


function WhenClick(elem, listener) {
	let theEl = document.querySelector(elem);
	if(theEl) {
		theEl.addEventListner('click', listener);
	}
}

export default {getCurrentDateTime, fetchDataList, WhenClick};