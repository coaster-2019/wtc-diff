#!/bin/bash
ps -ef | grep geth | awk '{ print $2 }' | sudo xargs kill -9
