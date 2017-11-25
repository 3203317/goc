#!/usr/bin/env bash

if [ ! -f install ]; then
echo 'install.sh must be run from its folder' 1>&2
exit 1
fi

CURDIR=`pwd`
OLDGOPATH="$GOPATH"
export GOPATH="$CURDIR"

gofmt -w src

go install goc

export GOPATH="$OLDGOPATH"

echo 'finished'
