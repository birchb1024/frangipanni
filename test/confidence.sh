#!/bin/bash
set -euo pipefail
scriptdir="$(readlink -f $( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd ))"
(cd "$scriptdir"/.. ; go build)
tempfile=$(mktemp)
cd $scriptdir   # So all paths arelative to here in test data file
for tf in fixtures/*
do
    for sw in '' -chars '-breaks /' -separators -counts -no-fold '-level 2' '-depth 2' 
    do
        echo "$tf--- $sw ----------------------------------------------------------------------------------------------------" >> "$tempfile"
        head -200 "$tf" | ../frangipanni $sw >> "$tempfile"
    done
    for sw in '' '-breaks /' -separators -counts '-level 2' '-depth 2'
    do
        echo "$tf--- -format json -order alpha $sw ----------------------------------------------------------------------------------------------------" >> "$tempfile"
        head -50 "$tf" | ../frangipanni $sw -format json -order alpha >> "$tempfile"
    done
done
set -x
#cp "$tempfile" $scriptdir/fixtures.log
diff "$tempfile" $scriptdir/fixtures.log || kdiff3 "$tempfile" $scriptdir/fixtures.log
