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
const platform = os.platform()
const arch = os.arch()
const binary = chooseBinary()
const mainScript = `${__dirname}/dist/${binary}`
const spawnSyncReturns = childProcess.spawnSync(mainScript, { stdio: 'inherit' })
