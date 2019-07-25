// Copyright 2012 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/* A little code to ease navigation of these documents.
 *
 * On window load we:
 *  + Generate a table of contents (generateTOC)
 *  + Bind foldable sections (bindToggles)
 *  + Bind links to foldable sections (bindToggleLinks)
 */

(function() {
'use strict';

var headerEl = document.querySelector('.js-header');
var menuButtonEl = document.querySelector('.js-headerMenuButton');
menuButtonEl.addEventListener('click', function(e) {
  e.preventDefault();
  headerEl.classList.toggle('is-active');
  menuButtonEl.setAttribute('aria-expanded', headerEl.classList.contains('is-active'));
});

/* Generates a table of contents: looks for h2 and h3 elements and generates
 * links. "Decorates" the element with id=="nav" with this table of contents.
 */
function generateTOC() {
  if ($('#manual-nav').length > 0) {
    return;
  }

  // For search, we send the toc precomputed from server-side.
  // TODO: Ideally, this should always be precomputed for all pages, but then
  // we need to do HTML parsing on the server-side.
  if (location.pathname === '/search') {
    return;
  }

  var nav = $('#nav');
  if (nav.length === 0) {
    return;
  }

  var toc_items = [];
  $(nav).nextAll('h2, h3').each(function() {
    var node = this;
    if (node.id == '')
      node.id = 'tmp_' + toc_items.length;
    var link = $('<a/>').attr('href', '#' + node.id).text($(node).text());
    var item;
    if ($(node).is('h2')) {
      item = $('<dt/>');
    } else { // h3
      item = $('<dd class="indent"/>');
    }
    item.append(link);
    toc_items.push(item);
  });
  if (toc_items.length <= 1) {
    return;
  }
  var dl1 = $('<dl/>');
  var dl2 = $('<dl/>');

  var split_index = (toc_items.length / 2) + 1;
  if (split_index < 8) {
    split_index = toc_items.length;
  }
  for (var i = 0; i < split_index; i++) {
    dl1.append(toc_items[i]);
  }
  for (/* keep using i */; i < toc_items.length; i++) {
    dl2.append(toc_items[i]);
  }

  var tocTable = $('<table class="unruled"/>').appendTo(nav);
  var tocBody = $('<tbody/>').appendTo(tocTable);
  var tocRow = $('<tr/>').appendTo(tocBody);

  // 1st column
  $('<td class="first"/>').appendTo(tocRow).append(dl1);
  // 2nd column
  $('<td/>').appendTo(tocRow).append(dl2);
}

function bindToggle(el) {
  $('.toggleButton', el).click(function() {
    if ($(this).closest(".toggle, .toggleVisible")[0] != el) {
      // Only trigger the closest toggle header.
      return;
    }

    if ($(el).is('.toggle')) {
      $(el).addClass('toggleVisible').removeClass('toggle');
    } else {
      $(el).addClass('toggle').removeClass('toggleVisible');
    }
  });
}

function bindToggles(selector) {
  $(selector).each(function(i, el) {
    bindToggle(el);
  });
}

function bindToggleLink(el, prefix) {
  $(el).click(function() {
    var href = $(el).attr('href');
    var i = href.indexOf('#'+prefix);
    if (i < 0) {
      return;
    }
    var id = '#' + prefix + href.slice(i+1+prefix.length);
    if ($(id).is('.toggle')) {
      $(id).find('.toggleButton').first().click();
    }
  });
}
function bindToggleLinks(selector, prefix) {
  $(selector).each(function(i, el) {
    bindToggleLink(el, prefix);
  });
}

function setupInlinePlayground() {
	'use strict';
	// Set up playground when each element is toggled.
	$('div.play').each(function (i, el) {
		// Set up playground for this example.
		var setup = function() {
			var code = $('.code', el);
			playground({
				'codeEl':   code,
				'outputEl': $('.output', el),
				'runEl':    $('.run', el),
				'fmtEl':    $('.fmt', el),
				'shareEl':  $('.share', el),
				'shareRedirect': '//play.golang.org/p/'
			});

			// Make the code textarea resize to fit content.
			var resize = function() {
				code.height(0);
				var h = code[0].scrollHeight;
				code.height(h+20); // minimize bouncing.
				code.closest('.input').height(h);
			};
			code.on('keydown', resize);
			code.on('keyup', resize);
			code.keyup(); // resize now.
		};

		// If example already visible, set up playground now.
		if ($(el).is(':visible')) {
			setup();
			return;
		}

		// Otherwise, set up playground when example is expanded.
		var built = false;
		$(el).closest('.toggle').click(function() {
			// Only set up once.
			if (!built) {
				setup();
				built = true;
			}
		});
	});
}

// fixFocus tries to put focus to #page so that keyboard navigation works.
function fixFocus() {
  var page = $('#page');
  var topbar = $('div#topbar');
  page.css('outline', 0); // disable outline when focused
  page.attr('tabindex', -1); // and set tabindex so that it is focusable
  $(window).resize(function (evt) {
    // only focus page when the topbar is at fixed position (that is, it's in
    // front of page, and keyboard event will go to the former by default.)
    // by focusing page, keyboard event will go to page so that up/down arrow,
    // space, etc. will work as expected.
    if (topbar.css('position') == "fixed")
      page.focus();
  }).resize();
}

function toggleHash() {
  var id = window.location.hash.substring(1);
  // Open all of the toggles for a particular hash.
  var els = $(
    document.getElementById(id),
    $('a[name]').filter(function() {
      return $(this).attr('name') == id;
    })
  );

  while (els.length) {
    for (var i = 0; i < els.length; i++) {
      var el = $(els[i]);
      if (el.is('.toggle')) {
        el.find('.toggleButton').first().click();
      }
    }
    els = el.parent();
  }
}

function personalizeInstallInstructions() {
  var prefix = '?download=';
  var s = window.location.search;
  if (s.indexOf(prefix) != 0) {
    // No 'download' query string; detect "test" instructions from User Agent.
    if (navigator.platform.indexOf('Win') != -1) {
      $('.testUnix').hide();
      $('.testWindows').show();
    } else {
      $('.testUnix').show();
      $('.testWindows').hide();
    }
    return;
  }

  var filename = s.substr(prefix.length);
  var filenameRE = /^go1\.\d+(\.\d+)?([a-z0-9]+)?\.([a-z0-9]+)(-[a-z0-9]+)?(-osx10\.[68])?\.([a-z.]+)$/;
  var m = filenameRE.exec(filename);
  if (!m) {
    // Can't interpret file name; bail.
    return;
  }
  $('.downloadFilename').text(filename);
  $('.hideFromDownload').hide();

  var os = m[3];
  var ext = m[6];
  if (ext != 'tar.gz') {
    $('#tarballInstructions').hide();
  }
  if (os != 'darwin' || ext != 'pkg') {
    $('#darwinPackageInstructions').hide();
  }
  if (os != 'windows') {
    $('#windowsInstructions').hide();
    $('.testUnix').show();
    $('.testWindows').hide();
  } else {
    if (ext != 'msi') {
      $('#windowsInstallerInstructions').hide();
    }
    if (ext != 'zip') {
      $('#windowsZipInstructions').hide();
    }
    $('.testUnix').hide();
    $('.testWindows').show();
  }

  var download = "https://dl.google.com/go/" + filename;

  var message = $('<p class="downloading">'+
    'Your download should begin shortly. '+
    'If it does not, click <a>this link</a>.</p>');
  message.find('a').attr('href', download);
  message.insertAfter('#nav');

  window.location = download;
}

function updateVersionTags() {
  var v = window.goVersion;
  if (/^go[0-9.]+$/.test(v)) {
    $(".versionTag").empty().text(v);
    $(".whereTag").hide();
  }
}

function addPermalinks() {
  function addPermalink(source, parent) {
    var id = source.attr("id");
    if (id == "" || id.indexOf("tmp_") === 0) {
      // Auto-generated permalink.
      return;
    }
    if (parent.find("> .permalink").length) {
      // Already attached.
      return;
    }
    parent.append(" ").append($("<a class='permalink'>&#xb6;</a>").attr("href", "#" + id));
  }

  $("#page .container").find("h2[id], h3[id]").each(function() {
    var el = $(this);
    addPermalink(el, el);
  });

  $("#page .container").find("dl[id]").each(function() {
    var el = $(this);
    // Add the anchor to the "dt" element.
    addPermalink(el, el.find("> dt").first());
  });
}

$(".js-expandAll").click(function() {
  if ($(this).hasClass("collapsed")) {
    toggleExamples('toggle');
    $(this).text("(Collapse All)");
  } else {
    toggleExamples('toggleVisible');
    $(this).text("(Expand All)");
  }
  $(this).toggleClass("collapsed")
});

function toggleExamples(className) {
  // We need to explicitly iterate through divs starting with "example_"
  // to avoid toggling Overview and Index collapsibles.
  $("[id^='example_']").each(function() {
    // Check for state and click it only if required.
    if ($(this).hasClass(className)) {
      $(this).find('.toggleButton').first().click();
    }
  });
}

$(document).ready(function() {
  generateTOC();
  addPermalinks();
  bindToggles(".toggle");
  bindToggles(".toggleVisible");
  bindToggleLinks(".exampleLink", "example_");
  bindToggleLinks(".overviewLink", "");
  bindToggleLinks(".examplesLink", "");
  bindToggleLinks(".indexLink", "");
  setupInlinePlayground();
  fixFocus();
  toggleHash();
  personalizeInstallInstructions();
  updateVersionTags();

  // godoc.html defines window.initFuncs in the <head> tag, and root.html and
  // codewalk.js push their on-page-ready functions to the list.
  // We execute those functions here, to avoid loading jQuery until the page
  // content is loaded.
  for (var i = 0; i < window.initFuncs.length; i++) window.initFuncs[i]();
});

})();
