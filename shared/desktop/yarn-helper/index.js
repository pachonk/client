// @flow
// Helper for cross platform yarn run script commands
import buildCommands from './build'
import electronComands from './electron'
import fontCommands from './font'
import visdiff from './visdiff'
import {execSync} from 'child_process'
import {padEnd} from 'lodash'

const [, , command, ...rest] = process.argv

const commands = {
  ...buildCommands,
  ...fontCommands,
  ...electronComands,
  ...visdiff,
  help: {
    code: () => {
      const len = Object.keys(commands).reduce((acc, i) => Math.max(i.length, acc), 1)
      console.log(
        Object.keys(commands)
          .map(c => commands[c].help && `yarn run ${padEnd(c + ': ', len + 5)}${commands[c].help || ''}`)
          .filter(Boolean)
          .join('\n')
      )
    },
  },
  postinstall: {
    code: () => {
      // Inject dummy module
      exec('node make-shim net')
      exec('node make-shim tls')
      exec('node make-shim msgpack')
    },
    help: 'all: install global eslint. dummy modules',
  },
}

function exec(command, env, options) {
  console.log(execSync(command, {encoding: 'utf8', env: env || process.env, stdio: 'inherit', ...options}))
}

const decorateInfo = info => {
  let temp = {
    ...info,
    env: {
      ...process.env,
      ...info.env,
    },
  }

  if (info.nodeEnv) {
    temp.env.NODE_ENV = info.nodeEnv
  }

  if (rest.length && temp.shell) {
    temp.shell = `${temp.shell} ${rest.join(' ')}`
  }

  return temp
}

function main() {
  let info = commands[command]

  if (!info) {
    console.log('Unknown command: ', command)
    process.exit(1)
  }

  info = decorateInfo(info)

  if (info.shell) {
    exec(info.shell, info.env, info.options)
  }

  if (info.code) {
    info.code()
  }
}

main()
