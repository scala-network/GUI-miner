This file is a small tutorial for developers. But you still need to read the README.md file first.  

### Intro

1. The **miner** source folder is `~/GUI-miner`  
2. The **libraries** source folder is `~/go/src/github.com/furiousteam/BLOC-GUI-Miner/`  
3. The **xmr-stak** binaries folder is `~/xmr-stak-linux-2.10.7-cpu/`  

**Notes:**  
* The **.go** files you need to change are from the **miner** folder  
* The **html**, **css**, **javascript** and **images** files you need to change are from the **libraries** folder  
* The `export DISPLAY=:0` is needed when the command is run over SSH  
* The `~/GUI-miner/bin/linux-amd64/config.json` is the GUI miner config file  
* The `~/GUI-miner/bin/linux-amd64/miner/config.txt`, `~/GUI-miner/bin/linux-amd64/miner/cpu.txt` and `~/GUI-miner/bin/linux-amd64/miner/pools.txt` are the xmr-stak config files  

### To compile and run the miner clean

```shell
clear; export DISPLAY=:0; cd ~/GUI-miner; \
make clean; make; mkdir ~/GUI-miner/bin/linux-amd64/miner; \
cp -r ~/xmr-stak-linux-2.10.7-cpu/xmr-stak ~/GUI-miner/bin/linux-amd64/miner; \
cd ~/GUI-miner/bin/linux-amd64/; ./BLOC\ GUI\ Miner\ v0.0.4 -d
``` 

### To compile and run the miner with old configuration

```shell
clear; export DISPLAY=:0; cd ~/GUI-miner; \
cp -u ~/GUI-miner/bin/linux-amd64/config.json /tmp/config.json; \
cp -u ~/GUI-miner/bin/linux-amd64/miner/config.txt /tmp/config.txt; \
cp -u ~/GUI-miner/bin/linux-amd64/miner/cpu.txt /tmp/cpu.txt; \
cp -u ~/GUI-miner/bin/linux-amd64/miner/pools.txt /tmp/pools.txt; make clean; make; \
mkdir ~/GUI-miner/bin/linux-amd64/miner; \
cp -r ~/xmr-stak-linux-2.10.7-cpu/xmr-stak ~/GUI-miner/bin/linux-amd64/miner; \
cp -u /tmp/config.json ~/GUI-miner/bin/linux-amd64/config.json; \
cp -u /tmp/config.txt ~/GUI-miner/bin/linux-amd64/miner/config.txt; \
cp -u /tmp/cpu.txt ~/GUI-miner/bin/linux-amd64/miner/cpu.txt; \
cp -u /tmp/pools.txt ~/GUI-miner/bin/linux-amd64/miner/pools.txt; \
cd ~/GUI-miner/bin/linux-amd64/; ./BLOC\ GUI\ Miner\ v0.0.4 -d
```

### To run the miner with clean configuration

```shell
clear; export DISPLAY=:0; cd ~/GUI-miner/bin/linux-amd64/; \
rm -f ~/GUI-miner/bin/linux-amd64/config.json; \
rm -f ~/GUI-miner/bin/linux-amd64/miner/config.txt; \
rm -f ~/GUI-miner/bin/linux-amd64/miner/cpu.txt; \
rm -f ~/GUI-miner/bin/linux-amd64/miner/pools.txt; \
./BLOC\ GUI\ Miner\ v0.0.4 -d
```

### To run the miner with old configuration

```shell
clear; export DISPLAY=:0; cd ~/GUI-miner/bin/linux-amd64/; ./BLOC\ GUI\ Miner\ v0.0.4 -d
```
