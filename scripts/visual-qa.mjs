import { createServer } from 'node:http';
import { readFileSync, mkdirSync, statSync, existsSync } from 'node:fs';
import { resolve, extname, join } from 'node:path';
import { pathToFileURL } from 'node:url';
import { chromium } from '@playwright/test';

const output = resolve('.cache/qa');
const dist = resolve('dist');
const reference = resolve('docs/design-reference.html');
mkdirSync(output, { recursive: true });
const types = { '.html': 'text/html; charset=utf-8', '.css': 'text/css', '.js': 'text/javascript', '.svg': 'image/svg+xml', '.woff2': 'font/woff2' };
const server = createServer((request, response) => {
  const pathname = decodeURIComponent(new URL(request.url, 'http://localhost').pathname);
  let file = resolve(dist, '.' + pathname);
  if (!file.startsWith(dist)) { response.writeHead(403).end(); return; }
  try {
    if (statSync(file).isDirectory()) file = join(file, 'index.html');
    response.writeHead(200, { 'content-type': types[extname(file)] || 'application/octet-stream' });
    response.end(readFileSync(file));
  } catch {
    response.writeHead(404, { 'content-type': types['.html'] });
    response.end(readFileSync(join(dist, '404.html')));
  }
});
await new Promise(resolveListen => server.listen(0, '127.0.0.1', resolveListen));
const base = `http://127.0.0.1:${server.address().port}`;
const browser = await chromium.launch({ headless: true });
try {
  for (const width of [1440, 1024, 768, 390]) {
    const height = width === 390 ? 844 : width === 1024 ? 768 : 900;
    const page = await browser.newPage({ viewport: { width, height }, deviceScaleFactor: 1 });
    if (existsSync(reference)) {
      await page.goto(pathToFileURL(reference).href);
      await page.screenshot({ path: join(output, `reference-${width}.png`), fullPage: true });
    }
    await page.goto(base);
    await page.evaluate(() => document.fonts.ready);
    await page.screenshot({ path: join(output, `site-${width}.png`), fullPage: true });
    const overflow = await page.evaluate(() => document.documentElement.scrollWidth > innerWidth);
    if (overflow) throw new Error(`Horizontal overflow at ${width}px`);
    await page.close();
  }
  const page = await browser.newPage();
  await page.goto(base);
  await page.getByRole('button', { name: 'Alternar tema' }).click();
  if (await page.locator('html').getAttribute('data-theme') !== 'light') throw new Error('Theme toggle failed');
  await page.reload();
  if (await page.locator('html').getAttribute('data-theme') !== 'light') throw new Error('Theme persistence failed');
  await page.keyboard.press('Control+k');
  if (!await page.locator('#palette').evaluate(dialog => dialog.open)) throw new Error('Search dialog did not open');
  await page.locator('#palette-search').fill('Projetos');
  await page.keyboard.press('ArrowDown');
  if (!await page.locator('#palette .palette-list a:focus').count()) throw new Error('Keyboard selection failed');
  await page.keyboard.press('Escape');
  if (await page.locator('#palette').evaluate(dialog => dialog.open)) throw new Error('Escape did not close search');
  await page.getByRole('link', { name: 'en', exact: true }).click();
  if (new URL(page.url()).pathname !== '/en/') throw new Error('Locale switch failed');
  await page.goto(base + '/projects/fobras-next/');
  await page.getByRole('link', { name: 'en', exact: true }).click();
  if (new URL(page.url()).pathname !== '/en/projects/fobras-next/') throw new Error('Project locale switch failed');
  await page.getByRole('link', { name: 'Projects', exact: true }).first().click();
  if (new URL(page.url()).pathname !== '/en/projects/') throw new Error('Secondary navigation failed');
  await page.goto(base + '/missing/');
  if (!await page.getByText('Página não encontrada / Page not found').count()) throw new Error('404 fallback failed');
  await page.close();
  console.log(`Browser checks passed. Screenshots: ${output}`);
} finally {
  await browser.close();
  await new Promise(resolveClose => server.close(resolveClose));
}
