#  Stellite GUI Miner

The (Unofficial) Stellite GUI miner is a simple, cross-platform, wrapper around the fantastic
[xmr-stak](https://github.com/fireice-uk/xmr-stak) miner.

It is aimed at getting people that have never mined Stellite into the crypto
game by making it really simple to get started.

## Compiling

The miner GUI is built using [Electron](https://electronjs.org), [Go](https://golang.org) using the [Astilectron app framework](https://github.com/asticode/astilectron).

* Install Go

[https://golang.org/dl/](https://golang.org/dl/)

* Install required Go packages

```shell
go get -u github.com/asticode/go-astilectron
go get -u github.com/asticode/go-astilectron-bundler/...
go get -u github.com/asticode/go-astichartjs
go get -u github.com/asticode/go-astilectron-bootstrap
```

* Clone and build the app

```shell
git clone git@github.com:donovansolms/stellite-gui-miner.git
cd stellite-gui-miner
make
```
If all goes well the binary will be available in the `bin` folder.
