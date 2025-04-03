# Changelog

All notable changes to this project will be documented in this file.

The format is based on [Keep a Changelog](https://keepachangelog.com/en/1.1.0/),
and this project adheres to [Semantic Versioning](https://semver.org/spec/v2.0.0.html).

## [Unreleased]

### TODO

+ Modify examples to use only one example with different routes, 
so then you only need to import the project one time and you can show how to use the middleware per route group, etc for the wrapped functionality
+ If you put a "wrong" wrapper field then all the json is shown, when maybe we should return the default behaviour
+ If you chose indented json, you lose the indentation

## [v0.0.2] 2025-11-01

- Updating Content-Length header when json is modified
- Refactor examples to use one single go.mod
- Updating dependencies

## [v0.0.1]

### Added

- Adding partial response gin framework.
  - Support for arrays. 
  - Support for nested objects.
  - Support json wrapper response