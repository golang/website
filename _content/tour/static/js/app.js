/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

angular.module('tour', ['ui', 'tour.services', 'tour.controllers', 'tour.directives', 'tour.values', 'ng']).

config(['$routeProvider', '$locationProvider',
    function($routeProvider, $locationProvider) {
        var notFound = {
            templateUrl: '/tour/static/partials/notfound.html',
        };

        // The server sets window.tourNotFound when the requested URL is not a
        // tour page, in which case it also responds with a 404 status. Show the
        // not found page at that URL, instead of redirecting elsewhere, so that
        // the address bar keeps matching the response.
        if (window.tourNotFound) {
            var path = window.location.pathname;
            try {
                // Routes are matched against the decoded path.
                path = decodeURIComponent(path);
            } catch (e) {
                // Leave the path as is if it is not validly encoded.
            }
            // Registered first so that it takes precedence over the
            // lesson routes below.
            $routeProvider.when(path, notFound);
        }

        $routeProvider.
            when('/tour/', {
                redirectTo: '/tour/welcome/1'
            }).
            when('/tour/list', {
                templateUrl: '/tour/static/partials/list.html',
            }).
            when('/tour/notfound', notFound).
            when('/tour/:lessonId/:pageNumber', {
                templateUrl: '/tour/static/partials/editor.html',
                controller: 'EditorCtrl'
            }).
            when('/tour/:lessonId', {
                redirectTo: '/tour/:lessonId/1'
            }).
            otherwise(notFound);

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
