function requestHaTs(cookieName, triggerId, bucketSample, promptSample) {
  var inBucket;

  var cookies = decodeURIComponent(document.cookie).split(';');

  for (let i = 0; i < cookies.length; i++) {
    var c = cookies[i];

    while (c.charAt(0) == ' ') c = c.substring(1);

    if (c.indexOf(cookieName + '=') == 0)
      inBucket = c.substring((cookieName + '=').length, c.length);
  }

  if (typeof inBucket === 'undefined') {
    inBucket = String(Math.random() < bucketSample);
    document.cookie =
      cookieName + '=' + inBucket + '; path=/; max-age=2592000;';
  }

  if (inBucket === 'true') {
    var shouldPrompt = Math.random() < promptSample;
    if (shouldPrompt) {
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
          triggerId: triggerId,
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
  }
}

(function () {
  // HaTS - go.dev
  if (location.pathname !== '/play/') {
    requestHaTs('HaTS_BKT', 'RLVVv5Lf10njVvnD1rP0QUpmtosS', 0.03, 0.5);
  }

  // All download links on the go.dev homepage may trigger a survey.
  [
    document.querySelector('.js-downloadBtn'),
    document.querySelector('.js-downloadWin'),
    document.querySelector('.js-downloadMac'),
    document.querySelector('.js-downloadLinux'),
  ].forEach(function (el) {
    if (el) {
      el.addEventListener('click', () => {
        // HaTS - Core Go distribution
        requestHaTs('HaTS_BKT_DIST', 'dz6fkRxyz0njVvnD1rP0QxCXzhSX', 0.2, 1);
      });
    }
  });
})();
