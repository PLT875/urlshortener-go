const form = document.getElementById('form') as HTMLFormElement;
const input = document.getElementById('url') as HTMLInputElement;
const result = document.getElementById('result') as HTMLDivElement;

form.addEventListener('submit', async (e) => {
    e.preventDefault();
    result.textContent = 'Submitting...';
    try {
        const resp = await fetch('/shorten', {
            method: 'POST',
            headers: { 'Content-Type': 'application/json' },
            body: JSON.stringify({ url: input.value })
        });
        if (resp.ok) {
            const data = await resp.json();
            const short = data.short;
            const loc = `${location.origin}/${short}`;
            result.innerHTML = `Short URL: <a href="${loc}" target="_blank" rel="noopener">${loc}</a>`;
        } else {
            const txt = await resp.text();
            result.textContent = `Error: ${txt}`;
        }
    } catch (err) {
        result.textContent = `Error: ${err}`;
    }
});
