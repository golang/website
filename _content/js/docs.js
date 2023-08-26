// when page loads
window.addEventListener('load', () => {
  // Load the toc outline for a documentation page
  const tocRoot = document.querySelector("[role='tree'].js-toc-tree");
  if (tocRoot) {
    const hasNestedList = addLinksToTOC(document, tocRoot);
    const toggleTreeRef = document.querySelector('.js-expand-toc');
    let expanded = false;
    toggleTreeRef.children[1].style.display = 'none';
    if (!hasNestedList) toggleTreeRef.children[0].style.display = 'none';
    if (hasNestedList) {
      toggleTreeRef.addEventListener('click', () => {
        expanded = !expanded;
        toggleTreeRef.children[expanded ? 0 : 1].style.display = 'none';
        toggleTreeRef.children[!expanded ? 0 : 1].style.display = 'inline-block';
      });
    }

    new Tree(tocRoot, {
      useScrollObservers: true,
      toggleTreeRef: hasNestedList ? toggleTreeRef : null,
    });
  }

  // Load the doc outline for the site
  const outlineRoot = document.querySelector("[role='tree'].js-outline-tree");
  if (outlineRoot) {
    new Tree(outlineRoot, {
      useScrollObservers: false,
      allowMultipleTreeExpansion: true,
    });
  }
});

/**
 * Tree is the navigation index component of the documentation page.
 * It adds accessiblity attributes to a tree, observes the heading elements
 * focus the topmost link for headings visible on the page, and implements the
 * WAI-ARIA Treeview Design Pattern with full
 * [keyboard support](https://www.w3.org/TR/wai-aria-practices/examples/treeview/treeview-2/treeview-2a.html#kbd_label).
 */

/**
 * @typedef {object} TreeOpt
 * @property {boolean} useScrollObservers
 * @property {boolean} allowMultipleTreeExpansion
 * @property {HTMLElement} toggleTreeRef
 */

class Tree {
  /**
   * 
   * @param {HTMLElement} node 
   * @param {TreeOpt} opts
   */
  constructor(node, opts) {
    const {
      useScrollObservers = false,
      allowMultipleTreeExpansion = false,
      toggleTreeRef,
    } = opts;

    /** @type {HTMLElement} */
    this.node = node;

    /** @type {HTMLElement} */
    this.toggleTreeRef = toggleTreeRef;

    /** @type {boolean} */
    this.useScrollObservers = useScrollObservers;

    /** @type {boolean} */
    this.allowMultipleTreeExpansion = allowMultipleTreeExpansion;

    /** @type {Treeitem[]} */
    this.treeitems = [];

    /** @type {string[]} */
    this.firstChars = [];

    /** @type {Treeitem} */
    // this.firstTreeitem = null;
    /** @type {Treeitem} */
    this.lastTreeitem = null;
    /** @type {Treeitem} */
    this.selectedItem = null;
    this.observerCallbacks = [];

    this.init();
  }

  init() {
    /** @type {Treeitem} */
    let activeTreeitem = null;
    /**
     * 
     * @param {Element} node 
     * @param {Tree} tree 
     * @param {Treeitem} group 
     */
    const findTreeitems = (node, group) => {
      let elem = node.firstElementChild;
      let ti = group;
      const urlPath = window.location.pathname;
      const urlHash = window.location.hash;
      while (elem) {
        if (elem.tagName === 'A' || elem.tagName === 'SPAN') {
          ti = new Treeitem(elem, this, group);
          if (!this.useScrollObservers) {
            if (urlPath.includes(ti.node.getAttribute('href'))) {
              activeTreeitem = ti;
            }
          } else {
            if (urlHash.includes(ti.node.getAttribute('href'))) {
              activeTreeitem = ti;
            }
          }
          this.treeitems.push(ti);
          this.firstChars.push(ti.label.substring(0, 1).toLowerCase());
        }
        if (elem.firstElementChild) {
          findTreeitems(elem, ti);
        }
        elem = elem.nextElementSibling;
      }
    }

    findTreeitems(this.node, null);
    // only proceed if treeitems exist
    if (!this.treeitems.length) {
      // if there was a treeitems toggler, let's remove it
      if (this.toggleTreeRef) {
        this.toggleTreeRef.parentElement.style.display = 'none';
      }
      return;
    }

    this.treeitems.forEach((ti, idx) => (ti.index = idx));
    this.updateVisibleTreeitems();
    if (this.useScrollObservers) {
      this.observeTargets();
      this.firstTreeitem.node.tabIndex = 0;
    }

    if (activeTreeitem) {
      this.setSelectedToItem(activeTreeitem);
      let groupedItem = activeTreeitem.groupTreeitem;
      while (groupedItem) {
        this.expandTreeitem(groupedItem);
        groupedItem = groupedItem.groupTreeitem;
      }
    }

    /** Add event handlers to manage all treeitems */
    if (this.toggleTreeRef) {
      this.toggleTreeRef.addEventListener('click', () => {
        this.allowMultipleTreeExpansion = !this.allowMultipleTreeExpansion;
        this.treeitems.forEach((ti) => {
          if (ti.isExpandable) {
            if (ti.isExpanded()) {
              ti.tree.collapseTreeitem(ti);
            } else {
              ti.tree.expandTreeitem(ti);
            }
          }
        });
      });
    }
  }

