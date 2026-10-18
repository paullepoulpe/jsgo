/** @const */ var STORAGE_KEY = 'small-app-count';
/** @return {number} */
function loadCount() {
  var raw = window.localStorage.getItem(STORAGE_KEY);
  var count = raw === null ? 0 : Number(raw);
  return isNaN(count) ? 0 : count;
}
/** @param {!Element} button */
function installCounter(button) {
  var count = loadCount();
  button.textContent = String(count);
  button.addEventListener('click', function() {
    count++;
    button.textContent = String(count);
    window.localStorage.setItem(STORAGE_KEY, String(count));
  });
}
var button = document.getElementById('counter');
if (button) installCounter(button);
