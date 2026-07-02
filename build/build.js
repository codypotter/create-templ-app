import * as esbuild from 'esbuild'
import { execSync, spawn } from 'child_process'
import { createHash } from 'crypto'
import { readFileSync, writeFileSync, mkdirSync, copyFileSync, rmSync } from 'fs'
import { basename } from 'path'

const watch = process.argv.includes('--watch')
const outdir = 'internal/assets/dist'
const manifest = {}

// Resolve the tailwindcss CLI entry point directly — avoids symlink issues in npm scripts.
const twCLI = new URL('../node_modules/@tailwindcss/cli/dist/index.mjs', import.meta.url).pathname

// Start clean so stale hashed files from previous builds don't linger.
rmSync(outdir, { recursive: true, force: true })
mkdirSync(outdir, { recursive: true })

// --- CSS via tailwindcss CLI ---

function buildCSS() {
  const tmpOut = `${outdir}/_main.css`
  execSync(
    `${process.execPath} ${twCLI} -i internal/assets/src/css/main.css -o ${tmpOut}${watch ? '' : ' --minify'}`,
    { stdio: 'inherit' }
  )
  const css = readFileSync(tmpOut)
  const hash = createHash('sha256').update(css).digest('hex').slice(0, 8)
  const hashedName = `main-${hash}.css`
  copyFileSync(tmpOut, `${outdir}/${hashedName}`)
  rmSync(tmpOut)
  manifest['main.css'] = hashedName
}

// --- JS via esbuild ---

function buildJS() {
  const result = esbuild.buildSync({
    entryPoints: ['internal/assets/src/js/main.js'],
    bundle: true,
    minify: !watch,
    entryNames: '[name]-[hash]',
    outdir,
    metafile: true,
  })
  for (const [outputPath, output] of Object.entries(result.metafile.outputs)) {
    if (output.entryPoint) {
      manifest[basename(output.entryPoint)] = basename(outputPath)
    }
  }
}

function writeManifest() {
  writeFileSync(`${outdir}/manifest.json`, JSON.stringify(manifest, null, 2))
  console.log('manifest:', manifest)
}

if (watch) {
  console.log('starting watch mode...')
  buildCSS()
  buildJS()
  writeManifest()

  // Tailwind in watch mode — spawns node directly with the resolved CLI path.
  const tw = spawn(
    process.execPath,
    [twCLI, '-i', 'internal/assets/src/css/main.css', '-o', `${outdir}/main.css`, '--watch'],
    { stdio: 'inherit' }
  )
  tw.on('exit', code => process.exit(code ?? 0))

  // esbuild JS watch
  const ctx = await esbuild.context({
    entryPoints: ['internal/assets/src/js/main.js'],
    bundle: true,
    outdir,
    metafile: false,
  })
  await ctx.watch()
} else {
  buildCSS()
  buildJS()
  writeManifest()
}
