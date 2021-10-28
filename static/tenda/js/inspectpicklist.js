import {createTable } from "./helper.js"
console.log("[ inspectpicklist.js ]")

let qtys = document.querySelectorAll(".qty")
var qtyes = Array.from(qtys)
let arrqty = []
qtyes.forEach( (x,i)=> {
	let currNo = parseInt(x.textContent, 10);
	arrqty.push(currNo);
});
const sum = arrqty.reduce(add, 0);

function add(accumulator, a) {
	return accumulator + a;
}

let total = document.querySelector("#total");
total.textContent = sum;

let PNOinfo = document.querySelectorAll(".PNOinfo");
let PNOinfos = Array.from(PNOinfo);

PNOinfos.forEach( (elm, i)=> {
	elm.addEventListener("click", function() {
		const PNO = elm.textContent;
		fetch("https://gzhang.dev/tenda/api/picklist/PNO/"+PNO)
		.then( resp => {
			return resp.json();
		})
		.then( data => {
			console.log(data)
			let titles = ['PID', 'PNO', 'Customer', 'Sales', 'Location', 'Model', 'Quantity','Created_at','--------']
			let objnames = ['PID','PNO','customer','sales_mgr','location','model','qty','created_at', ]; 
            let table = createTable(titles,data,objnames);
            let outb = document.querySelector("#outb");
            outb.innerHTML = "";
            outb.appendChild(table);
		})
	})
})