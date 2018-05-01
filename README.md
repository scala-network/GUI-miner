#  Stellite GUI Miner

The Stellite GUI miner is a beautiful, easy to use, interface for mining Stellite.
It is aimed at getting people that have never mined Stellite into the crypto
game by making it really simple to get started.

We currently support two very popular miner backends:

1. [xmr-stak](https://github.com/fireice-uk/xmr-stak)
2. [xmrig](https://github.com/xmrig/xmrig) (including [xmrig-nvidia](https://github.com/xmrig/xmrig-nvidia) and [xmrig-amd](https://github.com/xmrig/xmrig-amd))

If you'd like to fork this miner for you own coin, please see the __forking__
section later.

## Compiling

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

In the spirit of open source we've made it really simple to brand the miner
for your own coin. Together with this guide you'll have your own miner in 30 minutes.

For this guide we'll assume your coin is named `awesomecoin`

### How to

1. Fork the repository to your awesomecoin's account
2. Clone the repository
3. Style the miner

The miner's stylesheet is generated from a LESS file. You can find it in
`repo-root/less`. Duplicate the `stellite.less` file and rename to
`awesomecoin.less`

Open `awesomecoin.less`. The colours and font is defined at the top, simply update
it to your colour scheme

```less
/* Colours */
@color-blue-dark: #0B0C22;
@color-blue-darker: #06071D;
@color-blue-dark-lighter: #0D112D;
@color-purple-grey: #8284a5;
@color-link: #00FFFF;
@color-stellite-purple: #9012F0;
@color-stellite-purple-dark: #540C87;
@color-selected: #E8D400;
@color-text: #fafafa;

/* Fonts */
@font-family-main: 'Poppins', sans-serif;
```

Compile the .less file into `repo-root/src/resources/app/static/css/main.css`

If you do `make run` now you'll see your new colour scheme in effect.

4. Logo

TODO: Change logo

5. Icons

TODO: Icons

6. Pools API

TODO: Pools API

7. Stats API

TODO: Stats API

8. API Endpoint

TODO: change in main.go
