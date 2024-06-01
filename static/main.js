let form = document.querySelector('form');
let output = document.querySelector('.output');

form.onsubmit = async (ev) => {
    ev.preventDefault()
    output.textContent = "Generating...";

    var data = new FormData(form);
    var request = new XMLHttpRequest();
    request.open("POST", "/api/generate");
    request.onload = function () {
        // Read the response and interpret the output as markdown.
        let md = window.markdownit();
        output.innerHTML = md.render(request.responseText);
    };
    request.send(data);
    return false;
}
function readURL(input) {
    if (input.files && input.files[0]) {
      var reader = new FileReader();
  
      reader.onload = function (e) {
        var img = document.querySelector("#img");
        img.src = reader.result;
        img.style.display="inline";
      };
  
      reader.readAsDataURL(input.files[0]);
    }
  }
function hideImg() {
    document.getElementById("img")
                      .style.display = "none";
   }