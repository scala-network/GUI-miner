# BLOC GUI Miner

BLOC GUI miner is a beautiful, easy to use, Graphical User interface for mining the cryptocurrency [BLOC](https://bloc.money).
It is aimed at getting people that have never tried mining before with a focus on accessibility, security and simplicity.
BLOC GUI miner makes getting started with [BLOC](https://bloc.money) mining easier than ever.

## Getting started

[What is mining ?](#) If you want to learn about cryptocurrencies, mining is a great place to start!. Start mining today and get rewarded in BLOC.

What is [XMR-STAK](#) and how to mine with it
What is [XMRIG](#) and how to mine with it

## **Make sure you have a wallet**

If you have not yet downloaded and ran the [BLOC](https://bloc.money/download) software to sync the blockchain and create a wallet, you need to [create your wallet](#) before start mining.

## **Choose your mining pool**

You can find a complete list of the BLOC mining pools available on the [BLOC MINING](https://bloc.money/mining) section of our website. We suggest you to select the nearest mining pool following your location for the best mining experience and results.


## **Linux**

BLOC GUI Miner comes with XMR-STAK already built in including configuration files for **CPU mining ONLY**.

- Download and install [BLOC GUI Miner](https://github.com/furiousteam/GUI-miner/releases/latest) from GitHub
- From the [Download Area](https://bloc.money/download) of BLOC.MONEY
- Unzip the file to your desktop

It must looks like this:

![Screenshot](https://i.imgur.com/ruK7z4Y.png "Screenshot")

Inside the miner folder:

![Screenshot](https://i.imgur.com/ruK7z4Y.png "Screenshot")


You are free to use your own XMR-STAK binaries as long as it is the same version compatible with the BLOC GUI Miner.




## **Downloading and Installing**

BLOC GUI Miner comes with XMR-STAK already built in including configuration files for CPU and GPU mining in the most cases.

1. Download and install [BLOC GUI Miner](https://github.com/furiousteam/GUI-miner/releases/latest)
2. Unzip the file to your computer in your desktop for exemple

## **Start the BLOC GUI Miner**

- Windows: Double click the icon BLOC GUI Miner.exe
- Linux: Double click the icon app image Bloc GUI Miner.App image
- MacOs: Double click the icon BLOC GUI Miner






It will auto-detect your hardware, and tune everything for you.
2. Create a folder called `miner` (if it's not already created) next to the `BLOC GUI Miner.exe` and unzip the files you just downloaded for XMR-Stak in there.
3. Double-click on `xmr-stak.exe`.

1. Download and install [XMR-Stak Unified Miner](https://github.com/fireice-uk/xmr-stak/releases/latest). It will auto-detect your hardware, and tune everything for you.
2. Create a folder called `miner` (if it's not already created) next to the `BLOC GUI Miner.exe` and unzip the files you just downloaded for XMR-Stak in there.
3. Double-click on `xmr-stak.exe`.














![Screenshot](https://i.imgur.com/ruK7z4Y.png "Screenshot")

## BLOC GUI Features

BLOC GUI miner is a wrapper for the most popular cryptonote based coins miner XMR-STAK and XMRIG.

## Supported Miners

We currently support two very popular miner backends:

1. [xmr-stak](https://github.com/fireice-uk/xmr-stak)
2. [xmrig](https://github.com/xmrig/xmrig) (note: [xmrig-nvidia](https://github.com/xmrig/xmrig-nvidia) and [xmrig-amd](https://github.com/xmrig/xmrig-amd) does not support our v4 proof-of-work algorithm yet)

## Compiling on Linux (Ubuntu)

Compiling on Linux will generate the binaries for Windows, macOS and Linux.

The miner GUI is built using [Electron](https://electronjs.org) and
[Go](https://golang.org) using the
[Astilectron app framework](https://github.com/asticode/astilectron).

### Install dependencies

```shell
sudo apt-get update
sudo apt-get install gcc make python libmicrohttpd10 libnss3 -y
```

- gcc and make are required for go packages  
- python is required for GUI-miner  
- libmicrohttpd is required for xmrig  
- libnss3 is required for electron  

### Install Go

1. [https://golang.org/dl/](https://golang.org/dl/)

2. or follow the next lines

download and unpack golang binaries

```shell
cd ~
wget https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.11.2.linux-amd64.tar.gz
```

add Go to current $PATH, by editing the current user's `.bashrc`

```shell
nano ~/.bashrc 
```

add the following

```shell
# golang
export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:~/go/bin
```

run the .bashrc file (otherwise you need to logout and login again for it to get triggered)

```shell
source ~/.bashrc
```

now you could check the go installation by running

```shell
go version
```

### Clone the app

```shell
cd ~
git clone https://github.com/furiousteam/GUI-miner.git
```

### Install required Go packages

```shell
go get -u github.com/asticode/go-astilectron
go get -u github.com/asticode/go-astilectron-bundler/...
go get -u github.com/asticode/go-astichartjs
go get -u github.com/asticode/go-astilectron-bootstrap
go get -u github.com/google/uuid
go get -u github.com/mitchellh/go-ps
go get -u github.com/furiousteam/gui-miner/src/gui
go get -u github.com/konsorten/go-windows-terminal-sequences
go get -u github.com/mattn/go-colorable
```

### Update electron version

edit `~/go/src/github.com/asticode/go-astilectron/astilectron.go` file

```shell
nano ~/go/src/github.com/asticode/go-astilectron/astilectron.go
```

and change `VersionElectron         = "1.8.1"` to `VersionElectron         = "3.0.8"`

then, recompile go-astilectron-bundler

```shell
cd ~/go/src/github.com/asticode/go-astilectron-bundler
make
```

### Compile the miner

```shell
cd ~/GUI-miner
make
```

If all goes well, the binaries for Windows, macOS and Linux will be available in the `bin` folder.

### Attach the miner

before you start the GUI-miner, make sure you have copied the binaries of [xmrig](https://github.com/xmrig/xmrig) or [xmr-stak](https://github.com/fireice-uk/xmr-stak) into the `miner` subfolder right next to the main GUI-miner executable