// Copyright 2021 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

window.addEventListener('DOMContentLoaded', () => {
  // Set up playground if enabled.
  if (window.playground) {
    window.playground({
      "codeEl":        ".js-playgroundCodeEl",
      "outputEl":      ".js-playgroundOutputEl",
      "runEl":         ".js-playgroundRunEl",
      "fmtEl":         ".js-playgroundFmtEl",
      "shareEl":       ".js-playgroundShareEl",
      "shareRedirect": "/play/p/",
      "toysEl":        ".js-playgroundToysEl"
    });

    // The pre matched below is added by the code above. Style it appropriately.
    document.querySelector(".js-playgroundOutputEl pre").classList.add("Playground-output");
    $('.js-playgroundToysEl').val("hello.go").trigger("change")
  } else {
    $(".Playground").hide();
  }
});
