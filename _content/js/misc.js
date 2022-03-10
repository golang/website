// 'More projects' button on use case pages
(() => {
  const button = document.querySelector('.js-moreProjectsBtn');
  if (!button) return;
  const hiddenProjects = document.querySelectorAll('.js-featuredUsersRow[hidden]');
  button.addEventListener('click', () => {
    button.setAttribute('hidden', true);
    hiddenProjects.forEach(project => {
      project.removeAttribute('hidden');
    });
  });
})();

// Use case pages section navigation
(() => {
  const stickyNav = document.querySelector('.js-useCaseStickyNav');
  if (!stickyNav) return;
  const linkData = {
    'overview': 'Overview',
    'key-benefits': 'Key Benefits',
    'use-case': 'Use Case',
    'featured-users': 'Featured Users',
    'get-started': 'Get Started',
  };
  const container = document.querySelector('.js-useCaseSubnav');
  const subNavAnchorLinks = document.querySelector('.js-useCaseSubnavLinks');
  const siteHeader = document.querySelector('.js-siteHeader');
  const header = document.querySelector('.js-useCaseSubnavHeader');
  const icon = document.querySelector('.js-useCaseSubnavMenuIcon');
  const menu = document.querySelector('.js-useCaseSubnavMenu');
  const contentBody = document.querySelector('.js-useCaseContentBody');
  const headerHeightPx = 56;
  const sectionHeadings = Array.from(
    document.querySelectorAll('.sectionHeading')
  );
  let distanceFromTop =
    window.pageYOffset +
    contentBody.getBoundingClientRect().top -
    headerHeightPx;
  if (!header || !menu) return;
  container.addEventListener('click', handleClick);
  container.addEventListener('keydown', handleKeydown);
  changeScrollPosition();

  function handleClick(event) {
    if (event.target === header) {
      toggleMenu();
    } else {
      closeMenu();
    }
  }

  function handleKeydown(event) {
    if (event.key === 'Enter') {
      closeMenu();
    } else {
      openMenu();
    }
  }

  function openMenu() {
    menu.classList.add('UseCaseSubNav-menu--open');
    icon.classList.add('UseCaseSubNav-menuIcon--open');
  }

  function closeMenu() {
    menu.classList.remove('UseCaseSubNav-menu--open');
    icon.classList.remove('UseCaseSubNav-menuIcon--open');
  }

  function toggleMenu() {
    menu.classList.toggle('UseCaseSubNav-menu--open');
    icon.classList.toggle('UseCaseSubNav-menuIcon--open');
  }

  sectionHeadings.forEach(heading => {
    let a = document.createElement('a');
    a.classList.add('UseCase-anchorLink', 'anchor-link');
    a.href = `${window.location.pathname}#${heading.id}`;
    a.textContent = linkData[heading.id];
    stickyNav.appendChild(a);
    a = a.cloneNode();
    a.textContent = linkData[heading.id];
    subNavAnchorLinks.appendChild(a);
  });

  // Selected section styles
  const anchorLinks = document.querySelectorAll('.anchor-link');
  anchorLinks.forEach(link => {
    link.addEventListener('click', () => {
      document
        .querySelectorAll('.anchor-link')
        .forEach(el => el.classList.remove('selected'));
      link.classList.add('selected');
    });
  });

  window.addEventListener('scroll', () => {
    // delay in case the user clicked the anchor link and we are autoscrolling
    setTimeout(setSelectedAnchor, 500);
  });

  function setSelectedAnchor() {
    for (heading of sectionHeadings) {
      const {offsetTop} = heading;
      if (offsetTop > window.scrollY) {
        anchorLinks.forEach(link => {
          const anchorId = link.href.slice(link.href.indexOf('#') + 1);
          if (anchorId === heading.id) {
            link.classList.add('selected');
          } else {
            link.classList.remove('selected');
          }
        });
        break;
      }
    }
  }

  /* sticky nav logic -- uses content for y position because reloading page
  when not scrolled to the top creates bug if using current y position of
  sticky nav */
  window.addEventListener('scroll', setStickyNav);

  window.addEventListener('resize', () => {
    distanceFromTop =
      window.pageYOffset +
      contentBody.getBoundingClientRect().top -
      headerHeightPx;

    changeScrollPosition()
  });

  /**
   * Changes scroll position according to the size of the header and menu
   * Also changes according to the user's browser
   */
  function changeScrollPosition() {
    const SUPPORTS_SCROLL_BEHAVIOR = document.body.style.scrollBehavior !== undefined;
    const WINDOW_WIDTH_BREAKPOINT_PX = 923;
    let scrollPosition = headerHeightPx;

    if (SUPPORTS_SCROLL_BEHAVIOR) {
      if (window.innerWidth < WINDOW_WIDTH_BREAKPOINT_PX) {
        scrollPosition += header.clientHeight;
      }
    } else {
      if (window.innerWidth >= WINDOW_WIDTH_BREAKPOINT_PX) {
        scrollPosition = siteHeader.clientHeight
      } else {
        scrollPosition = siteHeader.clientHeight + header.clientHeight;
      }
    }
  }

  function setStickyNav() {
    if (window.scrollY > distanceFromTop) {
      stickyNav.classList.add('UseCaseSubNav-anchorLinks--sticky');
    } else {
      stickyNav.classList.remove('UseCaseSubNav-anchorLinks--sticky');
    }
  }
})();
