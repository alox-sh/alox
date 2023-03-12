#!/bin/bash

echo "0: $0"
echo "1: $1"

find . -name *.go -exec sed -i -e "s/alox.sh/$1/g" {} \;
