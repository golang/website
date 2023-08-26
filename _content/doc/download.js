class DownloadsController {
  constructor() {
    // Parts of tabbed section.
    this.tablist = document.querySelector('.js-tabSection');
    this.tabs = this.tablist.querySelectorAll('[role="tab"]');
    this.panels = document.querySelectorAll('[role="tabpanel"]');

    // OS for which to display download and install steps.
    this.osName = 'Unknown OS';
    this.osNameFromQuery = '';

    // URL to JSON containing list of installer downloads.
    const fileListUrl = '/dl/?mode=json';
    this.activeTabIndex = 0;

    const dlQuery = new URL(document.URL).searchParams.get('dc') || '';
    if (dlQuery !== '') {
      const [queryOS] = dlQuery.split('-');
      if (queryOS === 'darwin') this.osNameFromQuery = 'mac';
      if (queryOS === 'windows') this.osNameFromQuery = 'windows';
      if (queryOS === 'linux') this.osNameFromQuery = 'linux';
    }

    // Get the install file list, then get names and sizes
    // for each OS supported on the install page.
    fetch(fileListUrl)
      .then((response) => response.json())
      .then((data) => {
        const files = data[0]['files'];
        for (var i = 0; i < files.length; i++) {
          let file = files[i].filename;
          if (file.match('.linux-amd64.tar.gz$')) {
            this.linuxFileName = file;
          }
        }
        this.detectOS();
        const osTab = document.getElementById(this.osName);
        if (osTab !== null) {
          osTab.click();
        }
        this.setVersion(data[0].version);
      })
      .catch(console.error);
      this.setEventListeners();
  }

  setEventListeners() {
    this.tabs.forEach((tabEl) => {
      tabEl.addEventListener('click', e => this.handleTabClick((e)));
    });
  }

  // Set the download button UI version.
  setVersion(latest) {
    document.querySelector('.js-downloadDescription').textContent =
      `Download (${this.parseVersionNumber(latest)})`;
  }

  // Updates install tab with dynamic data.
  setInstallTabData(osName) {
    const fnId = `#${osName}-filename`;
    const el = document.querySelector(fnId);
    if (!el) {
      return;
    }
    switch(osName) {
      case 'linux':
        // Update filename for linux installation step
        if (this.linuxFileName) {
          el.textContent = this.linuxFileName;
        }
        break;
    }
  }

  // Detect the users OS for installation default.
  detectOS() {
    if (this.osNameFromQuery !== '') {
      this.osName = this.osNameFromQuery;
      return;
    }
    if (navigator.userAgent.indexOf('Linux') !== -1) {
      this.osName = 'linux';
    } else if (navigator.userAgent.indexOf('Mac') !== -1) {
      this.osName = 'mac';
    } else if (navigator.userAgent.indexOf('X11') !== -1) {
      this.osName = 'unix';
    } else if (navigator.userAgent.indexOf('Win') !== -1) {
      this.osName = 'windows';
    }
  }

  // Activates the tab at the given index.
  activateTab(index) {
    this.tabs[this.activeTabIndex].setAttribute('aria-selected', 'false');
    this.tabs[this.activeTabIndex].setAttribute('tabindex', '-1');
    this.panels[this.activeTabIndex].setAttribute('hidden', '');
    this.tabs[index].setAttribute('aria-selected', 'true');
    this.tabs[index].setAttribute('tabindex', '0');
    this.panels[index].removeAttribute('hidden');
    this.tabs[index].focus();
    this.activeTabIndex = index;
  }

  // Handles clicks on tabs.
  handleTabClick(e) {
    const el = (e.target);
    this.activateTab(Array.prototype.indexOf.call(this.tabs, el));
    this.setInstallTabData(el.id);
  }

  // get version number.
  parseVersionNumber(string) {
    const rx = /(\d+\.)(\d+)(\.\d+)?/g;
    const matches = rx.exec(string);
    if (matches?.[0]) {
      return matches[0];
    } else {
      return '';
    }
  }

}

// Instantiate controller for page event handling.
new DownloadsController();
