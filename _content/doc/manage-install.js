document.addEventListener('DOMContentLoaded', () => {
    const tablist = document.querySelector('.TabSection-tabList');
    if (!tablist) return;
    const tabs = tablist.querySelectorAll('[role="tab"]');
    const panels = document.querySelectorAll('[role="tabpanel"]');

    tabs.forEach((tab) => {
        tab.addEventListener('click', () => {
            tabs.forEach(t => {
                t.setAttribute('aria-selected', 'false');
                t.setAttribute('tabindex', '-1');
                t.classList.remove('active');
            });
            panels.forEach(p => p.setAttribute('hidden', ''));
            
            tab.setAttribute('aria-selected', 'true');
            tab.setAttribute('tabindex', '0');
            tab.classList.add('active');
            
            const controlsId = tab.getAttribute('aria-controls');
            const panel = document.getElementById(controlsId);
            if (panel) {
                panel.removeAttribute('hidden');
            }
        });
    });
});
