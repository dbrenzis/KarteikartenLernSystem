document.getElementById("selectOpt").addEventListener("change", post);

function post(method = "post") {
  var e = document.getElementById("selectOpt");
  var unter = e.options[e.selectedIndex].value;
  var ueber = e.selectedIndex;
  if (ueber <= 9) {
    ueber = "Naturwissenschaften";
  } else if (ueber <= 18) {
    ueber = "Sprachen";
  } else if (ueber <= 28) {
    ueber = "Gesellschaft";
  } else if (ueber <= 34) {
    ueber = "Wirtschaft";
  } else {
    ueber = "Geisteswissenschaften";
  }

  const form = document.createElement("form");
  form.method = method;
  form.action = "/karteikasten/filtered";

  const hiddenField = document.createElement("input");
  hiddenField.type = "text";
  hiddenField.name = "UnterKate";
  hiddenField.value = unter;

  form.appendChild(hiddenField);

  const hiddenField2 = document.createElement("input");
  hiddenField2.type = "text";
  hiddenField2.name = "UeberKate";
  hiddenField2.value = ueber;

  form.appendChild(hiddenField2);

  document.body.appendChild(form);
  form.submit();
}