  addObserver(fn, delay = 100) {
    this.observerCallbacks.push(debounce(fn, delay))
  }

  observeTargets() {
    this.addObserver(item => {
      this.expandTreeitem(item);
      this.setSelectedToItem(item);
    });

    const targets = {};
    const observer = new IntersectionObserver((entries) => {
      for (const entry of entries) {
        const intersecting = entry.isIntersecting || entry.intersectionRatio === 1;
        targets[entry.target.id] = intersecting;
      }
      for (const [id, isIntersecting] of Object.entries(targets)) {
        if (isIntersecting) {
          const active = this.treeitems.find(t => t.node?.href.endsWith(`#${id}`));
          if (active) {
            for (const fn of this.observerCallbacks) {
              fn(active);
            }
          }
          break;
        }
      }
    }, { threshold: 1, rootMargin: '0px 0px 0px 0px' });

    const hrefs = this.treeitems.map(t => t.node.getAttribute('href'));
    for (const href of hrefs) {
      if (href) {
        const id = href.replace(window.location.origin, '').replace('/', '').replace('#', '');
        const target = document.getElementById(id);
        if (target) observer.observe(target);
      }
    }
  }

  updateVisibleTreeitems() {
    this.firstTreeitem = this.treeitems[0];
    for (let i = 0; i < this.treeitems.length; i++) {
      const ti = this.treeitems[i];

      /** @type {HTMLElement} */
      let parent = ti.node.parentNode;
      ti.isVisible = true;

      while (parent && parent !== this.node) {
        if (parent.getAttribute('aria-expanded') === 'false') {
          ti.isVisible = false;
        }
        parent = parent.parentNode;
      }
      if (ti.isVisible) {
        this.lastTreeitem = ti;
      }
    }
  }

  /**
   * 
   * @param {Treeitem} currentItem 
   */
  setSelectedToItem(currentItem) {
    if (!this.allowMultipleTreeExpansion) {
      for (const l1 of this.node.querySelectorAll('[aria-expanded="true"]')) {
        if (l1 === currentItem.node) continue;
        if (!l1.nextElementSibling?.contains(currentItem.node)) {
          l1.setAttribute('aria-expanded', 'false');
        }
      }
    }
    for (const l1 of this.node.querySelectorAll('[aria-selected]')) {
      if (l1 !== currentItem.node) {
        l1.setAttribute('aria-selected', 'false');
      }
    }
    currentItem.node.setAttribute('aria-selected', 'true');
    this.setFocusToItem(currentItem, false);
  }

  /**
   * 
   * @param {Treeitem} treeitem 
   */
  setFocusToItem(treeitem, focusEl = true) {
    treeitem.node.tabIndex = 0;
    if (focusEl) treeitem.node.focus();

    for (const ti of this.treeitems) {
      if (ti !== treeitem) {
        ti.node.tabIndex = -1;
      }
    }
  }

  /**
   * 
   * @param {Treeitem} currentItem 
   */
  setFocusToNextItem(currentItem) {
    let nextItem = null;
    for (let i = this.treeitems.length - 1; i >= 0; i--) {
      const ti = this.treeitems[i];
      if (ti === currentItem) break;
      if (ti.isVisible) nextItem = ti;
    }
    if (nextItem) this.setFocusToItem(nextItem);
  }

  /**
   * 
   * @param {Treeitem} currentItem 
   */
  setFocusToPreviousItem(currentItem) {
    let prevItem = null;
    for (let i = 0; i < this.treeitems.length; i++) {
      const ti = this.treeitems[i];
      if (ti === currentItem) break;
      if (ti.isVisible) prevItem = ti;
    }
    if (prevItem) this.setFocusToItem(prevItem);
  }

  /**
   * 
   * @param {Treeitem} currentItem 
   */
  setFocusToParentItem(currentItem) {
    if (currentItem.groupTreeitem) {
      this.setFocusToItem(currentItem.groupTreeitem);
    }
  }

  setFocusToFirstItem() {
    this.setFocusToItem(this.firstTreeitem);
  }

  setFocusToLastItem() {
    this.setFocusToItem(this.lastTreeitem);
  }

  /**
   * 
   * @param {Treeitem} currentItem 
   */
  expandTreeitem(currentItem) {
    if (currentItem.isExpandable) {
      currentItem.node.setAttribute('aria-expanded', 'true');
      this.updateVisibleTreeitems();
    }
  }

