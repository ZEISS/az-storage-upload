function chooseBinary() {
  if (platform === 'linux' && arch === 'x64') {
    return `main-amd64-linux`
  }

  if (platform === 'linux' && arch === 'arm64') {
    return `main-arm64-linux`
  }

  if (platform === 'darwin' && arch === 'arm64') {
    return `main-arm64-darwin`
  }

  if (platform === 'darwin' && arch === 'x64') {
    return `main-amd64-darwin`
  }
}

const os = require('os');
const childProcess = require('child_process');
const resolve = require('path').resolve
const platform = os.platform()
const arch = os.arch()
const binary = chooseBinary()
const mainScript = resolve(`${__dirname}/../dist/${binary}`)
const proc = childProcess.spawnSync(mainScript, { stdio: 'inherit' })
process.exit(proc.status)
