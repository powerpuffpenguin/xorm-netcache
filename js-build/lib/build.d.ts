import { Command } from '../commander';
export declare function Build(program: Command, os: string, arch: Array<string>, name: string, ext: string, ...packs: Array<string>): void;
export declare function BuildTest(program: Command, pkg: string, unit: Array<string>, bench: Array<string>): void;
export declare function BuildVersion(program: Command): void;
