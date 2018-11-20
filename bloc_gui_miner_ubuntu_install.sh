#!/bin/bash

# This has been tested on Ubuntu 16.04 ONLY!
# Use it on your own risk!

# https://stackoverflow.com/a/5947802/1057527
GREEN='\033[0;32m'
YELLOW='\033[0;33m'
LIGHT_BLUE='\033[0;34m'
NC='\033[0m' # No Color
STEP=1

# https://gist.github.com/JamieMason/4761049
function echo_step {
	printf "${GREEN}${1}${NC}\n"
}
function echo_warning {
	printf "${YELLOW}${1}${NC}\n"
}
function echo_success {
	printf "${LIGHT_BLUE}${1}${NC}\n"
}

# Updating repositories list
echo_step "${STEP}. Updating repositories list...\n"
STEP=$((STEP + 1))

sudo apt-get update
printf "\n"

# Installing dependencies
echo_step "${STEP}. Installing dependencies...\n"
STEP=$((STEP + 1))
sudo apt-get install gcc make python libmicrohttpd10 libnss3 -y
printf "\n"

# Installing golang
echo_step "${STEP}. Installing golang...\n"
STEP=$((STEP + 1))

cd ~
wget https://dl.google.com/go/go1.11.2.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.11.2.linux-amd64.tar.gz
rm -f go1.11.2.linux-amd64.tar.gz

cp ~/.bashrc ~/.bashrc.original
echo "" >> ~/.bashrc
echo "# golang" >> ~/.bashrc
echo "export PATH=\$PATH:/usr/local/go/bin" >> ~/.bashrc
echo "export PATH=\$PATH:~/go/bin" >> ~/.bashrc

export PATH=$PATH:/usr/local/go/bin
export PATH=$PATH:~/go/bin

# Cloning GUI-miner
echo_step "${STEP}. Cloning GUI-miner...\n"
STEP=$((STEP + 1))

cd ~
rm -fr GUI-miner
git clone https://github.com/furiousteam/GUI-miner.git GUI-miner

# Installing GUI-miner dependencies
echo_step "${STEP}. Installing GUI-miner dependencies...\n"
STEP=$((STEP + 1))

cd ~
rm -fr ~/go
go get -u github.com/asticode/go-astilectron
go get -u github.com/asticode/go-astilectron-bundler/...
go get -u github.com/asticode/go-astichartjs
go get -u github.com/asticode/go-astilectron-bootstrap
go get -u github.com/google/uuid
go get -u github.com/mitchellh/go-ps
go get -u github.com/furiousteam/gui-miner/src/gui
go get -u github.com/konsorten/go-windows-terminal-sequences
go get -u github.com/mattn/go-colorable

# Updating GUI-miner electron version
echo_step "${STEP}. Updating GUI-miner electron version...\n"
STEP=$((STEP + 1))

cp ~/go/src/github.com/asticode/go-astilectron/astilectron.go ~/go/src/github.com/asticode/go-astilectron/astilectron.go.original
sed -i 's/1.8.1/3.0.8/' ~/go/src/github.com/asticode/go-astilectron/astilectron.go

cd ~/go/src/github.com/asticode/go-astilectron-bundler
make

# Compile GUI-miner
echo_step "${STEP}. Compile GUI-miner...\n"
STEP=$((STEP + 1))

cd ~/GUI-miner
make

echo_success "Done!"