  /**
   * 
   * @param {Treeitem} currentItem 
   */
  expandAllSiblingItems(currentItem) {
    for (let i = 0; i < this.treeitems.length; i++) {
      const ti = this.treeitems[i];
      if (ti.groupTreeitem === currentItem.groupTreeitem && ti.isExpandable) {
        this.expandTreeitem(ti);
      }
    }
  }

  /**
   * 
   * @param {Treeitem} currentItem 
   */
  collapseTreeitem(currentItem) {
    /** @type {Treeitem} */
    let groupTreeitem = null;
    if (currentItem.isExpanded()) {
      groupTreeitem = currentItem;
    } else {
      groupTreeitem = currentItem.groupTreeitem;
    }
    if (groupTreeitem) {
      groupTreeitem.node.setAttribute('aria-expanded', 'false');
      this.updateVisibleTreeitems();
      this.setFocusToItem(groupTreeitem);
    }
  }

  /**
   * 
   * @param {Treeitem} currentItem 
   * @param {string} char 
   */
  setFocusByFirstCharactor(currentItem, char) {
    let start, index;

    char = char.toLowerCase();
    // Get start index for search based on position of currentItem
    start = this.treeitems.indexOf(currentItem) + 1;
    if (start === this.treeitems.length) {
      start = 0;
    }
    // Check remaining slots in the menu
    index = this.getIndexFirstChars(start, char);
    if (index === -1) {
      index = this.getIndexFirstChars(0, char);
    }
    if (index > -1) this.setFocusToItem(this.treeitems[index]);
  }

  /**
   * 
   * @param {number} startIndex 
   * @param {string} char 
   */
  getIndexFirstChars(startIndex, char) {
    for (let i = startIndex; i < this.firstChars.length; i++) {
      if (this.treeitems[i].isVisible) {
        if (char === this.firstChars[i]) {
          return i;
        }
      }
    }
    return -1;
  }
}

class Treeitem {
  /**
   * 
   * @param {HTMLElement} node 
   * @param {Tree} treeObj 
   * @param {Treeitem} group 
   */
  constructor(node, treeObj, group) {
    node.tabIndex = -1;
    this.tree = treeObj;
    this.node = node;
    this.groupTreeitem = group;
    this.label = node.textContent.trim();
    this.depth = (group?.depth || 0) + 1;
    this.index = 0;

    const parent = node.parentElement;
    if (parent?.tagName.toLowerCase() === 'li') {
      parent?.setAttribute('role', 'none');
    }

    node.setAttribute('aria-level', this.depth + '');
    if (node.getAttribute('aria-label')) {
      this.label = node.getAttribute('aria-label').trim();
    }

    this.isExpandable = false;
    this.isVisible = false;
    this.inGroup = false;

    if (group) this.inGroup = true;
    
    let elem = node.nextElementSibling;
    while (elem) {
      if (elem.tagName === 'UL') {
        const groupId = `${group?.label ?? ''} index group ${this.label}`.replace(/[\W_]+/g, '_');
        node.setAttribute('aria-owns', groupId);
        node.setAttribute('aria-expanded', 'false');
        elem.setAttribute('role', 'group');
        elem.setAttribute('id', groupId);
        this.isExpandable = true;
        break;
      }
      elem = elem.nextElementSibling;
    }

    this.key = Object.freeze({
      SPACE: ' ',
      RETURN: 'Enter',
      ARROW_UP: 'ArrowUp',
      ARROW_DOWN: 'ArrowDown',
      ARROW_RIGHT: 'ArrowRight',
      ARROW_LEFT: 'ArrowLeft',
      HOME: 'Home',
      END: 'End',
    });

    this.init();
  }

  init() {
    this.node.tabIndex = -1;
    if (!this.node.getAttribute('role')) {
      this.node.setAttribute('role', 'treeitem');
    }
    this.node.addEventListener('keydown', this.handleKeydown.bind(this));
    this.node.addEventListener('click', this.handleClick.bind(this));
    this.node.addEventListener('focus', this.handleFocus.bind(this));
    this.node.addEventListener('blur', this.handleBlur.bind(this));

    if (!this.isExpandable) {
      this.node.addEventListener('mouseover', this.handleMouseOver.bind(this));
      this.node.addEventListener('mouseout', this.handleMouseOut.bind(this));
    }
  }

  isExpanded() {
    if (this.isExpandable) {
      return this.node.getAttribute('aria-expanded') === 'true';
    }
    return false;
  }

  isSelected() {
    return this.node.getAttribute('aria-selected') === 'true';
  }

