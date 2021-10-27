import {createTable } from "./helper.js"
console.log("--- inspectpicklist")

let qtys = document.querySelectorAll(".qty")
var qtyes = Array.from(qtys)
console.log(qtyes)


let arrqty = []

qtyes.forEach( (x,i)=> {
	let currNo = parseInt(x.textContent, 10);
	arrqty.push(currNo);
});

console.log(arrqty)

const sum = arrqty.reduce(add, 0);

function add(accumulator, a) {
	return accumulator + a;
}

let total = document.querySelector("#total");
total.textContent = sum;

// /tenda/api/picklist/PNO/{{.PNO}}
let PNOinfo = document.querySelectorAll(".PNOinfo");

let PNOinfos = Array.from(PNOinfo);

PNOinfos.forEach( (elm, i)=> {
	console.log("elem:", elm)
	elm.addEventListener("click", function() {
		console.log("elem PNO:", elm.textContent)
		const PNO = elm.textContent;
		fetch("https://gzhang.dev/tenda/api/picklist/PNO/"+PNO)
		.then( resp => {
			return resp.json();
		})
		.then( data => {
			console.log(data)
			let titles = ['PID','PNO','customer','sales_mgr','location','model','qty','created_at'
			]; 
            let table = createTable(titles,data,titles);
            let outb = document.querySelector("#outb");
            outb.appendChild(table)
		})
	})
})