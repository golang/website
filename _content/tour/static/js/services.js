/* Copyright 2012 The Go Authors.   All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */
'use strict';

/* Services */

angular.module('tour.services', []).

// Google Analytics
factory('analytics', ['$window',
    function(win) {
        var track = win.trackPageview || (function() {});
        return {
            trackView: track
        };
    }
]).

// Internationalization
factory('i18n', ['translation',
    function(translation) {
        return {
            l: function(key) {
                if (translation[key]) return translation[key];
                return '(no translation for ' + key + ')';
            }
        };
    }
]).

// Running code
factory('run', ['$window', 'editor',
    function(win, editor) {
        var writeInterceptor = function(writer, done) {
            return function(write) {
                if (write.Kind == 'stderr') {
                    var lines = write.Body.split('\n');
                    for (var i in lines) {
                        var match = lines[i].match(/.*\.go:([0-9]+): ([^\n]*)/);
                        if (match !== null) {
                            editor.highlight(match[1], match[2]);
                        }
                    }
                }
                writer(write);
                if (write.Kind == 'end') done();
            };
        };
        return function(code, output, options, done) {
            // We want to build tour snippets in module mode, so append
            // a default go.mod file when it is not already included in
            // the txtar archive.
            //
            // The exercises use golang.org/x/tour/{pic,reader,tree,wc}
            // packages, so include the golang.org/x/tour module in the
            // build list.
            const hasGoMod = code.indexOf('\n-- go.mod --\n') !== -1 || code.startsWith('-- go.mod --\n');
            if (!hasGoMod) {
                code += '\n' +
                    '-- go.mod --\n' +
                    'module example\n' +
                    'require golang.org/x/tour v0.0.0-20201207214521-004403599411\n' +
                    '-- go.sum --\n' +
                    'golang.org/x/crypto v0.0.0-20190308221718-c2843e01d9a2/go.mod h1:djNgcEr1/C05ACkg1iLfiJU5Ep61QUkGW8qpdssI0+w=\n' +
                    'golang.org/x/net v0.0.0-20190311183353-d8887717615a/go.mod h1:t9HGtf8HONx5eT2rtn7q6eTqICYqUVnKs3thJo3Qplg=\n' +
                    'golang.org/x/sys v0.0.0-20190215142949-d0b11bdaac8a/go.mod h1:STP8DvDyc/dI5b8T5hshtkjS+E42TnysNCUPdjciGhY=\n' +
                    'golang.org/x/text v0.3.0/go.mod h1:NqM8EUOU14njkJ3fqMW+pc6Ldnwhi/IjpwHt7yyuwOQ=\n' +
                    'golang.org/x/tools v0.0.0-20190312164927-7b79afddac43/go.mod h1:LCzVGOaR6xXOjkQ3onu1FJEFr0SW1gC7cKk1uF8kGRs=\n' +
                    'golang.org/x/tour v0.0.0-20201207214521-004403599411 h1:dJ4kVwSGlrLZXW6eo2IOer4Pm3wl2GIG4fytRziMgL8=\n' +
                    'golang.org/x/tour v0.0.0-20201207214521-004403599411/go.mod h1:qMugOFWX59KzC8Nx7f2uvXxKxAqJfi1J6ZUHAWKnrRA=\n';
            }

            // PlaygroundOutput is defined in playground.js which is prepended
            // to the generated script.js in gotour/tour.go.
            // The next line removes the jshint warning.
            // global PlaygroundOutput
            return win.transport.Run(code, writeInterceptor(PlaygroundOutput(output), done), options);
        };
    }
]).

// Formatting code
factory('fmt', ['$http',
    function($http) {
        return function(body, imports) {
            var params = $.param({
                'body': body,
                'imports': imports,
            });
            var headers = {
                'Content-Type': 'application/x-www-form-urlencoded'
            };
            return $http.post('/_/fmt', params, {
                headers: headers
            });
        };
    }
]).

