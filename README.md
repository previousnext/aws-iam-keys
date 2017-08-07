AWS IAM: Authorized Keys
========================

[![CircleCI](https://circleci.com/gh/previousnext/aws-iam-keys.svg?style=svg)](https://circleci.com/gh/previousnext/aws-iam-keys)

**Maintainer**: Nick Schuch

A simple `authorized_keys` file generator backed by an IAM group.

## How it works

* Writes a new `authorized_keys` file
* Ensures permissions are correct

## Development

### Principles

* Code lives in the `workspace` directory

### Tools

* **Dependency management** - https://getgb.io
* **Build** - https://github.com/mitchellh/gox
* **Linting** - https://github.com/golang/lint

### Workflow

(While in the `workspace` directory)

**Installing a new dependency**

```bash
gb vendor fetch github.com/foo/bar
```

**Running quality checks**

```bash
make lint test
```

**Building binaries**

```bash
make build
```
