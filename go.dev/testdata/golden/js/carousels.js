// Copyright 2020 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.

/**
 * Handles Carousel logic. Hides right/left buttons if user is at the end of the
 * list of slides. Transitions the movement with CSS, and blocks slide changes
 * during the transitions with `allowShift`.
 */
(() => {
  'use strict';

  /** TODO: refactor to slide container using transformX instead of left,
   * in order to leverage GPU rendering
   */

  function initCarousel(container, prev, next, breakpoints) {
    const slides = Array.from(container.children);
    let posInitial;
    const slidesCount = container.children.length;
    let distanceToMove = Math.floor(
      container.children[0].getBoundingClientRect().width
    );
    let index = 0;
    let allowShift = true;
    let groupCountAdjustment = getGroupCountAdjustment();

    adjustTabbableSlides(breakpoints);

    // Click events.
    prev.addEventListener('click', () => {
      shiftSlide(-1);
    });
    next.addEventListener('click', () => {
      shiftSlide(1);
    });

    // Transition events.
    container.addEventListener('transitionend', handleTransitionEnd);

    window.addEventListener(
      'resize',
      debounce(() => {
        // the next 3 lines reinitializes the slide container node
        container.style.display = 'none';
        container.offsetHeight;
        container.style.display = 'flex';

        groupCountAdjustment = getGroupCountAdjustment();
        distanceToMove = Math.floor(
          container.children[0].getBoundingClientRect().width
        );
        container.style.left = -(distanceToMove * index) + 'px';
      })
    );

    function shiftSlide(dir, action) {
      container.classList.add('shifting');
      if (allowShift) {
        if (!action) {
          posInitial = container.offsetLeft;
        }
        if (dir === 1) {
          prev.removeAttribute('hidden');
          container.style.left = posInitial - distanceToMove + 'px';
          index++;
        } else if (dir === -1) {
          next.removeAttribute('hidden');
          container.style.left = posInitial + distanceToMove + 'px';
          index--;
        } else {
          container.style.left = posInitial + 'px';
        }
        if (index === 0) {
          prev.setAttribute('hidden', true);
        } else if (index === slidesCount - 1 - groupCountAdjustment) {
          next.setAttribute('hidden', true);
        }
      }
      allowShift = false;
      adjustTabbableSlides(breakpoints);
    }

    /**
     * Only allow visible slides to be tabbable
     */
    function adjustTabbableSlides(breakpoints) {
      let count = 1;
      if (breakpoints) {
        for (const bp of breakpoints) {
          if (window.innerWidth > bp.px) {
            count = bp.groupCount;
          }
        }
      }
      slides.forEach((slide, i) => {
        const links = slide.querySelectorAll('a');
        if (i >= index && i < index + count) {
          links.forEach(link => (link.tabIndex = 0));
        } else {
          links.forEach(link => (link.tabIndex = -1));
        }
      });
    }

    function handleTransitionEnd() {
      container.classList.remove('shifting');
      allowShift = true;
    }

    function getGroupCountAdjustment() {
      let groupCountAdjustment = 0;

      if (breakpoints) {
        for (const bp of breakpoints) {
          if (window.innerWidth > bp.px) {
            groupCountAdjustment = bp.groupCount - 1;
          }
        }
      }
      return groupCountAdjustment;
    }
  }

  // Build quotes carousel.
  const quotesSliderContainer = document.querySelector(
    '.js-testimonialsGoQuotes'
  );
  const quotesPrev = document.querySelector('.js-testimonialsPrev');
  const quotesNext = document.querySelector('.js-testimonialsNext');
  // Build events carousel.
  const eventsSliderContainer = document.querySelector(
    '.js-goCarouselEventsSlides'
  );
  const eventsPrev = document.querySelector('.js-eventsCarouselPrev');
  const eventsNext = document.querySelector('.js-eventsCarouselNext');

  // Build Solutions hero carousel.
  const solutionsCarouselSliderContainer = document.querySelector(
    '.js-solutionsHeroCarouselSlides'
  );
  const solutionsCarouselPrev = document.querySelector(
    '.js-solutionsHeroCarouselPrev'
  );
  const solutionsCarouselNext = document.querySelector(
    '.js-solutionsHeroCarouselNext'
  );

  if (quotesSliderContainer) {
    initCarousel(quotesSliderContainer, quotesPrev, quotesNext);
  }
  if (eventsSliderContainer) {
    const breakpoints = [
      {px: 768, groupCount: 2},
      {px: 1068, groupCount: 3},
    ];
    initCarousel(eventsSliderContainer, eventsPrev, eventsNext, breakpoints);
  }
  if (solutionsCarouselSliderContainer) {
    initCarousel(
      solutionsCarouselSliderContainer,
      solutionsCarouselPrev,
      solutionsCarouselNext
    );
  }
})();

/**
 * Debounce functions for better performance
 * (c) 2018 Chris Ferdinandi, MIT License, https://gomakethings.com
 * @param  {Function} fn The function to debounce
 */
function debounce(fn) {
  let timeout;
  return (...args) => {
    if (timeout) {
      window.cancelAnimationFrame(timeout);
    }
    timeout = window.requestAnimationFrame(function () {
      fn(args);
    });
  };
}