// Local storage, persistent to page refreshing.
factory('storage', ['$window',
    function(win) {
        try {
            // This will raise an exception if cookies are disabled.
            win.localStorage = win.localStorage;
            return {
                get: function(key) {
                    return win.localStorage.getItem(key);
                },
                set: function(key, val) {
                    win.localStorage.setItem(key, val);
                }
            };
        } catch (e) {
            return {
                get: function() {
                    return null;
                },
                set: function() {}
            };
        }
    }
]).

// Editor context service, kept through the whole app.
factory('editor', ['$window', 'storage',
    function(win, storage) {
        var ctx = {
            imports: storage.get('imports') === 'true',
            toggleImports: function() {
                ctx.imports = !ctx.imports;
                storage.set('imports', ctx.imports);
            },
            syntax: storage.get('syntax') === 'true',
            toggleSyntax: function() {
                ctx.syntax = !ctx.syntax;
                storage.set('syntax', ctx.syntax);
                ctx.paint();
            },
            paint: function() {
                var mode = ctx.syntax && 'text/x-go' || 'text/x-go-comment';
                // Wait for codemirror to start.
                var set = function() {
                    if ($('.CodeMirror').length > 0) {
                        var cm = $('.CodeMirror')[0].CodeMirror;
                        if (cm.getOption('mode') == mode) {
                            cm.refresh();
                            return;
                        }
                        cm.setOption('mode', mode);
                    }
                    win.setTimeout(set, 10);
                };
                set();
            },
            highlight: function(line, message) {
                $('.CodeMirror-code > div:nth-child(' + line + ')')
                    .addClass('line-error').attr('title', message);
            },
            onChange: function() {
                $('.line-error').removeClass('line-error').attr('title', null);
            }
        };
        // Set in the window so the onChange function in the codemirror config
        // can call it.
        win.codeChanged = ctx.onChange;
        return ctx;
    }
]).

// Table of contents management and navigation
factory('toc', ['$http', '$q', '$log', 'tableOfContents', 'storage',
    function($http, $q, $log, tableOfContents, storage) {
        var modules = tableOfContents;

        var lessons = {};

        var prevLesson = function(id) {
            var mod = lessons[id].module;
            var idx = mod.lessons.indexOf(id);
            if (idx < 0) return '';
            if (idx > 0) return mod.lessons[idx - 1];

            idx = modules.indexOf(mod);
            if (idx <= 0) return '';
            mod = modules[idx - 1];
            return mod.lessons[mod.lessons.length - 1];
        };

        var nextLesson = function(id) {
            var mod = lessons[id].module;
            var idx = mod.lessons.indexOf(id);
            if (idx < 0) return '';
            if (idx + 1 < mod.lessons.length) return mod.lessons[idx + 1];

            idx = modules.indexOf(mod);
            if (idx < 0 || modules.length <= idx + 1) return '';
            mod = modules[idx + 1];
            return mod.lessons[0];
        };

        $http.get('/tour/lesson/').then(
            function(data) {
                lessons = data.data;
                for (var m = 0; m < modules.length; m++) {
                    var module = modules[m];
                    module.lesson = {};
                    for (var l = 0; l < modules[m].lessons.length; l++) {
                        var lessonName = module.lessons[l];
                        var lesson = lessons[lessonName];
                        lesson.module = module;
                        module.lesson[lessonName] = lesson;

                        // replace file contents with locally stored copies.
                        for (var p = 0; p < lesson.Pages.length; p++) {
                            var page = lesson.Pages[p];
                            for (var f = 0; f < page.Files.length; f++) {
                                page.Files[f].OrigContent = page.Files[f].Content;
                                var val = storage.get(page.Files[f].Hash);
                                if (val !== null) {
                                    page.Files[f].Content = val;
                                }
                            }
                        }
                    }
                }
                moduleQ.resolve(modules);
                lessonQ.resolve(lessons);
            },
            function(error) {
                $log.error('error loading lessons : ', error);
                moduleQ.reject(error);
                lessonQ.reject(error);
            }
        );

        var moduleQ = $q.defer();
        var lessonQ = $q.defer();

        return {
            modules: moduleQ.promise,
            lessons: lessonQ.promise,
            prevLesson: prevLesson,
            nextLesson: nextLesson
        };
    }
]);
