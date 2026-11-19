function pad(n) { return n < 10 ? '0' + n : '' + n; }
function formatDate(d) { return d.getFullYear() + '-' + pad(d.getMonth() + 1) + '-' + pad(d.getDate()); }
window['fmt'] = formatDate(new Date(2026, 10, 14));
