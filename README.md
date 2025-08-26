# GitMan (Git Manager)

Your friendly neighbor who’s too lazy to use Git

## What is gitman

gitman is your friendly Git sidekick 🛠️.
A command-line tool that saves you from typing long commands and memorizing commit hashes.
Powered by fzf, it helps you navigate branches, commits, and actions with ease — so you can spend more time coding, and less time wrestling with Git.

![gitman-log](./demo/gitman-log-demo.png)

## Dependence

| name | version |
| --- | --- |
| golang | >= 1.24.2 |
| fzf | >= 0.65.1 |
| git | >= 1.51.0 |

## Installation

### With Go

```bash
go install github.com/o-kaisan/gitman@latest
```

### Other

```
make install
# or
bash install.sh
```

## Features

### Commit Action

```
gitman log
# or
gitman l
```

- select commit
![gitman-log](./demo/gitman-log-demo.png)

- select git command
![gitman-log-action](./demo/gitman-log-select-action-demo.png)


### Branch Action

```
gitman branch
# or
gitman br
```

- select branch
![gitman-log](./demo/gitman-branch-demo.png)

- select branch action
![gitman-log-action](./demo/gitman-log-select-action-demo.png)


## Environment variable

| value | type | default | description |
| -- | -- | -- | -- |
| GITMAN_DEBUG | bool | false |  debug mode (default: "false") |
| GITMAN_BRANCH_ALIAS | string | br | change branch command alias (default: "br") |
| GITMAN_LOG_ALIAS | string | l | change log command alias (default: "l") |
| GITMAN_LOG_DISPLAY_LIMIT | string | 100 |change log display limit (default: "100")|
