#! /bin/bash

ROOT=`cd $(dirname $0) && pwd`
cd $ROOT
ln -sf $1 ./$2.log
echo  link new log file \"$1\" to \"./$2.log\"
