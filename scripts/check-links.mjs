import { readFileSync, readdirSync, statSync, existsSync } from 'node:fs';
import { join, relative, resolve } from 'node:path';

const root = resolve('dist');
const files = [];
function walk(dir) {
  for (const name of readdirSync(dir)) {
    const file = join(dir, name);
    if (statSync(file).isDirectory()) walk(file);
    else if (file.endsWith('.html')) files.push(file);
  }
}
walk(root);
const errors = [];
for (const file of files) {
  const html = readFileSync(file, 'utf8');
  const rel = relative(root, file).replaceAll('\\', '/');
  const route = '/' + rel.replace(/index\.html$/, '');
  for (const [, href] of html.matchAll(/href="([^"]+)"/g)) {
    if (href.startsWith('mailto:') || href.startsWith('https:')) continue;
    const url = new URL(href, `https://local.invalid${route}`);
    if (url.origin !== 'https://local.invalid') continue;
    const target = resolve(root, '.' + url.pathname, url.pathname.endsWith('/') ? 'index.html' : '');
    if (!target.startsWith(root) || !existsSync(target)) {
      errors.push(`${rel}: missing ${href}`);
      continue;
    }
    if (url.hash && target.endsWith('.html')) {
      const id = decodeURIComponent(url.hash.slice(1));
      if (!readFileSync(target, 'utf8').includes(`id="${id}"`)) errors.push(`${rel}: missing anchor ${href}`);
    }
  }
}
if (errors.length) {
  console.error(errors.join('\n'));
  process.exitCode = 1;
} else console.log(`Checked internal links in ${files.length} HTML files.`);
