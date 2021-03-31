import { promises } from 'fs';
import { join, normalize, isAbsolute } from 'path';
import { Command } from '../commander';
import { ExecFile, ClearDirectory, RmDirectory } from './utils';

class GRPC {
    names = new Array<string>()
    constructor(public readonly builders: Array<Builder>) {
        builders.forEach((builder) => {
            this.names.push(builder.name)
        })
    }
    async build(language: string, output: string, includes: Array<string>): Promise<void> {
        const builders = this.builders
        for (let i = 0; i < builders.length; i++) {
            const builder = builders[i]
            if (language === builder.name) {
                return builder.build(output, includes)
            }
        }
        const e = new Error(`not supported language : ${language}`)
        console.warn(e)
        throw e
    }
}
class WalkResult {
    files = new Array<string>()
    dirs = new Array<string>()
}
async function Walk(root: string, dir: string): Promise<WalkResult> {
    const result = new WalkResult()
    const names = await promises.readdir(join(root, dir))
    for (let i = 0; i < names.length; i++) {
        const name = names[i]
        const filename = join(root, dir, name)
        const stat = await promises.stat(filename)
        if (stat.isDirectory()) {
            result.dirs.push(join(dir, name))
        } else if (stat.isFile()) {
            result.files.push(join(dir, name))
        }
    }
    return result
}
class Builder {
    private cwd_ = ''
    get cwd(): string {
        return this.cwd_
    }
    private root_ = ''
    get root(): string {
        return this.root_
    }
    private first_ = true
    private include_: Promise<Array<string>> = null
    getInclude(includes: Array<string>): Promise<Array<string>> {
        if (!this.include_) {
            this.include_ = new Promise<Array<string>>(async (resolve, reject) => {
                try {
                    const set = new Set<string>()
                    const args = new Array<string>()
                    args.push('-I', this.root)
                    set.add(this.root)
                    if (Array.isArray(includes)) {
                        includes.forEach((v) => {
                            if (!set.has(v)) {
                                args.push('-I', v)
                                set.add(v)
                            }
                        })
                    }
                    resolve(args)
                } catch (e) {
                    console.warn(e)
                    reject(e)
                }
            })
        }
        return this.include_
    }
    constructor(public readonly name: string,
        public readonly uuid: string,
        public readonly gateway: boolean,
    ) {
        this.cwd_ = normalize(join(__dirname, '..', '..'))
        const root = join('pb', this.uuid)
        this.root_ = normalize(join(this.cwd_, root))
    }
    async build(output: string, includes: Array<string>): Promise<void> {
        output = await this.getOutput(output)
        this.first_ = true
        await this._build(output, '.', includes)
        await this.done()
    }
    async _build(output: string, dir: string, includes: Array<string>): Promise<void> {
        const result = await Walk(this.root, dir)
        if (result.files.length > 0) {
            if (this.first_) {
                this.first_ = false
            } else {
                console.log()
            }
            console.log(`------ build : ${dir} ------`)
            await this.buildGRPC(output, includes, ...result.files)
        }
        const dirs = result.dirs
        for (let i = 0; i < dirs.length; i++) {
            await this._build(output, dirs[i], includes)
        }
    }
    async getOutput(output: string): Promise<string> {
        throw Error('getOutput grpc not impl')
    }
    async buildGRPC(output: string, includes: Array<string>, ...files: Array<string>) {
        throw Error('build grpc not impl')
    }
    async done() {

    }
}
class Dart extends Builder {
    constructor(uuid: string, gateway: boolean) {
        super('dart', uuid, gateway)
    }
    async getOutput(output: string): Promise<string> {
        if (typeof output != 'string' || output.length == 0) {
            output = join('bin', 'protocol', 'dart')
        }
        try {
            let filename: string
            if (isAbsolute(output)) {
                filename = output
            } else {
                filename = normalize(join(this.cwd, output))
            }
            await ClearDirectory(filename)
        } catch (e) {
            console.warn(e)
        }
        return output
    }
    async buildGRPC(output: string, includes: Array<string>, ...files: Array<string>) {
        const args = [
            ...await this.getInclude(includes),
            `--dart_out=grpc:${output}`,
            ...files,
        ]
        console.log(`protoc ${args.join(' ')}`)
        await ExecFile('protoc', args, {
            cwd: this.cwd,
        })
    }
}
class Go extends Builder {
    constructor(public readonly pkg: string,
        uuid: string, gateway: boolean,
        public readonly gin: boolean,
    ) {
        super('go', uuid, gateway)
    }
    async getOutput(output: string): Promise<string> {
        output = join(this.cwd, '.tmp')
        try {
            await ClearDirectory(output)
            if (this.gin) {
                await RmDirectory(join(this.cwd, 'static', 'document', 'api'))
                await ClearDirectory(join(output, 'api'))
            }
            await RmDirectory(join(this.cwd, 'protocol'))
        } catch (e) {
            console.warn(e)
        }
        return output
    }
    async buildGRPC(output: string, includes: Array<string>, ...files: Array<string>) {
        let args = [
            ...await this.getInclude(includes),
            `--go_out=plugins=grpc:${output}`,
            ...files,
        ]
        console.log(`protoc ${args.join(' ')}`)
        await ExecFile('protoc', args, {
            cwd: this.cwd,
        })
        if (this.gateway) {
            args = [
                ...await this.getInclude(includes),
                `--grpc-gateway_out=logtostderr=true:${output}`,
                ...files,
            ]
            console.log(`protoc ${args.join(' ')}`)
            await ExecFile('protoc', args, {
                cwd: this.cwd,
            })
        }

        if (this.gin) {
            args = [
                ...await this.getInclude(includes),
                `--openapiv2_out=logtostderr=true:${join(output, 'api')}`,
                ...files,
            ]
            console.log(`protoc ${args.join(' ')}`)
            await ExecFile('protoc', args, {
                cwd: this.cwd,
            })
        }
    }
    async done() {
        await promises.rename(normalize(join(this.cwd, '.tmp', this.pkg, 'protocol')), join(this.cwd, 'protocol'))
        if (this.gin) {
            await promises.rename(normalize(join(this.cwd, '.tmp', 'api')), join(this.cwd, 'static', 'document', 'api'))
        }
        await RmDirectory(join(this.cwd, '.tmp'))
    }
}
export function BuildGRPC(program: Command, pkg: string, uuid: string, gateway: boolean, gin: boolean) {
    const grpc = new GRPC([
        new Go(pkg, uuid, gateway, gin),
        new Dart(uuid, gateway),
    ])
    program.command('grpc')
        .description('build *.proto to grpc code')
        .option(`-l,--language [${grpc.names.join(' ')}]`, 'grpc target language')
        .option('-o,--output []', 'grpc output directory')
        .option('-i, --include [includes...]', 'protoc include path')
        .action(function () {
            const opts = this.opts()
            grpc.build(opts['language'], opts['output'], opts["include"]).catch(() => {
                process.exit(1)
            })
        })
}