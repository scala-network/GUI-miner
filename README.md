![](git-assets/miner-logo.png)

#  Torque GUI Miner

The Torque GUI miner is a beautiful, easy to use, interface for mining Torque.
It is aimed at getting people that have never mined Torque into the crypto
game by making it really simple to get started.

![Screenshot](https://i.imgur.com/ruK7z4Y.png "Screenshot")

We currently support two very popular miner backends:

1. [xtc-stak](https://github.com/contribute-torque/xtc-stak)
2. [xtcrig](https://github.com/contribute-torque/xtcrig) (note: [xmrig-nvidia](https://github.com/xmrig/xmrig-nvidia) and [xmrig-amd](https://github.com/xmrig/xmrig-amd) does not support our v5 proof-of-work algorithm yet)

If you'd like to fork this miner for you own coin, please see the __forking__
section later.

## Compiling

### Linux

The miner GUI is built using [Electron](https://electronjs.org) and
[Go](https://golang.org) using the
[Astilectron app framework](https://github.com/asticode/astilectron).

* Install Go

[https://golang.org/dl/](https://golang.org/dl/)

* Install required Go packages

```shell
go get -u github.com/asticode/go-astilectron
go get -u github.com/asticode/go-astilectron-bundler/...
go get -u github.com/asticode/go-astichartjs
go get -u github.com/asticode/go-astilectron-bootstrap
go get -u github.com/google/uuid
go get -u github.com/mitchellh/go-ps
```

* Clone and build the app

```shell
git clone git@github.com:contribute-torque/gui-miner.git contribute-torque/gui-miner
cd contribute-torque/gui-miner
make
```

NOTE: Ensure you clone the GUI miner into your working $GOPATH

If all goes well the binaries for Windows, macOS and Linux will be available in the `bin` folder.

## Forking

In the spirit of open source we'll be making it really simple to fork and
brand the miner for your own coin. Some structural changes need to be made to
simplify the process. Subscribe to issue [#3][i3] to follow the progress on this
guide.

[i3]: https://github.com/contribute-torque/gui-miner/issues/3
