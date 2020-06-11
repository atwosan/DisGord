var counter = function(str,seq){
  return str.split(seq).length;
}
for (let step = 0; step < counter(d, "\n"); step++) {
  d = d.replace("\n", "<br>");
}
document.write(d);
