$(document).ready(function() {
    $('#ArticleTextArea').on('input', function () {
        this.style.height = 'auto';
        this.style.height = (this.scrollHeight) + 'px';
    });
});