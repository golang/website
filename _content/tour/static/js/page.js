window.transport = {{.Transport}}();
window.socketAddr = "{{.SocketAddr}}";

function highlight(selector) {
    var speed = 50;
    var obj = $(selector).stop(true, true)
    for (var i = 0; i < 5; i++) {
        obj.addClass("highlight", speed)
        obj.delay(speed)
        obj.removeClass("highlight", speed)
    }
}

function highlightAndClick(selector) {
    highlight(selector);
    setTimeout(function() {
        $(selector)[0].click()
    }, 750);
}

function click(selector) {
    $(selector)[0].click();
}
