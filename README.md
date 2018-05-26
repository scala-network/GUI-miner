![](git-assets/miner-logo.png)

#  Stellite GUI Miner

The Stellite GUI miner is a beautiful, easy to use, interface for mining Stellite.
It is aimed at getting people that have never mined Stellite into the crypto
game by making it really simple to get started.

![Screenshot](https://i.imgur.com/ruK7z4Y.png "Screenshot")

We currently support two very popular miner backends:

1. [xmr-stak](https://github.com/fireice-uk/xmr-stak)
2. [xmrig](https://github.com/xmrig/xmrig) (including [xmrig-nvidia](https://github.com/xmrig/xmrig-nvidia) and [xmrig-amd](https://github.com/xmrig/xmrig-amd))

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
git clone git@github.com:donovansolms/stellite-gui-miner.git
cd stellite-gui-miner
make
```
If all goes well the binaries for Windows, macOS and Linux will be available in the `bin` folder.

## Forking

In the spirit of open source we'll be making it really simple to fork and
brand the miner for your own coin. Some structural changes need to be made to
simplify the process. Subscribe to issue [#3][i3] to follow the progress on this
guide.

[i3]: https://github.com/stellitecoin/gui-miner/issues/3
