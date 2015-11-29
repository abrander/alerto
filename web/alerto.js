var alerto = angular.module('alerto', ['ngResource', 'ui.bootstrap']);

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
	return $resource('/monitor/:id', {id: '@id'}, {
		save: {
			url: '/monitor/new',
			method: 'POST'
		}
	});
};

alerto.factory('MonitorService', Alerto.Factory.MonitorService);

/**
 * @ngInject
 * @constructor
 */
Alerto.Controller.MainController = function(MonitorService, $http, $uibModal, $scope) {
	var self = this;

	this.monitors = MonitorService.query();
	this.uptime = 0;

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

	/**
	 * @expose
	 * @param {string} agentId
	 */
	this.addMonitor = function(agentId) {
		var modalInstance = $uibModal.open({
			animation: true,
			templateUrl: 'newMonitorModalTemplate',
			controller: 'NewMonitorController',
			size: 'lg',
			resolve: {
				agent: function() {
					var agent = self.agents[agentId];
					agent.agentId = agentId;
					return agent;
					}
				}
			});
		modalInstance.result.then(function(result) {
			// Convert to seconds
			result.interval *= 1000000000;
			result.agent.agentId = agentId;
			MonitorService.save(result);
		});
	};

	// Construct websocket url
	var url = '';
	if (window.location.protocol.replace(/:/g, '') === 'https')
		url += 'wss://';
	else
		url += 'ws://';
	url += window.location.host + '/ws';

	var socket = new WebSocket(url);

	socket.onmessage = function(msg) {
		var message = JSON.parse(msg.data);
		$scope.$apply(function() {

			switch (message.type) {
				case 'status':
					self.uptime = message.payload.uptime;
					break;
			}
		});
	};
};

alerto.controller('MainController', Alerto.Controller.MainController);

/**
 * @ngInject
 * @constructor
 */
Alerto.Controller.NewMonitorController = function($scope, $uibModalInstance, agent) {
	$scope.agent = agent;
	$scope.newMonitor = {};

	$scope.ok = function() {
		$uibModalInstance.close($scope.newMonitor);
	};

	$scope.cancel = function() {
		$uibModalInstance.dismiss('cancel');
	};
};

alerto.controller('NewMonitorController', Alerto.Controller.NewMonitorController);

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
		if (value == undefined)
			return '-';

		var d = parseInt(value);

		if (d < Microsecond) {
			return d + "ns";
		}

		if (d < Millisecond) {
			return d / Microsecond + "µs";
		}

		if (d < Second) {
			return d / Millisecond + "ms";
		}

		if (d < Minute) {
			return d / Second + "s";
		}

		if (d < Hour) {
			return d / Minute + "m";
		}

		return d / Hour + "h";
	};

	return filter;
};

alerto.filter('goDuration', Alerto.Filter.GoDuration);
