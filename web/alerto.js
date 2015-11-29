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
 * Object to hold all angular factories
 */
Alerto.Factory = {};

/**
 * Object to hold all angular filters
 */
Alerto.Filter = {};

/**
 * Object to hold all angular services
 */
Alerto.Service = {};

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

/**
 * @ngInject
 * @constructor
 */
Alerto.Filter.GoDuration = function() {
	var Nanosecond = 1;
	var Microsecond = 1000 * Nanosecond;
	var Millisecond = 1000 * Microsecond;
	var Second = 1000 * Millisecond;
	var Minute = 60 * Second;
	var Hour = 60 * Minute;

	/**
	 * @param {string} value
	 * @return {string}
	 */
	var filter = function(value) {
		console.dir(value);
		if (value == undefined)
			return "-";

		var d = parseInt(value);

		if (d < Microsecond) {
			return d + 'ns';
		}

		if (d < Millisecond) {
			return d / Microsecond + 'µs';
		}

		if (d < Second) {
			return d / Millisecond + "ms";
		}

		if (d < Minute) {
			return d/Second + "s";
		}

		if (d < Hour) {
			return d/Minute + "m";
		}

		return d/Hour + "h";
	};

	return filter;
};

alerto.filter('goDuration', Alerto.Filter.GoDuration);
