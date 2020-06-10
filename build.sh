#!/bin/bash
set -xeuo pipefail
scriptdir="$(readlink -f $( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd ))"

git describe 
go build

test/confidence.sh

GOOS=windows GOARCH=386 go build -o frangipanni.exe frangipanni.go
GOOS=darwin GOARCH=amd64 go build -o frangipanni_mac frangipanni.go

rm -f frangipanni.zip rm -f frangipanni.tgz
tar zcvf frangipanni.tgz frangipanni *.lua
zip frangipanni.zip frangipanni.exe frangipanni_mac *.lua

ls -ltr
