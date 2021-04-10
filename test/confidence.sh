#!/bin/bash
# shellcheck disable=SC2129
set -euo pipefail
#set -x
scriptdir="$(readlink -f "$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )")"
(cd "$scriptdir"/.. ; go build -ldflags="-X 'main.Version=$(git describe)'")
cd "$scriptdir"   # So all paths are relative to here in test data file
V="$(../frangipanni -version)"
if [[ "$V" != "$(git describe)" ]]
then
  echo "Incorrect version number '$V' was expecting '$(git describe)'"
  exit 1
fi
tempfile=$(mktemp)

echo "--- -breaks / -counts -depth 3 -level 20 -lua args.lua ----------------------------------------------------------------------------------------------------" >> "$tempfile"
../frangipanni -breaks / -counts -depth 3 -level 20 -lua ../args.lua arg1 arg2 </dev/null | jp '@' >> "$tempfile"

for tf in fixtures/*
do
    for sw in '' -chars '-breaks /' -separators -counts -no-fold '-level 2' '-depth 2' '-spacer ^ -indent 1 -counts'
    do
        echo "$tf--- $sw ----------------------------------------------------------------------------------------------------" >> "$tempfile"
        # shellcheck disable=SC2086
        head -200 "$tf" | ../frangipanni $sw >> "$tempfile"
    done
    for sw in '' '-breaks /' -separators -counts '-level 2' '-depth 3' '-no-fold' '-down'
    do
        echo "$tf--- -format json -sort alpha $sw ----------------------------------------------------------------------------------------------------" >> "$tempfile"
        # shellcheck disable=SC2086
        head -50 "$tf" | ../frangipanni $sw -format json -sort alpha | jq '.' > /dev/null
        # shellcheck disable=SC2086
        head -50 "$tf" | ../frangipanni $sw -format json -sort alpha >> "$tempfile"
        echo "$tf--- -counts -sort counts -down $sw ----------------------------------------------------------------------------------------------------" >> "$tempfile"
        # shellcheck disable=SC2086
        head -50 "$tf" | ../frangipanni $sw -counts -sort counts -down >> "$tempfile"
    done
done
# Lua
    echo "fixtures/simplechars.txt--- -lua json.lua ----------------------------------------------------------------------------------------------------" >> "$tempfile"
    <fixtures/simplechars.txt ../frangipanni -lua ../json.lua | jp '@' >> "$tempfile"

    echo "fixtures/simplechars.txt--- -lua xml.lua ----------------------------------------------------------------------------------------------------" >> "$tempfile"
    <fixtures/simplechars.txt ../frangipanni -lua ../xml.lua >> "$tempfile"
# -skip
    echo "fixtures/log-file.txt--- -skip 5 ----------------------------------------------------------------------------------------------------" >> "$tempfile"
    <fixtures/log-file.txt ../frangipanni -sort input -skip 5 >> "$tempfile"

#set -x
if diff "$tempfile" "$scriptdir"/fixtures.log # || meld "$tempfile" $scriptdir/fixtures.log
then
  : "PASS"
else
  echo "FAIL"
  echo "maybe - meld $tempfile $scriptdir/fixtures.log"
fi