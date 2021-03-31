#!/usr/bin/env node
"use strict";
const { Command } = require('./js-build/commander')
const { Append } = require('./js-build/lib/utils')
const { Build, BuildTest, BuildVersion } = require('./js-build/lib/build')
const { BuildGRPC } = require('./js-build/lib/grpc')
const Name = 'xormcache'
const PackageName = 'github.com/powerpuffpenguin/xormcache'
const TestBench = [
    'utils'
]
const TestUnit = Append(TestBench)
const program = new Command('./build.js')
const Arch = ['amd64', '386']
Build(program, 'linux', Arch,
    Name, '',
    `${Name}.jsonnet`,
)
Build(program, 'freebsd', Arch,
    Name, '',
    `${Name}.jsonnet`,
)
Build(program, 'darwin', Arch,
    Name, '',
    `${Name}.jsonnet`,
)
Build(program, 'windows', Arch,
    Name, '.exe',
    `${Name}.jsonnet`,
)

BuildVersion(program)
BuildTest(program, PackageName, TestUnit, TestBench)
BuildGRPC(program, PackageName, '307f9480-91f3-11eb-962b-b5b119b956b2', false, false)
program.parse(process.argv)