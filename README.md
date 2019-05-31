![](images/wtc_logo.jpg)

# GO WTC
Waltonchain Mainnet User Manual


## Steps

### 1. Run docker container
Install latest distribution of [Go](https://golang.org "Go") if you don't have it already.  
`# sudo apt-get install -y build-essential`  

### 2. Compile source code
`# cd /usr/local/src`  
`# git clone https://github.com/WaltonChain/WaltonChain_Geth_Src.git`  
`# cd WaltonChain_Geth_Src`  
`# make geth`  
`# ./build/bin/geth version`  

### 3. Deploy
`# cd /usr/local/src/WaltonChain_Geth_Src/geth_bin/`  
`# cp ../build/bin/geth ./bin/geth`  
`# ./backend.sh`

### 4. Enter console
`# cd /usr/local/src/WaltonChain_Geth_Src/geth_bin/`  
`# ./bin/geth attach ./data/geth.ipc`

### 5. View information of the connected node
`# admin.peers`

### 6. Create account
`# personal.newAccount()`  
`# ******`  ---- Enter new account password  
`# ******`  ---- Confirm the new account password  

### 7. Mine
`# miner.start()`

### 8. Query
`# wtc.getBalance(wtc.coinbase)`

### 9. Unlock account
`# personal.unlockAccount(wtc.coinbase)`

### 10. Transfer
`# wtc.sendTransaction({from: wtc.accounts[0], to: wtc.accounts[1], value: web3.toWei(1)})`

### 11. Exit console
`# exit`

### 12. View log
`# cd /usr/local/src/WaltonChain_Geth_Src/geth_bin/`  
`# tail -f geth.log`

### 13. Stop geth
`# cd /usr/local/src/WaltonChain_Geth_Src/geth_bin/`  
`# ./stop.sh` 


## Acknowledgement
We hereby thank:  
· [Ethereum](https://www.ethereum.org/ "Ethereum")




