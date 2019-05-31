#!/bin/sh
if [ ! -d "./data" ]; then
	 ./bin/geth --datadir ./data/ init ./settings/wtc.json
fi
if [ "$1" = "--mine" ]; then
	nohup ./bin/geth --networkid 15 --datadir ./data/ --identity "wtc" $1 --etherbase $2 > geth.log 2>&1 &
else
	nohup ./bin/geth --networkid 15 --datadir ./data/ --identity "wtc" > geth.log 2>&1 &
fi
