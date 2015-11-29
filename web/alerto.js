var alerto = angular.module('alerto', ['ngResource']);

var Alerto = {};

/**
 * Object to hold all angular controllers
 */
Alerto.Controller = {};

/**
 * Object to hold all angular directives
 */
Alerto.Directive = {};

/**
 * Object to hold all angular services
 */
Alerto.Service = {};

/**
 * Object to hold all angular factories
 */
Alerto.Factory = {};

/**
 * @constructor
 * @param {*} $resource
 * @ngInject
 * @suppress {checkTypes}
 */
Alerto.Factory.MonitorService = function($resource) {
    return $resource('/monitor/:id', {id: '@id'});
};

alerto.factory('MonitorService', Alerto.Factory.MonitorService);

/**
 * @ngInject
 * @constructor
 */
Alerto.Controller.MainController = function(MonitorService, $http) {
	var self = this;

	this.monitors = MonitorService.query();

	this.agents = {};
	$http.get('/agent/').then(function(response) {
		self.agents = response.data;
	});

	/**
	 * @expose
	 */
	this.deleteMonitor = function(id) {
		MonitorService.delete({id: id});
	};
};

alerto.controller('MainController', Alerto.Controller.MainController);
