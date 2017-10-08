#!/bin/bash
# -*- coding: utf-8 -*-

set -eux

mainsrc=./md4pt.go
dirname=./dist

if [ -e $dirname ]; then
  : $dirname is exist
else
  mkdir $dirname
fi

GOOS=linux GOARCH=amd64 go build -o $dirname/md4p_linux $mainsrc
GOOS=darwin GOARCH=amd64 go build -o $dirname/md4p_mac $mainsrc
GOOS=windows GOARCH=amd64 go build -o $dirname/md4p_win.exe $mainsrc

ls $dirname
: build completed

