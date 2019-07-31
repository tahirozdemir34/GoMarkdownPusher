function compileMarkdown() {
    var text = document.getElementById('markdown').value,
        target = document.getElementById('compiledMarkdown'),
        converter = new showdown.Converter(),
        html = converter.makeHtml(text);

    target.innerHTML = html;
}