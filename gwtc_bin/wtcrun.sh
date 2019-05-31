#!/bin/sh
if [ ! -d "./data/geth" ]; then
	./bin/geth --datadir ./data/ init ./settings/wtc.json
fi
if [ "$1" = "--mine" ]; then
./bin/geth --networkid 15 --datadir ./data/ --identity "wtc" $1 --etherbase $2 console
else
./bin/geth --networkid 15 --datadir ./data/ --identity "wtc" console
fi
