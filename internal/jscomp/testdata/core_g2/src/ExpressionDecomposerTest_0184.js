let temp_const$jscomp$0;
let temp_const$jscomp$1;
var temp_const$jscomp$3 =
    (temp_const$jscomp$0 = x) == null
        ? void 0 : (temp_const$jscomp$1 = temp_const$jscomp$0.y).z;
// This double .call is unfortunate but not incorrect.
// Maybe we should make the ExpressionDecomposer recognize `.call` and avoid creating
// another one? We'd run the risk of breaking "real" methods called `.call`, which
// are allowed.
var temp_const$jscomp$2 = temp_const$jscomp$3.call;
// The temp_const$jscomp$1 argument gets type '?' when we decompose,
// which is correct, but TypeInference on this expected code appears to give it
// `undefined`. This is why we've set `shouldTestTypes = false;` above
a = temp_const$jscomp$2.call(temp_const$jscomp$3, temp_const$jscomp$1, foo());
