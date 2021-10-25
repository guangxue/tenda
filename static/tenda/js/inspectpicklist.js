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