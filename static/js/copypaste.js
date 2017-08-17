var copyButton = document.querySelector('.copy-button');

copyButton.addEventListener('click', function(event) {
  var copyTextarea = document.querySelector('.text');
  copyTextarea.select();

  try {
    var successful = document.execCommand('copy');
    var msg = successful ? 'successful' : 'unsuccessful';
    console.log('copied');
  } catch (err) {
    alert("Browser does not support copy and paste automation :(")
    console.log('failed copy');
  }
});