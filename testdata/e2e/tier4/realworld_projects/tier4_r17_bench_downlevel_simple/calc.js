const sumAll = (...args) => args.reduce((acc, x) => acc + x, 0);
class Geometry { static area(w, h) { return sumAll(w * h); } }
window['area'] = Geometry.area(10, 20);
