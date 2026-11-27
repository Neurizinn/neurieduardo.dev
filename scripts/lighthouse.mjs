import { createServer } from 'node:http';
import { readFileSync, mkdirSync, statSync } from 'node:fs';
import { resolve, extname, join } from 'node:path';
import { spawn } from 'node:child_process';
import { chromium } from '@playwright/test';

const dist = resolve('dist');
const output = resolve('.cache/qa/lighthouse.json');
mkdirSync(resolve('.cache/qa'), { recursive: true });
const types = { '.html': 'text/html; charset=utf-8', '.css': 'text/css', '.js': 'text/javascript', '.svg': 'image/svg+xml', '.woff2': 'font/woff2' };
const server = createServer((request, response) => {
  const pathname = decodeURIComponent(new URL(request.url, 'http://localhost').pathname);
  let file = resolve(dist, '.' + pathname);
  if (!file.startsWith(dist)) { response.writeHead(403).end(); return; }
  try {
    if (statSync(file).isDirectory()) file = join(file, 'index.html');
    response.writeHead(200, { 'content-type': types[extname(file)] || 'application/octet-stream' });
    response.end(readFileSync(file));
  } catch { response.writeHead(404).end(); }
});
await new Promise(resolveListen => server.listen(0, '127.0.0.1', resolveListen));
try {
  const url = `http://127.0.0.1:${server.address().port}/`;
  const args = [resolve('node_modules/lighthouse/cli/index.js'), url, '--output=json', `--output-path=${output}`, '--quiet', '--chrome-flags=--headless=new --disable-gpu --no-sandbox'];
  const child = spawn(process.execPath, args, { env: { ...process.env, CHROME_PATH: chromium.executablePath() }, stdio: ['ignore', 'pipe', 'pipe'] });
  let log = '';
  child.stdout.on('data', chunk => { log += chunk; });
  child.stderr.on('data', chunk => { log += chunk; });
  const code = await new Promise(resolveExit => child.on('close', resolveExit));
  if (code !== 0) throw new Error(`Lighthouse exited ${code}: ${log.slice(-2000)}`);
  const report = JSON.parse(readFileSync(output, 'utf8'));
  for (const key of ['performance', 'accessibility', 'best-practices', 'seo']) {
    console.log(`${key}: ${Math.round(report.categories[key].score * 100)}`);
  }
  console.log(`Report: ${output}`);
} finally {
  await new Promise(resolveClose => server.close(resolveClose));
}
