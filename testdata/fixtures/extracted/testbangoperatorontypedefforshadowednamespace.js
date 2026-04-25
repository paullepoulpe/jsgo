/** @typedef {!service.Service} */
var serviceService;
/** @constructor @param {!service.Service} service */
var Controller = function(service) {
  /** @private {!serviceService} */
  this.service = service;
};
/** @const */
var service = {};
/** @constructor */
service.Service = function() {};
