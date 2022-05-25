// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

// for /play; play.js is for embedded play widgets

window.addEventListener('load', () => {
  // Set up playground if enabled.
  if (window.playground) {
    window.playground({
      "codeEl":        ".js-playgroundCodeEl",
      "outputEl":      ".js-playgroundOutputEl",
      "runEl":         ".js-playgroundRunEl",
      "fmtEl":         ".js-playgroundFmtEl",
      "shareEl":       ".js-playgroundShareEl",
      "shareURLEl":    ".js-playgroundShareURLEl",
      "toysEl":        ".js-playgroundToysEl",
      "versionEl":     ".js-playgroundVersionEl",
      'enableHistory': true,
      'enableShortcuts': true,
      'enableVet': true
    });

    // The pre matched below is added by the code above. Style it appropriately.
    document.querySelector(".js-playgroundOutputEl pre").classList.add("Playground-output");

    $('#code').linedtextarea();
    $('#code').attr('wrap', 'off');
    $('#code').resize(function() { $('#code').linedtextarea(); });
  }
});
