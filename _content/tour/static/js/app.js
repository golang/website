/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

angular.module('tour', ['ui', 'tour.services', 'tour.controllers', 'tour.directives', 'tour.values', 'ng']).

config(['$routeProvider', '$locationProvider',
    function($routeProvider, $locationProvider) {
        $routeProvider.
        when('/tour/', {
            redirectTo: '/tour/welcome/1'
        }).
        when('/tour/list', {
            templateUrl: '/tour/static/partials/list.html',
        }).
        when('/tour/:lessonId/:pageNumber', {
            templateUrl: '/tour/static/partials/editor.html',
            controller: 'EditorCtrl'
        }).
        when('/tour/:lessonId', {
            redirectTo: '/tour/:lessonId/1'
        }).
        otherwise({
            redirectTo: '/tour/'
        });

        $locationProvider.html5Mode(true).hashPrefix('!');
    }
]).

// handle mapping from old paths (#42) to the new organization.
run(function($rootScope, $location, mapping) {
    $rootScope.$on( "$locationChangeStart", function(event, next) {
        var url = document.createElement('a');
        url.href = next;
        if (url.pathname == '/') {
            window.location.href = next;
            return;
        }
        if (url.pathname != '/tour/' || url.hash == '') {
            return;
        }
        $location.hash('');
        var m = mapping[url.hash];
        if (m === undefined) {
            console.log('unknown url, redirecting home');
            $location.path('/tour/welcome/1');
            return;
        }
        $location.path('/tour' + m);
    });
});
