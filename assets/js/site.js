(() => {
  const root = document.documentElement;
  try {
    const saved = localStorage.getItem('theme');
    if (saved === 'light' || saved === 'dark') root.dataset.theme = saved;
  } catch {}

  document.getElementById('theme')?.addEventListener('click', () => {
    const next = root.dataset.theme === 'light' ? 'dark' : 'light';
    root.dataset.theme = next;
    try { localStorage.setItem('theme', next); } catch {}
  });

  const dialog = document.getElementById('palette');
  const input = document.getElementById('palette-search');
  const links = [...dialog.querySelectorAll('.palette-list a')];
  const visible = () => links.filter(link => !link.hidden);
  const open = () => {
    if (!dialog.open) dialog.showModal();
    input.value = '';
    links.forEach(link => { link.hidden = false; });
    input.focus();
  };
  document.getElementById('palette-open')?.addEventListener('click', open);
  document.getElementById('palette-close')?.addEventListener('click', () => dialog.close());
  input.addEventListener('input', () => {
    const query = input.value.trim().toLocaleLowerCase();
    links.forEach(link => { link.hidden = !link.textContent.toLocaleLowerCase().includes(query); });
  });
  document.addEventListener('keydown', event => {
    if ((event.ctrlKey || event.metaKey) && event.key.toLowerCase() === 'k') {
      event.preventDefault();
      open();
    }
  });
  dialog.addEventListener('keydown', event => {
    if (event.key === 'Enter' && document.activeElement === input && visible()[0]) {
      visible()[0].click();
    }
    if (event.key !== 'ArrowDown' && event.key !== 'ArrowUp') return;
    const items = visible();
    if (!items.length) return;
    event.preventDefault();
    const current = items.indexOf(document.activeElement);
    const next = event.key === 'ArrowDown' ? (current + 1) % items.length : (current - 1 + items.length) % items.length;
    items[next].focus();
  });
})();
