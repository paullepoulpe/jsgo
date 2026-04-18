/**
 * @this {Array<T>|{length:number}}
 * @return {T}
 * @template T
 */
Array.prototype.method = function() {};
(function(){
  Array.prototype.method.call(arguments);
})();
