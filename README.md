# Mooze

A command-line REST api test tool

This is very beta version of the final goal of the program

Many bugs and limited features

## Preview

![0-1-1-image](./asset/image/0-1-1(2).gif)

## Install

mooze requires go to be installed  
do not supports windows

```
$ go get github.com/rudy3091/mooze
$ MOOZE_ROOT="$GOPATH/src/github.com/rudy3091/mooze"
$ cd $MOOZE_ROOT
$ go build
$ ln -s $MOOZE_ROOT/mooze /usr/local/bin/mooze # for macOS
$ mooze
```

## Keybindings
- j, k, l, h: move around inside pane
- tab: move focus on pane
- s: send request
