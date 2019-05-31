#!/bin/sh
if [ ! -d "./data" ]; then
	 ./bin/geth --datadir ./data/ init ./settings/ethereum.json
fi
if [ "$1" = "--mine" ]; then
	nohup ./bin/geth --networkid 15 --datadir ./data/ --identity "ethereum" $1 --etherbase $2 > geth.log 2>&1 &
else
	nohup ./bin/geth --networkid 15 --datadir ./data/ --identity "ethereum" > geth.log 2>&1 &
fi
