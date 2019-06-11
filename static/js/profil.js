document.getElementById("profilDelButton").addEventListener("click", modalset);
document.getElementById("alertCancel").addEventListener("click", modaldel);

let modal = document.getElementById("myModal");

function modalset(e) {
  modal.style.display = "block";
  e.preventDefault();
}

function modaldel(e) {
  modal.style.display = "none";
  e.preventDefault();
}
