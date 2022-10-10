(() => {
    'use strict';
    const copyButtons = document.querySelectorAll('.CopyPaste button');
    function copyToClipboard(copyText) {
        if (typeof navigator?.clipboard?.writeText !== 'function') return;
        navigator.clipboard.writeText(copyText);
    }
    for (let btn of copyButtons) {
        btn.addEventListener('click', () => {
            const content = btn?.previousElementSibling?.textContent ?? '';
            const text = content.substring(content?.[0] === '$' ? 1 : 0);
            copyToClipboard(text);
        });
    }
})();
