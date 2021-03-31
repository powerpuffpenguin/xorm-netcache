"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BuildGRPC = void 0;
const fs_1 = require("fs");
const path_1 = require("path");
const utils_1 = require("./utils");
class GRPC {
    constructor(builders) {
        this.builders = builders;
        this.names = new Array();
        builders.forEach((builder) => {
            this.names.push(builder.name);
        });
    }
    async build(language, output, includes) {
        const builders = this.builders;
        for (let i = 0; i < builders.length; i++) {
            const builder = builders[i];
            if (language === builder.name) {
                return builder.build(output, includes);
            }
        }
        const e = new Error(`not supported language : ${language}`);
        console.warn(e);
        throw e;
    }
}
class WalkResult {
    constructor() {
        this.files = new Array();
        this.dirs = new Array();
    }
}
async function Walk(root, dir) {
    const result = new WalkResult();
    const names = await fs_1.promises.readdir(path_1.join(root, dir));
    for (let i = 0; i < names.length; i++) {
        const name = names[i];
        const filename = path_1.join(root, dir, name);
        const stat = await fs_1.promises.stat(filename);
        if (stat.isDirectory()) {
            result.dirs.push(path_1.join(dir, name));
        }
        else if (stat.isFile()) {
            result.files.push(path_1.join(dir, name));
        }
    }
    return result;
}
class Builder {
    constructor(name, uuid, gateway) {
        this.name = name;
        this.uuid = uuid;
        this.gateway = gateway;
        this.cwd_ = '';
        this.root_ = '';
        this.first_ = true;
        this.include_ = null;
        this.cwd_ = path_1.normalize(path_1.join(__dirname, '..', '..'));
        const root = path_1.join('pb', this.uuid);
        this.root_ = path_1.normalize(path_1.join(this.cwd_, root));
    }
    get cwd() {
        return this.cwd_;
    }
    get root() {
        return this.root_;
    }
    getInclude(includes) {
        if (!this.include_) {
            this.include_ = new Promise(async (resolve, reject) => {
                try {
                    const set = new Set();
                    const args = new Array();
                    args.push('-I', this.root);
                    set.add(this.root);
                    if (Array.isArray(includes)) {
                        includes.forEach((v) => {
                            if (!set.has(v)) {
                                args.push('-I', v);
                                set.add(v);
                            }
                        });
                    }
                    resolve(args);
                }
                catch (e) {
                    console.warn(e);
                    reject(e);
                }
            });
        }
        return this.include_;
    }
    async build(output, includes) {
        output = await this.getOutput(output);
        this.first_ = true;
        await this._build(output, '.', includes);
        await this.done();
    }
    async _build(output, dir, includes) {
        const result = await Walk(this.root, dir);
        if (result.files.length > 0) {
            if (this.first_) {
                this.first_ = false;
            }
            else {
                console.log();
            }
            console.log(`------ build : ${dir} ------`);
            await this.buildGRPC(output, includes, ...result.files);
        }
        const dirs = result.dirs;
        for (let i = 0; i < dirs.length; i++) {
            await this._build(output, dirs[i], includes);
        }
    }
    async getOutput(output) {
        throw Error('getOutput grpc not impl');
    }
    async buildGRPC(output, includes, ...files) {
        throw Error('build grpc not impl');
    }
    async done() {
    }
}
class Dart extends Builder {
    constructor(uuid, gateway) {
        super('dart', uuid, gateway);
    }
    async getOutput(output) {
        if (typeof output != 'string' || output.length == 0) {
            output = path_1.join('bin', 'protocol', 'dart');
        }
        try {
            let filename;
            if (path_1.isAbsolute(output)) {
                filename = output;
            }
            else {
                filename = path_1.normalize(path_1.join(this.cwd, output));
            }
            await utils_1.ClearDirectory(filename);
        }
        catch (e) {
            console.warn(e);
        }
        return output;
    }
    async buildGRPC(output, includes, ...files) {
        const args = [
            ...await this.getInclude(includes),
            `--dart_out=grpc:${output}`,
            ...files,
        ];
        console.log(`protoc ${args.join(' ')}`);
        await utils_1.ExecFile('protoc', args, {
            cwd: this.cwd,
        });
    }
}
class Go extends Builder {
    constructor(pkg, uuid, gateway, gin) {
        super('go', uuid, gateway);
        this.pkg = pkg;
        this.gin = gin;
    }
    async getOutput(output) {
        output = path_1.join(this.cwd, '.tmp');
        try {
            await utils_1.ClearDirectory(output);
            if (this.gin) {
                await utils_1.RmDirectory(path_1.join(this.cwd, 'static', 'document', 'api'));
                await utils_1.ClearDirectory(path_1.join(output, 'api'));
            }
            await utils_1.RmDirectory(path_1.join(this.cwd, 'protocol'));
        }
        catch (e) {
            console.warn(e);
        }
        return output;
    }
    async buildGRPC(output, includes, ...files) {
        let args = [
            ...await this.getInclude(includes),
            `--go_out=plugins=grpc:${output}`,
            ...files,
        ];
        console.log(`protoc ${args.join(' ')}`);
        await utils_1.ExecFile('protoc', args, {
            cwd: this.cwd,
        });
        if (this.gateway) {
            args = [
                ...await this.getInclude(includes),
                `--grpc-gateway_out=logtostderr=true:${output}`,
                ...files,
            ];
            console.log(`protoc ${args.join(' ')}`);
            await utils_1.ExecFile('protoc', args, {
                cwd: this.cwd,
            });
        }
        if (this.gin) {
            args = [
                ...await this.getInclude(includes),
                `--openapiv2_out=logtostderr=true:${path_1.join(output, 'api')}`,
                ...files,
            ];
            console.log(`protoc ${args.join(' ')}`);
            await utils_1.ExecFile('protoc', args, {
                cwd: this.cwd,
            });
        }
    }
    async done() {
        await fs_1.promises.rename(path_1.normalize(path_1.join(this.cwd, '.tmp', this.pkg, 'protocol')), path_1.join(this.cwd, 'protocol'));
        if (this.gin) {
            await fs_1.promises.rename(path_1.normalize(path_1.join(this.cwd, '.tmp', 'api')), path_1.join(this.cwd, 'static', 'document', 'api'));
        }
        await utils_1.RmDirectory(path_1.join(this.cwd, '.tmp'));
    }
}
function BuildGRPC(program, pkg, uuid, gateway, gin) {
    const grpc = new GRPC([
        new Go(pkg, uuid, gateway, gin),
        new Dart(uuid, gateway),
    ]);
    program.command('grpc')
        .description('build *.proto to grpc code')
        .option(`-l,--language [${grpc.names.join(' ')}]`, 'grpc target language')
        .option('-o,--output []', 'grpc output directory')
        .option('-i, --include [includes...]', 'protoc include path')
        .action(function () {
        const opts = this.opts();
        grpc.build(opts['language'], opts['output'], opts["include"]).catch(() => {
            process.exit(1);
        });
    });
}
exports.BuildGRPC = BuildGRPC;
