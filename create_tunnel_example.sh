#!/bin/bash

domain=$1
serverPort=$2
localPort=$3

ssh -t -R $serverPort:localhost:$localPort $domain sirtunnel $domain $serverPort
