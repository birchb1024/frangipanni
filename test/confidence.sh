#!/bin/bash
set -euo pipefail
scriptdir="$( cd "$( dirname "${BASH_SOURCE[0]}" )" >/dev/null 2>&1 && pwd )"

tempfile=$(mktemp)
for tf in $scriptdir/fixtures/*
do
    for sw in '' -chars '-breaks /' -separators -counts '-level 2' '-depth 2' '-format json'
    do
        echo "$tf---$sw----------------------------------------------------------------------------------------------------" >> "$tempfile"
        $scriptdir/../frangipanni $sw <"$tf" >> "$tempfile"
    done
done
set -x
diff "$tempfile" $scriptdir/fixtures.log || kdiff3 "$tempfile" $scriptdir/fixtures.log
