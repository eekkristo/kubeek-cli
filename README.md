# kubeekcli

`kubeekcli` is a small, fast CLI for **template-driven configuration**.

It lets you:

- Keep reusable templates with placeholders like `_clustername_`.
- Render placeholders to real values from `config.json` or `argument based --defaults _placeholder_=value`.
- **Generate** new folders interactively (prompt-based).
- **Generate template from real based environment**

## Why yet another templating engine?

One of my reasoning was a specific case where I had to work with a lot of terraform and helm chart files. As systems scale, it requires more and more copy-paste and seperations to avoid configuration maintenance and/or scalability issues. + those damn errors you sometimes make :)

Then it always comes into keeping template(boilerplate) up to date. If you are doing daily development/changes/testings it can become very annoying to update them.

This templater allows me to generate/scaffold new environments basically for any configuration tool. It doesn't care if it is terraform, ansible, helm etc. It just expects to have `_placeholder_` and it can read it from the file and it should simply work as long as you do not have `_<any-other-param>_` like this, 

**GULP I hope at least.**

# Table of Contents
---
- [kubeekcli](#kubeekcli)
  - [Why yet another templating engine?](#why-yet-another-templating-engine)
- [Table of Contents](#table-of-contents)
  - [Requirements](#requirements)
  - [Install](#install)
  - [Compire yourself](#compire-yourself)
- [Quick Start (TL;DR)](#quick-start-tldr)
- [Notes:](#notes)
  - [Features to be added](#features-to-be-added)
- [Releases](#releases)
- [License](#license)


## Requirements

- Go **1.21+**
- A shell (Bash, Zsh, or PowerShell)

Dependencies are managed with Go modules. Run:

```bash
go mod tidy
```

## Install

Download from [release](https://github.com/eekkristo/kubeek-cli/releases)

## Compire yourself
Clone and build:

```
git clone https://github.com/your-org/kubeekcli.git
cd kubeekcli
go build -o bin/kubeekcli ./cmd/kubeekcli
```

# Quick Start (TL;DR)

1. Put template files in ./templates/ with placeholders using underscores `_PLACEHOLDER_`, e.g `_clustername_`.
2. kubeekcli generate --name <name-of-folder> 
3. Follow prompt

# Notes:

This is still in early phase. Additionally I am by no means any senior developer. I just like to develop from time to time and mess around.

## Features to be added

The following features are being implemented in couple of new releases
- revert generated config back to template or wise/versa
- finish param based option to integrating it into CI/CD tools
- Add exclude files and directories (looking at you terraform!)

# Releases

Check out the [Changelog](CHANGELOG.md) to see what has changed, added.

# License

This cli is using https://github.com/urfave/cli
under [MIT](https://github.com/urfave/cli/blob/main/LICENSE)


