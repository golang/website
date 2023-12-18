(() => {
  'use strict';
  const BREAKPOINT = 512;
  const logo = document.querySelector('.js-headerLogo');
  const form = document.querySelector('.js-searchForm');
  if (!form) return;
  const button = document.querySelector('.js-searchFormSubmit');
  const input = form.querySelector('input');


  renderForm();

  window.addEventListener('resize', renderForm);

  function renderForm() {
    if (window.innerWidth > BREAKPOINT) {
      logo.classList.remove('Header-logo--hidden');
      form.classList.remove('SearchForm--open');
      input.removeEventListener('focus', showSearchBox);
      input.removeEventListener('keypress', handleKeypress);
      input.removeEventListener('focusout', hideSearchBox);
    } else {
      button.addEventListener('click', handleSearchClick);
      input.addEventListener('focus', showSearchBox);
      input.addEventListener('keypress', handleKeypress);
      input.addEventListener('focusout', hideSearchBox);
    }
  }

  /**
   * Submits form if Enter key is pressed
   * @param {KeyboardEvent} e
   */
  function handleKeypress(e) {
    if (e.key === 'Enter') form.submit();
  }

  /**
   * Shows the search box when it receives focus (expands it from
   * just the spyglass if we're on mobile).
   */
  function showSearchBox() {
    logo.classList.add('Header-logo--hidden');
    form.classList.add('SearchForm--open');
  }

  /**
   * Hides the search box (shrinks to just the spyglass icon).
   */
  function hideSearchBox() {
    logo.classList.remove('Header-logo--hidden');
    form.classList.remove('SearchForm--open');
  }

  /**
   * Expands the searchbox so input is visible and gives
   * the input focus.
   * @param {MouseEvent} e
   */
  function handleSearchClick(e) {
    e.preventDefault();

    showSearchBox();
    input.focus();
  }

  /**
   * Focus search input if '/' key is pressed.
   */
  document.addEventListener('keypress', e => {
    const t = e.target
    if (
      !handler.target &&
      (t?.tagName === 'INPUT' || t?.tagName === 'SELECT' || t?.tagName === 'TEXTAREA')
    ) {
      return;
    }
    if (t?.isContentEditable) {
      return;
    }
    if (
      (handler.withMeta && !(e.ctrlKey || e.metaKey)) ||
      (!handler.withMeta && (e.ctrlKey || e.metaKey))
    ) {
      return;
    }
    if (e.key === '/') {
      const searchInput = Array.from(
        document.querySelectorAll('.js-searchForm input')
      ).pop();
      // Favoring the Firefox quick find feature over search input
      // focus. See: go.dev/issue/41093.
      if (searchInput && !window.navigator.userAgent.includes('Firefox')) {
        e.preventDefault();
        searchInput.focus();
      }
    }
  })
})();
