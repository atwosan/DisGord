var texts = d.split("\n")
var text = ""
for (let i = 0; i < texts.length; i++) {
  text += "<p>" + texts[i] + "</p>" 
}
document.write(text);
