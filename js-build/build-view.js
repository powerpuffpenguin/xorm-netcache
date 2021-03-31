#!/usr/bin/env node
"use strict";
const { Command } = require('./commander')
const { BuildSource } = require('./lib/source')
const { BuildView } = require('./lib/view')


const program = new Command('./build-view.js')

BuildSource(program, true)
BuildView(program)

program.parse(process.argv)