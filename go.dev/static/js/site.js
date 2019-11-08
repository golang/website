// Copyright 2019 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/**
 * A bit of navigation related code for handling dismissible elements.
 */
(function() {
  'use strict';

  function registerHeaderListeners() {
    const header = document.querySelector('.js-header');
    const menuButtons = document.querySelectorAll('.js-headerMenuButton');
    menuButtons.forEach(button => {
      button.addEventListener('click', e => {
        e.preventDefault();
        header.classList.toggle('is-active');
        button.setAttribute(
          'aria-expanded',
          header.classList.contains('is-active')
        );
      });
    });

    const scrim = document.querySelector('.js-scrim');
    scrim.addEventListener('click', e => {
      e.preventDefault();
      header.classList.remove('is-active');
      menuButtons.forEach(button => {
        button.setAttribute(
          'aria-expanded',
          header.classList.contains('is-active')
        );
      });
    });
  }

  window.addEventListener('DOMContentLoaded', () => {
    registerHeaderListeners();
  });

  // Register feedback listeners.
  window.addEventListener('load', () => {
    const buttons = document.querySelectorAll('.js-feedbackButton');
    buttons.forEach(button => {
      button.addEventListener('click', sendFeedback);
    });
  });

  // Launches the feedback interface.
  function sendFeedback() {
    userfeedback.api.startFeedback({ productId: '5131929', bucket: 'Default' });
  }

  window.dataLayer = window.dataLayer || [];
  function gtag() {
    dataLayer.push(arguments);
  }
  gtag('js', new Date());
  gtag('config', 'UA-141356704-1');
})();
