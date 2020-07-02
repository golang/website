// 'More projects' button on use case pages
(() => {
  const button = document.querySelector('.js-more-projects-btn');
  if (!button) return;
  const hiddenProjects = document.querySelectorAll('.FeaturedUsers-row.hidden');
  button.addEventListener('click', () => {
    button.classList.add('hidden');
    hiddenProjects.forEach(project => {
      project.classList.remove('hidden');
    });
  });
})();

/* Header tags generated with markdown - inner span needed for correct scroll
  position */
(() => {
  const headingHashes = Array.from(document.querySelectorAll('h2[id]'));
  headingHashes.forEach(h2 => {
    const text = h2.textContent;
    const id = h2.id;
    h2.id = id + '-h2';
    h2.textContent = '';
    const span = document.createElement('span');
    span.textContent = text;
    span.id = id;
    h2.appendChild(span);
  });
})();

// Use case pages section navigation
(() => {
  const stickyNav = document.querySelector('.js-useCaseStickyNav');
  if (stickyNav) {
    const linkData = {
      'overview': 'Overview',
      'key-benefits': 'Key Benefits',
      'use-case': 'Use Case',
      'featured-users': 'Featured Users',
      'get-started': 'Get Started',
    };

    const container = document.querySelector('.js-useCaseSubnav');
    const subNavAnchorLinks = document.querySelector('.js-useCaseSubnavLinks');
    const header = document.querySelector('.js-useCaseSubnavHeader');
    const icon = document.querySelector('.js-useCaseSubnavMenuIcon');
    const menu = document.querySelector('.js-useCaseSubnavMenu');
    const contentBody = document.querySelector('.js-useCaseContentBody');
    const headerHeightPx = 56;
    const sectionHeadings = Array.from(
      document.querySelectorAll('.sectionHeading')
    ).map(h => h.firstChild);
    let distanceFromTop =
      window.pageYOffset +
      contentBody.getBoundingClientRect().top -
      headerHeightPx;

    if (!header || !menu) return;
    container.addEventListener('click', handleClick);
    container.addEventListener('keydown', handleKeydown);

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
    });

    function setStickyNav() {
      if (window.scrollY > distanceFromTop) {
        stickyNav.classList.add('UseCaseSubNav-anchorLinks--sticky');
      } else {
        stickyNav.classList.remove('UseCaseSubNav-anchorLinks--sticky');
      }
    }
  }
})();
