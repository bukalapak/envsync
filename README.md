[![Build Status](https://travis-ci.org/bukalapak/envsync.svg?branch=master)](https://travis-ci.org/bukalapak/envsync)
[![codecov](https://codecov.io/gh/bukalapak/envsync/branch/master/graph/badge.svg)](https://codecov.io/gh/bukalapak/envsync)
[![Go Report Card](https://goreportcard.com/badge/github.com/bukalapak/envsync)](https://goreportcard.com/report/github.com/bukalapak/envsync)
[![Documentation](https://godoc.org/github.com/bukalapak/envsync?status.svg)](http://godoc.org/github.com/bukalapak/envsync)

# Envsync

## Description

Envsync is a tool to synchronize sample env and actual env file.

## Installation

1. Download the executable file in the given link below. Open the given link via your favorite browser. Choose **envsync_linux_amd64** for linux or **envsync_darwin_amd64** for OSX. 

    ```sh
    https://github.com/bukalapak/envsync/releases/latest
    ```

2. For the next steps, please, change `OS` with "linux" or "darwin" (depend on what you have downloaded before).

3. Give the executable file the permission to execute

    ```sh
    chmod +x ~/Downloads/envsync_[OS]_amd64
    ```

4. Move to /usr/local/bin

    ```sh
    mv ~/Downloads/envsync_[OS]_amd64 /usr/local/bin/envsync
    ```

## Usage

```
envsync -s <source file> -t <target file>
```

Source file is the sample env. If the -s flag isn't provided, envsync will use the default value which is **env.sample**.
Target file is the actual env. If the -t flag isn't provided, envsync will use the default value which is **.env**.