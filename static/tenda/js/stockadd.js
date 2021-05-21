let stockAddBtn = document.querySelector("#stockAddBtn")

stockAddBtn.addEventListener('click', function() {
	fetch("https://gzhang.dev/tenda/api/stock/add")
});