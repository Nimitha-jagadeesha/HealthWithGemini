let form1 = document.querySelector('.form1');
let output1 = document.querySelector('.output1');
let output2 = document.querySelector('.output2');
let form2 = document.querySelector('.form2');
form1.onsubmit = async (ev) => {
    ev.preventDefault()
    output1.textContent = "Generating...";

    var data = new FormData(form1);
    var request = new XMLHttpRequest();
    request.open("POST", "/api/generate");
    request.onload = function () {
        // Read the response and interpret the output as markdown.
        let md = window.markdownit();
        output1.innerHTML = md.render(request.responseText);
    };
    request.send(data);
    return false;
}
form2.onsubmit = async (ev) => {
  ev.preventDefault()
  output2.textContent = "Generating...";

  var data = new FormData(form2);
  var request = new XMLHttpRequest();
  request.open("POST", "/api/generate");
  request.onload = function () {
      // Read the response and interpret the output as markdown.
      let md = window.markdownit();
      output2.innerHTML = md.render(request.responseText);
  };
  request.send(data);
  return false;
}
function readURL(input, val) {
    if (input.files && input.files[0]) {
      var reader = new FileReader();
  
      reader.onload = function (e) {
        var img;
        if(val === 1)
          img = document.querySelector("#img1");
        else
          img = document.querySelector("#img2");
        img.src = reader.result;
        img.style.display="inline";
      };
  
      reader.readAsDataURL(input.files[0]);
    }
  }