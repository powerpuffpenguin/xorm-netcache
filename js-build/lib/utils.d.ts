/// <reference types="node" />
import { ExecFileOptions } from 'child_process';
export interface Dict {
    [key: string]: any;
}
export declare function Merge(...objs: Array<Dict>): Dict;
export declare function Env(obj: Dict): Dict;
export declare function ExecFile(file: string, args: Array<string>, opts?: ExecFileOptions): Promise<void>;
export declare function Exec(file: string, args: Array<string>, opts?: ExecFileOptions): Promise<string>;
export declare function Append(items: Array<any>, ...elems: Array<any>): any[];
export declare function ClearDirectory(filename: string): Promise<void>;
export declare function RmDirectory(filename: string): Promise<void>;
