"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.BuildVersion = exports.BuildTest = exports.Build = void 0;
const fs_1 = require("fs");
const path_1 = require("path");
const utils_1 = require("./utils");
class Target {
    constructor(os, arch, name, ext, debug) {
        this.os = os;
        this.arch = arch;
        this.name = name;
        this.ext = ext;
        this.debug = debug;
        this.name_ = '';
        if (debug) {
            this.name_ = name + 'd';
        }
        else {
            this.name_ = name;
        }
        if (typeof this.ext === "string") {
            this.name_ += this.ext;
        }
    }
    build() {
        const name = path_1.join('bin', this.name_);
        console.log(`go build -o "${name}"`);
        const env = utils_1.Env({
            GOOS: this.os,
            GOARCH: this.arch,
        });
        return utils_1.ExecFile(`go`, ['build', '-o', name], {
            cwd: path_1.normalize(path_1.join(__dirname, '..', '..')),
            env: env,
        });
    }
    pack(algorithm, ...names) {
        if (typeof algorithm === "string") {
            let output = `${this.os}-${this.arch}`;
            let file;
            const args = [];
            switch (algorithm) {
                case '7z':
                    file = '7z';
                    args.push('a');
                    output += `.7z`;
                    break;
                case 'zip':
                    file = 'zip';
                    args.push('-r');
                    output += `.zip`;
                    break;
                case 'gz':
                    file = 'tar';
                    args.push('-zcvf');
                    output += `.tar.gz`;
                    break;
                case 'bz2':
                    file = 'tar';
                    args.push('-jcvf');
                    output += `.tar.bz2`;
                    break;
                case 'xz':
                    file = 'tar';
                    args.push('-Jcvf');
                    output += `.tar.xz`;
                    break;
                default:
                    throw new Error(`not supported pack algorithm : ${algorithm}`);
            }
            args.push(output);
            args.push(this.name_);
            for (let i = 0; i < names.length; i++) {
                args.push(names[i]);
            }
            console.log(file, ...args);
            return utils_1.ExecFile(file, args, {
                cwd: path_1.normalize(path_1.join(__dirname, '..', '..', 'bin')),
            });
        }
    }
}
function Build(program, os, arch, name, ext, ...packs) {
    const pack = '7z gz bz2 xz zip';
    program.command(os)
        .description(`build code to ${os}`)
        .option(`--arch [${arch.join(' ')}]`, 'GOARCH default use amd64')
        .option(`-p,--pack [${pack}]`, 'Pack to compressed package')
        .option('--debug', 'build as debug')
        .action(function () {
        const opts = this.opts();
        if (opts['Version']) {
            buildVersion();
        }
        else {
            let arch = opts['arch'];
            if (arch === undefined) {
                arch = 'amd64';
            }
            const taget = new Target(os, arch, name, ext, opts['debug']);
            taget.build().then(() => {
                return taget.pack(opts['pack'], ...packs);
            }).catch(() => {
                process.exit(1);
            });
        }
    });
}
exports.Build = Build;
async function test(pkg, args, ...names) {
    const cwd = path_1.normalize(path_1.join(__dirname, '..', '..'));
    for (let i = 0; i < names.length; i++) {
        const name = names[i];
        if (name.length > 0) {
            await utils_1.ExecFile('go', utils_1.Append(args, `${pkg}/${name}`), {
                cwd: cwd,
            });
        }
    }
}
function BuildTest(program, pkg, unit, bench) {
    program.command('test')
        .description(`run go test`)
        .option(`-v`, 'prints the full')
        .option(`-run []`, 'run func name')
        .option(`-b, --bench`, 'test bench')
        .action(function () {
        const opts = this.opts();
        const args = ['test'];
        if (opts['v']) {
            args.push('-v');
        }
        const run = opts['Run'];
        let names = bench;
        if (opts['bench']) {
            if (typeof run === "string") {
                args.push(`-run=^$`);
            }
            args.push('-benchmem', '-bench');
            if (typeof run === "string") {
                args.push(`${run}`);
            }
            else {
                args.push('.');
            }
        }
        else {
            if (typeof run === "string") {
                args.push(`-run=${run}`);
            }
            names = unit;
        }
        test(pkg, args, ...names).catch(() => {
            process.exit(1);
        });
    });
}
exports.BuildTest = BuildTest;
async function buildVersion() {
    let tag = '';
    let commit = '';
    try {
        tag = await utils_1.Exec('git', ['describe']);
    }
    catch (e) {
    }
    try {
        commit = await utils_1.Exec('git', ['rev-parse', 'HEAD']);
    }
    catch (e) {
    }
    const date = new Date().toUTCString();
    const filename = path_1.normalize(path_1.join(__dirname, '..', '..', 'version', 'version.go'));
    const str = ['package version', '',
        '// Tag git tag', 'const Tag = `' + tag.trim() + '`', '',
        '// Commit git commit', 'const Commit = `' + commit.trim() + '`', '',
        '// Date build datetime', 'const Date = `' + date + '`', '',
    ].join("\r\n");
    console.log(str.trim());
    await fs_1.promises.writeFile(filename, str);
}
function BuildVersion(program) {
    program.command('version')
        .description('update ' + path_1.join('version', 'version.go'))
        .action(function () {
        buildVersion().catch((e) => {
            console.warn(e);
            process.exit(1);
        });
    });
}
exports.BuildVersion = BuildVersion;
