/**
 * @license
 * Copyright 2021 The Go Authors. All rights reserved.
 * Use of this source code is governed by a BSD-style
 * license that can be found in the LICENSE file.
 */

(function () {
  var cookieName = 'HaTS_BKT_DIST';
  var inBucket;

  var cookies = decodeURIComponent(document.cookie).split(';');

  for (let i = 0; i < cookies.length; i++) {
    var c = cookies[i];

    while (c.charAt(0) == ' ') {
      c = c.substring(1);
    }

    if (c.indexOf(cookieName + '=') == 0) {
      inBucket = c.substring((cookieName + '=').length, c.length);
    }
  }

  if (typeof inBucket === 'undefined') {
    inBucket = String(Math.random() < 0.01);
    document.cookie =
      cookieName + '=' + inBucket + '; path=/; max-age=2592000;';
  }

  if (inBucket === 'true') {
    document.cookie = cookieName + '=false ; path=/; max-age=2592000;';

    var tag = document.createElement('script');
    tag.src =
      'https://www.gstatic.com/feedback/js/help/prod/service/lazy.min.js';
    tag.type = 'text/javascript';
    document.head.appendChild(tag);

    tag.onload = function () {
      var helpApi = window.help.service.Lazy.create(0, {
        apiKey: 'AIzaSyDfBPfajByU2G6RAjUf5sjkMSSLNTj7MMc',
        locale: 'en-US',
      });

      helpApi.requestSurvey({
        triggerId: 'dz6fkRxyz0njVvnD1rP0QxCXzhSX',
        callback: function (requestSurveyCallbackParam) {
          if (!requestSurveyCallbackParam.surveyData) {
            return;
          }
          helpApi.presentSurvey({
            surveyData: requestSurveyCallbackParam.surveyData,
            colorScheme: 1, // light
            customZIndex: 10000,
          });
        },
      });
    };
  }
})();
