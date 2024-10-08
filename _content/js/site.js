// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/**
 * A bit of navigation related code for handling dismissible elements.
 */
window.initFuncs = [];

(() => {
  'use strict';

  function registerHeaderListeners() {
    // Desktop menu hover state
    const menuItemHovers = document.querySelectorAll('.js-desktop-menu-hover');
    menuItemHovers.forEach(menuItemHover => {
      // when user clicks on the dropdown menu item on desktop or mobile,
      // force the menu to stay open until the user clicks off of it.
      menuItemHover.addEventListener('mouseenter', e => {
        const forced = document.querySelector('.forced-open');
        if (forced && forced !== menuItemHover) {
          forced.blur();
          forced.classList.remove('forced-open');
        }
        // prevents menus that have been tabbed into from staying open
        // when you hover over another menu
        e.target.classList.remove('forced-closed');
        e.target.classList.add('forced-open');
      });
      const toggleForcedOpen = e => {
        const isForced = e.target.classList.contains('forced-open');
        const target = e.currentTarget;
        if (isForced) {
          target.removeEventListener('blur', e =>
            target.classList.remove('forced-open')
          );
          target.classList.remove('forced-open');
          target.classList.add('forced-closed');
          target.blur();
          target.parentNode.addEventListener('mouseout', () => {
            target.classList.remove('forced-closed');
          });
        } else {
          target.classList.remove('forced-closed');
          target.classList.add('forced-open');
          target.parentNode.removeEventListener('mouseout', () => {
            target.classList.remove('forced-closed');
          });
        }
        e.target.focus();
      };
      menuItemHover.addEventListener('click', toggleForcedOpen);
      menuItemHover.addEventListener('focus', e => {
        e.target.classList.add('forced-closed');
        e.target.classList.remove('forced-open');
      });

      // ensure focus is removed when esc is pressed
      const focusOutOnEsc = e => {
        if (e.key === 'Escape') {
          const textarea = document.getElementById('code');
          if (e.target == textarea) {
            e.preventDefault();
            textarea.blur();
          }
          else {
            const forcedOpenItem = document.querySelector('.forced-open');
            const target = e.currentTarget;
            if (forcedOpenItem) {
              forcedOpenItem.classList.remove('forced-open');
              forcedOpenItem.blur();
              forcedOpenItem.classList.add('forced-closed');
              e.target.focus();
            }
          }
        }
      };
      document.addEventListener('keydown', focusOutOnEsc);
    });

    // Mobile menu subnav menus
    const header = document.querySelector('.js-header');
    const headerbuttons = document.querySelectorAll('.js-headerMenuButton');
    headerbuttons.forEach(button => {
      button.addEventListener('click', e => {
        e.preventDefault();
        const isActive = header.classList.contains('is-active');
        if (isActive) {
          handleNavigationDrawerInactive(header);
        } else {
          handleNavigationDrawerActive(header);
        }
        button.setAttribute('aria-expanded', isActive);
      });
    });

    const scrim = document.querySelector('.js-scrim');
    scrim.addEventListener('click', e => {
      e.preventDefault();

      // find any active submenus and close them
      const activeSubnavs = document.querySelectorAll(
        '.NavigationDrawer-submenuItem.is-active'
      );
      activeSubnavs.forEach(subnav => handleNavigationDrawerInactive(subnav));

      handleNavigationDrawerInactive(header);

      headerbuttons.forEach(button => {
        button.setAttribute(
          'aria-expanded',
          header.classList.contains('is-active')
        );
      });
    });

    const getNavigationDrawerMenuItems = navigationDrawer => [
      navigationDrawer.querySelector('.NavigationDrawer-header > a'),
      ...navigationDrawer.querySelectorAll(
        ':scope > .NavigationDrawer-nav > .NavigationDrawer-list > .NavigationDrawer-listItem > a, :scope > .NavigationDrawer-nav > .NavigationDrawer-list > .NavigationDrawer-listItem > .Header-socialIcons > a'
      ),
    ];

    const getNavigationDrawerIsSubnav = navigationDrawer =>
      navigationDrawer.classList.contains('NavigationDrawer-submenuItem');

    const handleNavigationDrawerInactive = navigationDrawer => {
      const menuItems = getNavigationDrawerMenuItems(navigationDrawer);
      navigationDrawer.classList.remove('is-active');
      const parentMenuItem = navigationDrawer
        .closest('.NavigationDrawer-listItem')
        ?.querySelector(':scope > a');
      parentMenuItem?.focus();

      menuItems.forEach(item => item.setAttribute('tabindex', '-1'));

      menuItems[0].removeEventListener(
        'keydown',
        handleMenuItemTabLeft.bind(navigationDrawer)
      );
      menuItems[menuItems.length - 1].removeEventListener(
        'keydown',
        handleMenuItemTabRight.bind(navigationDrawer)
      );

      if (navigationDrawer === header) {
        headerbuttons[0]?.focus();
      }
    };

    const handleNavigationDrawerActive = navigationDrawer => {
      const menuItems = getNavigationDrawerMenuItems(navigationDrawer);

      navigationDrawer.classList.add('is-active');
      menuItems.forEach(item => item.setAttribute('tabindex', '0'));
      menuItems[0].focus();

      menuItems[0].addEventListener(
        'keydown',
        handleMenuItemTabLeft.bind(this, navigationDrawer)
      );
      menuItems[menuItems.length - 1].addEventListener(
        'keydown',
        handleMenuItemTabRight.bind(this, navigationDrawer)
      );
    };

    const handleMenuItemTabLeft = (navigationDrawer, e) => {
      if (e.key === 'Tab' && e.shiftKey) {
        e.preventDefault();
        handleNavigationDrawerInactive(navigationDrawer);
      }
    };

    const handleMenuItemTabRight = (navigationDrawer, e) => {
      if (e.key === 'Tab' && !e.shiftKey) {
        e.preventDefault();
        handleNavigationDrawerInactive(navigationDrawer);
      }
    };

    const prepMobileNavigationDrawer = navigationDrawer => {
      const isSubnav = getNavigationDrawerIsSubnav(navigationDrawer);
      const menuItems = getNavigationDrawerMenuItems(navigationDrawer);

      navigationDrawer.addEventListener('keyup', e => {
        if (e.key === 'Escape') {
          handleNavigationDrawerInactive(navigationDrawer);
        }
      });

      menuItems.forEach(item => {
        const parentLi = item.closest('li');
        if (
          parentLi &&
          parentLi.classList.contains('js-mobile-subnav-trigger')
        ) {
          const submenu = parentLi.querySelector(
            '.NavigationDrawer-submenuItem'
          );
          item.addEventListener('click', () => {
            handleNavigationDrawerActive(submenu);
          });
        }
      });
      if (isSubnav) {
        handleNavigationDrawerInactive(navigationDrawer);
        navigationDrawer
          .querySelector('.NavigationDrawer-header')
          .addEventListener('click', e => {
            e.preventDefault();
            handleNavigationDrawerInactive(navigationDrawer);
          });
      }
    };

    document
      .querySelectorAll('.NavigationDrawer')
      .forEach(drawer => prepMobileNavigationDrawer(drawer));
    handleNavigationDrawerInactive(header);
  }

  /**
   * Attempts to detect user's operating system and sets the download
   * links accordingly
   */
  async function setDownloadLinks() {
    const versionElement = document.querySelector('.js-latestGoVersion');
    if (versionElement) {
      const anchorTagWindows = document.querySelector('.js-downloadWin');
      const anchorTagMac = document.querySelector('.js-downloadMac');
      const anchorTagLinux = document.querySelector('.js-downloadLinux');
      const version = await getLatestVersion();

      const macDownloadUrl = `/dl/${version}.darwin-amd64.pkg`;
      const windowsDownloadUrl = `/dl/${version}.windows-amd64.msi`;
      const linuxDownloadUrl = `/dl/${version}.linux-amd64.tar.gz`;

      anchorTagWindows.href = windowsDownloadUrl;
      anchorTagMac.href = macDownloadUrl;
      anchorTagLinux.href = linuxDownloadUrl;

      /*
       * Note: we do not change .js-downloadBtn anymore
       * because it is impossible to tell reliably which architecture
       * the user's browser is running on.
       */
    }
  }

  function registerPortToggles() {
    for (const el of document.querySelectorAll('.js-togglePorts')) {
      el.addEventListener('click', () => {
        el.setAttribute('aria-expanded', el.getAttribute('aria-expanded') === 'true' ? 'false' : 'true')
      })
    }
  }

  /**
   * Retrieves list of Go versions & returns the latest
   */
  async function getLatestVersion() {
    let version = 'go1.17'; // fallback version if fetch fails
    try {
      const versionData = await (await fetch('/dl/?mode=json')).json();
      if (!versionData.length) {
        return version;
      }
      versionData.sort((v1, v2) => {
        return v2.version - v1.version;
      });
      version = versionData[0].version;
    } catch (err) {
      console.error(err);
    }
    return version;
  }

  /**
   * initialThemeSetup sets data-theme attribute based on preferred color
   */

  function initialThemeSetup() {
    const themeCookie = document.cookie.match(
      /prefers-color-scheme=(light|dark|auto)/
    );
    const theme = themeCookie && themeCookie.length > 0 && themeCookie[1];
    if (theme) {
      document.querySelector('html').setAttribute('data-theme', theme);
    }
  }

  /**
   * setThemeButtons sets click listeners for toggling theme buttons
   */
  function setThemeButtons() {
    for (const el of document.querySelectorAll('.js-toggleTheme')) {
      el.addEventListener('click', () => {
        toggleTheme();
      });
    }
  }

  /**
   * setVersionSpan sets the latest version in any span that has this selector.
   */
  async function setVersionSpans() {
    const spans = document.querySelectorAll('.GoVersionSpan');
    if (!spans) return;
    const version = await getLatestVersion();
    Array.from(spans).forEach(span => {
      span.textContent = `Download (${version.replace('go', '')})`
    });
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
    document.cookie = `prefers-color-scheme=${nextTheme};${domain}path=/;max-age=31536000;`;
  }

  function registerCookieNotice() {
    const themeCookie = document.cookie.match(/cookie-consent=true/);
    if (!themeCookie) {
      const notice = document.querySelector('.js-cookieNotice');
      const button = notice.querySelector('button');
      notice.classList.add('Cookie-notice--visible');
      button.addEventListener('click', () => {
        let domain = '';
        if (location.hostname === 'go.dev') {
          // Apply the cookie to *.go.dev.
          domain = 'domain=.go.dev;';
        }
        document.cookie = `cookie-consent=true;${domain}path=/;max-age=31536000`;
        notice.remove();
      });
    }
  }

  initialThemeSetup();

  const onPageLoad = () => {
    registerHeaderListeners();
    setDownloadLinks();
    setThemeButtons();
    setVersionSpans();
    registerPortToggles();
    registerCookieNotice();
  };

  // DOM might be already loaded when we try to setup the callback, hence the check.
  if (document.readyState !== 'loading') {
    onPageLoad();
  } else {
    document.addEventListener('DOMContentLoaded', onPageLoad);
  }
})();