  /**
   * 
   * @param {KeyboardEvent} event 
   */
  handleKeydown(event) {
    let flag = false;
    let char = event.key;

    const isPrintableCharacter = (str) => str.length === 1 && str.match(/\S/);
    /**
     * 
     * @param {Treeitem} item 
     */
    const printableCharacter = (item) => {
      if (char === "*") {
        item.tree.expandAllSiblingItems(item);
        flag = true;
      } else {
        if (isPrintableCharacter(char)) {
          item.tree.setFocusByFirstCharactor(item, char);
          flag = true;
        }
      }
    }

    if (event.altKey || event.ctrlKey || event.metaKey) return;

    switch (event.key) {
      case this.key.RETURN:
      case this.key.SPACE:
        if (this.isExpandable) {
          if (this.isExpanded() && this.isSelected()) {
            this.tree.collapseTreeitem(this)
          } else {
            this.tree.expandTreeitem(this)
          }
          flag = true;
        } else {
          event.stopPropagation();
        }
        this.tree.setSelectedToItem(this);
        break;
      
      case this.key.ARROW_UP:
        this.tree.setFocusToPreviousItem(this);
        flag = true;
        break;
      
      case this.key.ARROW_DOWN:
        this.tree.setFocusToNextItem(this);
        flag = true;
        break;
      
      case this.key.ARROW_RIGHT:
        if (this.isExpandable) {
          if (this.isExpanded()) this.tree.setFocusToNextItem(this);
          else this.tree.expandTreeitem(this);
        }
        flag = true;
        break;

      case this.key.ARROW_LEFT:
        if (this.isExpandable && this.isExpanded()) {
          this.tree.collapseTreeitem(this);
          flag = true;
        } else {
          if (this.inGroup) {
            this.tree.setFocusToParentItem(this);
            flag = true;
          }
        }
        break;

      case this.key.HOME:
        this.tree.setFocusToFirstItem();
        flag = true;
        break;

      case this.key.END:
        this.tree.setFocusToLastItem();
        flag = true;
        break;

      default:
        if (isPrintableCharacter(char)) {
          printableCharacter(this);
        }
        break;
    }

    if (flag) {
      event.stopPropagation();
      event.preventDefault();
    }
  }

  /**
   * 
   * @param {MouseEvent} event 
   */
  handleClick(event) {
    if (event.target !== this.node && event.target !== this.node.firstElementChild) {
      return;
    }

    if (this.isExpandable) {
      if (this.isExpanded() && this.isSelected()) {
        this.tree.collapseTreeitem(this)
      } else {
        this.tree.expandTreeitem(this);
      }
      event.stopPropagation();
    }
    this.tree.setSelectedToItem(this);
  }

  handleFocus() {
    let node = this.node;
    if (this.isExpandable) node = node.firstElementChild;
    node?.classList.add('focus');
  }

  handleBlur() {
    let node = this.node;
    if (this.isExpandable) node = node.firstElementChild;
    node?.classList.remove('focus');
  }

  /**
   * 
   * @param {MouseEvent} event 
   */
  handleMouseOver(event) {
    event.currentTarget?.classList.add('hover')
  }

  /**
   * 
   * @param {MouseEvent} event 
   */
  handleMouseOut(event) {
    event.currentTarget?.classList.remove('hover');
  }
}

function debounce(fn, wait) {
  let timeout;
  return (...args) => {
    const later = () => {
      timeout = null;
      fn(...args);
    };
    if (timeout) clearTimeout(timeout);
    timeout = setTimeout(later, wait);
  }
}

/**
 * 
 * @param {HTMLElement} contentEl where we find the header elements h1...h6
 * @param {HTMLElement} tocEl where we inject the table of content tree
 */
function addLinksToTOC(contentEl, tocEl) {
  const headings = contentEl.querySelectorAll("h2, h3, h4, h5, h6");
  let tocItems;
  let prevLevel = false;
  let parentEl = tocEl;
  let prevEl = parentEl;
  let tocHasNestedEl = false;
  for (let i = 0, j = headings.length; i < j; i++) {
    let heading = headings[i];
    const text = heading.textContent;
    const id = text.replace(/[\W_]+/g, '_');
    heading.id = id;
    tocItems = tocEl.getElementsByTagName("li");

    let level = parseInt(heading.tagName.replace(/\D/g, ''));
    if (prevLevel) {
      if (level > prevLevel) {
        tocHasNestedEl = true;
        let ul = document.createElement("ul");
        tocItems[tocItems.length - 1].appendChild(ul);
        parentEl = ul;
        prevEl = tocItems[tocItems.length - 1].parentElement;
      } else if (level < prevLevel) {
        parentEl = prevEl;
      }
    }
    prevLevel = level;

    const li = document.createElement("li");
    const a = document.createElement("a");
    const aText = document.createTextNode(text);
    a.href = "#" + heading.id;
    a.appendChild(aText);
    li.appendChild(a);
    parentEl.appendChild(li);
  }
  return tocHasNestedEl;
}
