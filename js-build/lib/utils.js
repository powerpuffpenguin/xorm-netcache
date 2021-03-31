"use strict";
Object.defineProperty(exports, "__esModule", { value: true });
exports.RmDirectory = exports.ClearDirectory = exports.Append = exports.Exec = exports.ExecFile = exports.Env = exports.Merge = void 0;
const child_process_1 = require("child_process");
const fs_1 = require("fs");
const path_1 = require("path");
function Merge(...objs) {
    const result = {};
    for (let i = 0; i < objs.length; i++) {
        const obj = objs[i];
        if (obj) {
            for (const k in obj) {
                result[k] = obj[k];
            }
        }
    }
    return result;
}
exports.Merge = Merge;
function Env(obj) {
    return Merge(process.env, obj);
}
exports.Env = Env;
function ExecFile(file, args, opts) {
    return new Promise((resolve, reject) => {
        child_process_1.execFile(file, args, opts, (e, stdout, stderr) => {
            process.stdout.write(stdout);
            process.stderr.write(stderr);
            if (e) {
                reject(e);
            }
            else {
                resolve();
            }
        });
    });
}
exports.ExecFile = ExecFile;
function Exec(file, args, opts) {
    return new Promise((resolve, reject) => {
        child_process_1.execFile(file, args, opts, (e, stdout, stderr) => {
            if (e) {
                if (typeof stderr === "string" && stderr.length > 0) {
                    reject(stderr);
                }
                else {
                    reject(e);
                }
            }
            else {
                resolve(stdout ?? '');
            }
        });
    });
}
exports.Exec = Exec;
function Append(items, ...elems) {
    const obj = [];
    obj.push(...items);
    obj.push(...elems);
    return obj;
}
exports.Append = Append;
async function ClearDirectory(filename) {
    let dirs;
    try {
        dirs = await fs_1.promises.opendir(filename);
    }
    catch (e) {
        if (e.code === `ENOENT`) {
            await fs_1.promises.mkdir(filename, {
                recursive: true,
                mode: 0o775,
            });
            return;
        }
    }
    for await (const dirent of dirs) {
        if (dirent.isDirectory()) {
            await RmDirectory(path_1.join(filename, dirent.name));
        }
        else {
            await fs_1.promises.rm(path_1.join(filename, dirent.name));
        }
    }
}
exports.ClearDirectory = ClearDirectory;
async function RmDirectory(filename) {
    let dirs;
    try {
        dirs = await fs_1.promises.opendir(filename);
    }
    catch (e) {
        if (e.code === `ENOENT`) {
            await fs_1.promises.mkdir(filename, {
                recursive: true,
                mode: 0o775,
            });
            return;
        }
    }
    for await (const dirent of dirs) {
        if (dirent.isDirectory()) {
            await RmDirectory(path_1.join(filename, dirent.name));
        }
        else {
            await fs_1.promises.rm(path_1.join(filename, dirent.name));
        }
    }
    await fs_1.promises.rmdir(filename);
}
exports.RmDirectory = RmDirectory;
