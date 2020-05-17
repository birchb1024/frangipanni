#!/bin/bash
set -xeuo pipefail
scriptdir="$(readlink -f $( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd ))"

git describe 
go build

test/confidence.sh

GOOS=windows GOARCH=386 go build -o frangipanni.exe frangipanni.go

rm -f frangipanni.zip rm -f frangipanni.tgz
tar zcvf frangipanni.tgz frangipanni
zip frangipanni.zip frangipanni.exe

ls -ltr
