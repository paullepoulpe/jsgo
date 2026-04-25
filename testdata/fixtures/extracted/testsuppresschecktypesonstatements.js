/** @param {string} s */
function takesString(s) {}

takesString(0); // verify this is an error
/** @suppress {checkTypes} */
takesString(null);
/** @suppress {checkTypes} */
(0, takesString(undefined));
