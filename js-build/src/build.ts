import { promises } from 'fs';
import { join, normalize } from 'path';
import { Append, Env, ExecFile, Exec } from './utils';
import { Command } from '../commander'
class Target {
    private name_ = ''
    constructor(public readonly os: string, public readonly arch: string,
        public readonly name: string, public readonly ext: string,
        public readonly debug: boolean) {
        if (debug) {
            this.name_ = name + 'd'
        } else {
            this.name_ = name
        }
        if (typeof this.ext === "string") {
            this.name_ += this.ext
        }
    }
    build() {
        const name = join('bin', this.name_)
        console.log(`go build -o "${name}"`)
        const env = Env({
            GOOS: this.os,
            GOARCH: this.arch,
        })
        return ExecFile(`go`, ['build', '-o', name], {
            cwd: normalize(join(__dirname, '..', '..')),
            env: env,
        })
    }
    pack(algorithm: string, ...names: Array<string>) {
        if (typeof algorithm === "string") {
            let output = `${this.os}-${this.arch}`
            let file
            const args = []
            switch (algorithm) {
                case '7z':
                    file = '7z'
                    args.push('a')
                    output += `.7z`
                    break;
                case 'zip':
                    file = 'zip'
                    args.push('-r')
                    output += `.zip`
                    break;
                case 'gz':
                    file = 'tar'
                    args.push('-zcvf')
                    output += `.tar.gz`
                    break
                case 'bz2':
                    file = 'tar'
                    args.push('-jcvf')
                    output += `.tar.bz2`
                    break
                case 'xz':
                    file = 'tar'
                    args.push('-Jcvf')
                    output += `.tar.xz`
                    break
                default:
                    throw new Error(`not supported pack algorithm : ${algorithm}`)
            }
            args.push(output)
            args.push(this.name_)
            for (let i = 0; i < names.length; i++) {
                args.push(names[i])
            }
            console.log(file, ...args)
            return ExecFile(file, args, {
                cwd: normalize(join(__dirname, '..', '..', 'bin')),
            })
        }
    }
}
export function Build(program: Command,
    os: string, arch: Array<string>,
    name: string, ext: string,
    ...packs: Array<string>) {
    const pack = '7z gz bz2 xz zip'
    program.command(os)
        .description(`build code to ${os}`)
        .option(`--arch [${arch.join(' ')}]`, 'GOARCH default use amd64')
        .option(`-p,--pack [${pack}]`, 'Pack to compressed package')
        .option('--debug', 'build as debug')
        .action(function () {
            const opts = this.opts()
            if (opts['Version']) {
                buildVersion()
            } else {
                let arch = opts['arch']
                if (arch === undefined) {
                    arch = 'amd64'
                }
                const taget = new Target(os, arch, name, ext, opts['debug'])
                taget.build().then(() => {
                    return taget.pack(opts['pack'],
                        ...packs,
                    )
                }).catch(() => {
                    process.exit(1)
                })
            }
        })
}
async function test(pkg: string, args: Array<string>, ...names: Array<string>) {
    const cwd = normalize(join(__dirname, '..', '..'))
    for (let i = 0; i < names.length; i++) {
        const name = names[i]
        if (name.length > 0) {
            await ExecFile('go', Append(args, `${pkg}/${name}`), {
                cwd: cwd,
            })
        }
    }
}
export function BuildTest(program: Command, pkg: string, unit: Array<string>, bench: Array<string>) {
    program.command('test')
        .description(`run go test`)
        .option(`-v`, 'prints the full')
        .option(`-run []`, 'run func name')
        .option(`-b, --bench`, 'test bench')
        .action(function () {
            const opts = this.opts()
            const args = ['test']
            if (opts['v']) {
                args.push('-v')
            }
            const run = opts['Run']
            let names = bench
            if (opts['bench']) {
                if (typeof run === "string") {
                    args.push(`-run=^$`)
                }
                args.push('-benchmem', '-bench')
                if (typeof run === "string") {
                    args.push(`${run}`)
                } else {
                    args.push('.')
                }

            } else {
                if (typeof run === "string") {
                    args.push(`-run=${run}`)
                }
                names = unit
            }
            test(pkg, args, ...names).catch(() => {
                process.exit(1)
            })
        })
}

async function buildVersion() {
    let tag = ''
    let commit = ''
    try {
        tag = await Exec('git', ['describe'])
    } catch (e) {

    }
    try {
        commit = await Exec('git', ['rev-parse', 'HEAD'])
    } catch (e) {

    }
    const date = new Date().toUTCString()
    const filename = normalize(join(__dirname, '..', '..', 'version', 'version.go'))
    const str = ['package version', '',
        '// Tag git tag', 'const Tag = `' + tag.trim() + '`', '',
        '// Commit git commit', 'const Commit = `' + commit.trim() + '`', '',
        '// Date build datetime', 'const Date = `' + date + '`', '',
    ].join("\r\n")
    console.log(str.trim())
    await promises.writeFile(filename, str)
}
export function BuildVersion(program: Command) {
    program.command('version')
        .description('update ' + join('version', 'version.go'))
        .action(function () {
            buildVersion().catch((e) => {
                console.warn(e)
                process.exit(1)
            })
        })
}