class DownloadsController {
  constructor() {
    // Parts of tabbed section.
    this.tablist = document.querySelector('.js-tabSection');
    this.tabs = this.tablist.querySelectorAll('[role="tab"]');
    this.panels = document.querySelectorAll('[role="tabpanel"]');

    // OS for which to display download and install steps.
    this.osName = 'Unknown OS';

    // URL to JSON containing list of installer downloads.
    const fileListUrl = '/dl/?mode=json';
    this.activeTabIndex = 0;

    // Get the install file list, then get names and sizes
    // for each OS supported on the install page.
    fetch(fileListUrl)
      .then((response) => response.json())
      .then((data) => {
        const files = data[0]['files'];
        for (var i = 0; i < files.length; i++) {
          let file = files[i].filename;
          let fileSize = files[i].size;
          if (file.match('.linux-amd64.tar.gz$')) {
            this.linuxFileName = file;
            this.linuxFileSize = Math.round(fileSize / Math.pow(1024, 2));
          }
          if (file.match('.darwin-amd64(-osx10.8)?.pkg$')) {
            this.macFileName = file;
            this.macFileSize = Math.round(fileSize / Math.pow(1024, 2));
          }
          if (file.match('.windows-amd64.msi$')) {
            this.windowsFileName = file;
            this.windowsFileSize = Math.round(fileSize / Math.pow(1024, 2));
          }
        }
        this.detectOS();
        const osTab = document.getElementById(this.osName);
        if (osTab !== null) {
          osTab.click();
        }
        this.setDownloadForOS(this.osName);
      })
      .catch(console.error);
      this.setEventListeners();
  }

  setEventListeners() {
    this.tabs.forEach((tabEl) => {
      tabEl.addEventListener('click', e => this.handleTabClick((e)));
    });
  }

  // Set the download button UI for a specific OS.
  setDownloadForOS(osName) {
    const baseURL = '/dl/';
    let download;

    switch(osName){
      case 'linux':
        document.querySelector('.js-downloadButton').textContent =
          'Download Go for Linux';
        document.querySelector('.js-downloadDescription').textContent =
          this.linuxFileName + ' (' + this.linuxFileSize + ' MB)';
        document.querySelector('.js-download').href = baseURL + this.linuxFileName;
        break;
      case 'mac':
        document.querySelector('.js-downloadButton').textContent =
          'Download Go for Mac';
        document.querySelector('.js-downloadDescription').textContent =
          this.macFileName + ' (' + this.macFileSize + ' MB)';
        document.querySelector('.js-download').href = baseURL + this.macFileName;
        break;
      case 'windows':
        document.querySelector('.js-downloadButton').textContent =
          'Download Go for Windows';
        document.querySelector('.js-downloadDescription').textContent =
          this.windowsFileName + ' (' + this.windowsFileSize + ' MB)';
        document.querySelector('.js-download').href = baseURL + this.windowsFileName;
        break;
      default:
        document.querySelector('.js-downloadButton').textContent = 'Download Go';
        document.querySelector('.js-downloadDescription').textContent =
          'Visit the downloads page.';
        document.querySelector('.js-download').href = baseURL;
        break;
    }
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
    this.setDownloadForOS(el.id);
    this.setInstallTabData(el.id);
  }

}

// Instantiate controller for page event handling.
new DownloadsController();
