// The input on which behave.py checks that a jsgo binary has the seeded faults
// (`go build -tags jsgo_faults`): with JSGO_FAST_FAULT=fold-negzero and
// --fast_mode, the -0 below is printed as 0.
var z = 0 * -1;
window["z"] = 1 / z;
