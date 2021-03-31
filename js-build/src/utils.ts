import { execFile, ExecFileOptions } from 'child_process'
import { Dir, promises } from 'fs'
import { join } from 'path'
export interface Dict {
    [key: string]: any;
}

export function Merge(...objs: Array<Dict>): Dict {
    const result = {} as Dict
    for (let i = 0; i < objs.length; i++) {
        const obj = objs[i]
        if (obj) {
            for (const k in obj) {
                result[k] = obj[k]
            }
        }
    }
    return result
}
export function Env(obj: Dict) {
    return Merge(process.env, obj)
}
export function ExecFile(file: string, args: Array<string>, opts?: ExecFileOptions): Promise<void> {
    return new Promise<void>((resolve, reject) => {
        execFile(file, args, opts, (e, stdout, stderr) => {
            process.stdout.write(stdout)
            process.stderr.write(stderr)
            if (e) {
                reject(e)
            } else {
                resolve()
            }
        })
    })
}
export function Exec(file: string, args: Array<string>, opts?: ExecFileOptions): Promise<string> {
    return new Promise<string>((resolve, reject) => {
        execFile(file, args, opts, (e, stdout, stderr) => {
            if (e) {
                if (typeof stderr === "string" && stderr.length > 0) {
                    reject(stderr)
                } else {
                    reject(e)
                }
            } else {
                resolve(stdout ?? '')
            }
        })
    })
}
export function Append(items: Array<any>, ...elems: Array<any>) {
    const obj = []
    obj.push(...items)
    obj.push(...elems)
    return obj
}
export async function ClearDirectory(filename: string): Promise<void> {
    let dirs: Dir
    try {
        dirs = await promises.opendir(filename)
    } catch (e) {
        if (e.code === `ENOENT`) {
            await promises.mkdir(filename,
                {
                    recursive: true,
                    mode: 0o775,
                },
            )
            return
        }
    }

    for await (const dirent of dirs) {
        if (dirent.isDirectory()) {
            await RmDirectory(join(filename, dirent.name))
        } else {
            await promises.rm(join(filename, dirent.name))
        }
    }
}
export async function RmDirectory(filename: string): Promise<void> {
    let dirs: Dir
    try {
        dirs = await promises.opendir(filename)
    } catch (e) {
        if (e.code === `ENOENT`) {
            await promises.mkdir(filename,
                {
                    recursive: true,
                    mode: 0o775,
                },
            )
            return
        }
    }
    for await (const dirent of dirs) {
        if (dirent.isDirectory()) {
            await RmDirectory(join(filename, dirent.name))
        } else {
            await promises.rm(join(filename, dirent.name))
        }
    }
    await promises.rmdir(filename)
}