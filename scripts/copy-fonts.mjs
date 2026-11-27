import { copyFileSync, mkdirSync } from 'node:fs';
import { join } from 'node:path';

const target = 'dist/assets/fonts';
mkdirSync(target, { recursive: true });
for (const family of ['geist', 'jetbrains-mono']) {
  const source = join('node_modules', '@fontsource-variable', family);
  for (const subset of ['latin', 'latin-ext']) {
    const name = `${family}-${subset}-wght-normal.woff2`;
    copyFileSync(join(source, 'files', name), join(target, name));
  }
  copyFileSync(join(source, 'LICENSE'), join(target, `${family}-LICENSE.txt`));
}
