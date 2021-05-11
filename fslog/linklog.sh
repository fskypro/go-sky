#! /bin/bash

src=$1
dst=./daylog.log
arg=$2

ln -sf $src $dst

echo "$src -> $dst"
echo "args: $arg"
