/**
 * This is a basic version of the jump links.
 * New version will follow with parity to package site.
 */
(function() {
    'use-strict'
    async function resetNav() {
        return new Promise((resolve, reject) => {
            const nav = document.querySelector('.LearnNav');
            if (!nav) {
                reject(new Error('.LearnNav element not found.'));
            }
            for (let a of Array.from(nav?.children ?? [])) {
                a.classList.remove('active');
            }
            resolve();
        });
    }
    async function setNav() {
        return new Promise((resolve, reject) => {
            const nav = document.querySelector('.LearnNav');
            if (!nav) {
                reject(new Error('.LearnNav element not found.'));
            }
            for (let a of Array.from(nav?.children ?? [])) {
                if (a.href === location.href) {
                    a.classList.add('active');
                    document.getElementById(location.hash.substring(1))?.scrollIntoView();
                    break;
                }
            }
            resolve();
        });
    }
    function observeSections() {
        const learnContent = document.querySelector('.Learn-body');
        if (learnContent && learnContent.children) {
            const callback = async (payload) => {
                if (payload && payload.length > 0) {
                    for (p of payload) {
                        if (p.isIntersecting) {
                            const {id} = p.target;
                            const link = document.querySelector(`[href='#${id}']`);
                            await resetNav();
                            link.classList.add('active');
                            break;
                        }
                    }
                }
            }
            // rootMargin is important when multiple sections are in the observable area **on page load**.
            // they will still be highlighted on scroll because of the root margin.
            const ob = new IntersectionObserver(callback, {
                threshold: 0,
                rootMargin: '0px 0px -50% 0px',
            });
            for (section of learnContent.children) {
                ob.observe(section);
            }
        }
    }

    setNav()
        .then(observeSections)
        .catch((error) => {
            console.error(error);
        });
   
})();
