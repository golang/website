window.transport = {{.Transport}}();
window.socketAddr = "{{.SocketAddr}}";

function highlight(selector) {
    var speed = 50;
    var obj = $(selector).stop(true, true)
    for (var i = 0; i < 5; i++) {
        obj.addClass("highlight", speed)
        obj.delay(speed)
        obj.removeClass("highlight", speed)
    }
}

function highlightAndClick(selector) {
    highlight(selector);
    setTimeout(function() {
        $(selector)[0].click()
    }, 750);
}

function click(selector) {
    $(selector)[0].click();
}

/**
 * toggleTheme switches the preferred color scheme between auto, light, and dark.
 */
function toggleTheme() {
    let nextTheme = 'dark';
    const theme = document.documentElement.getAttribute('data-theme');
    if (theme === 'dark') {
      nextTheme = 'light';
    } else if (theme === 'light') {
      nextTheme = 'auto';
    }
    let domain = '';
    if (location.hostname === 'go.dev') {
      // Include subdomains to apply the setting to pkg.go.dev.
      domain = 'domain=.go.dev;';
    }
    document.documentElement.setAttribute('data-theme', nextTheme);
    document.cookie =
    `prefers-color-scheme=${nextTheme};${domain}path=/;max-age=31536000;`;
  }


function setThemeButtons() {
    for (const el of document.querySelectorAll('.js-toggleTheme')) {
      el.addEventListener('click', () => {
        toggleTheme();
      });
    }
}

setThemeButtons();
