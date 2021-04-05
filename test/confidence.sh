#!/bin/bash
set -euo pipefail
#set -x
scriptdir="$(readlink -f $( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd ))"
(cd "$scriptdir"/.. ; go build)
tempfile=$(mktemp)
cd $scriptdir   # So all paths arelative to here in test data file
for tf in fixtures/*
do
    for sw in '' -chars '-breaks /' -separators -counts -no-fold '-level 2' '-depth 2' '-spacer ^ -indent 1 -counts'
    do
        echo "$tf--- $sw ----------------------------------------------------------------------------------------------------" >> "$tempfile"
        head -200 "$tf" | ../frangipanni $sw >> "$tempfile"
    done
    for sw in '' '-breaks /' -separators -counts '-level 2' '-depth 3' '-no-fold'
    do
        echo "$tf--- -format json -order alpha $sw ----------------------------------------------------------------------------------------------------" >> "$tempfile"
        head -50 "$tf" | ../frangipanni $sw -format json -order alpha | jq '.' > /dev/null
        head -50 "$tf" | ../frangipanni $sw -format json -order alpha >> "$tempfile"
    done
done
# Lua
    echo "fixtures/simplechars.txt--- -lua json.lua ----------------------------------------------------------------------------------------------------" >> "$tempfile"
    <fixtures/simplechars.txt ../frangipanni -lua ../json.lua | jp '@' >> "$tempfile"

    echo "fixtures/simplechars.txt--- -lua xml.lua ----------------------------------------------------------------------------------------------------" >> "$tempfile"
    <fixtures/simplechars.txt ../frangipanni -lua ../xml.lua >> "$tempfile"
# -skip
    echo "fixtures/log-file.txt--- -skip 3 ----------------------------------------------------------------------------------------------------" >> "$tempfile"
    <fixtures/log-file.txt ../frangipanni -skip 5 >> "$tempfile"

set -x
diff "$tempfile" $scriptdir/fixtures.log # || meld "$tempfile" $scriptdir/fixtures.log
