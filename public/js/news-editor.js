(() => {
    const title = document.getElementById('news-title');
    const slug = document.getElementById('news-slug');
    if (!title || !slug || slug.readOnly) return;
    let automatic = !slug.value;
    slug.addEventListener('input', () => { automatic = !slug.value; });
    title.addEventListener('input', () => {
        if (automatic) slug.value = title.value.toLowerCase().normalize('NFD')
            .replace(/[\u0300-\u036f]/g, '').replace(/[^a-z0-9]+/g, '-')
            .replace(/^-|-$/g, '').slice(0, 180).replace(/-$/, '');
    });

    const body = document.getElementById('news-body');
    if (!body) return;

    const wrap = (before, after, fallback) => {
        const start = body.selectionStart;
        const end = body.selectionEnd;
        const selected = body.value.slice(start, end) || fallback;
        body.setRangeText(before + selected + after, start, end, 'end');
        body.focus();
    };
    const prefixLines = (prefix) => {
        const start = body.value.lastIndexOf('\n', Math.max(0, body.selectionStart - 1)) + 1;
        const nextBreak = body.value.indexOf('\n', body.selectionEnd);
        const end = nextBreak === -1 ? body.value.length : nextBreak;
        const selected = body.value.slice(start, end) || 'Text';
        body.setRangeText(selected.split('\n').map((line, index) => typeof prefix === 'function' ? prefix(index) + line : prefix + line).join('\n'), start, end, 'select');
        body.focus();
    };

    document.querySelectorAll('[data-markdown]').forEach((button) => {
        button.addEventListener('click', () => {
            switch (button.dataset.markdown) {
                case 'heading': prefixLines('## '); break;
                case 'bold': wrap('**', '**', 'tučný text'); break;
                case 'italic': wrap('*', '*', 'text kurzívou'); break;
                case 'bullet': prefixLines('- '); break;
                case 'number': prefixLines((index) => `${index + 1}. `); break;
                case 'quote': prefixLines('> '); break;
                case 'link': wrap('[', '](https://)', 'text odkazu'); break;
            }
        });
    });
})();
