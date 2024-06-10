#!/bin/sh

[ -z "`which jq`" ] && echo "jq not found, try 'brew install jq' or get it from https://stedolan.github.io/jq/" && exit 1
[ -z "`which unzip`" ] && echo "unzip not found" && exit 1
[ -z "`which curl`" ] && echo "curl not found" && exit 1

echo "checking for latest SQLite3MultipleCiphers release"
release_json="`curl -s https://api.github.com/repos/utelle/SQLite3MultipleCiphers/releases/latest | jq '{name: .name} + (.assets[] | select(.name | test("-amalgamation.zip$")) | {filename: .name, url: .browser_download_url})'`"

[ -z "$release_json" ] && echo "could not get release assets from github repository, aborting" && exit 1

url="`echo "$release_json" | jq -r '.url'`"
name="`echo "$release_json" | jq -r '.name'`"
filename="`echo "$release_json" | jq -r '.filename'`"

echo "latest release: $name"

# get script directory
target="`cd $(dirname "$0"); pwd`"

# make temporary path
tmp="`mktemp -d`"

echo "downloading latest SQLite3 source code"
curl -s -L "$url" -o "$tmp/sqlite3mc-update.zip"

cd "$tmp"
echo "extracting files"
unzip -q "sqlite3mc-update.zip"


updated=0

[ ! -f "$target/sqlite3.c" ] && updated=1
[ ! -f "$target/sqlite3.h" ] && updated=1
[ ! -f "$target/sqlite3ext.h" ] && updated=1

if [ $updated -eq 0 ]; then
    diff -q "$tmp/sqlite3mc_amalgamation.c" "$target/sqlite3.c" > /dev/null
    [ $? -eq 1 ] && updated=1
fi
if [ $updated -eq 0 ]; then
    diff -q "$tmp/sqlite3mc_amalgamation.h" "$target/sqlite3.h" > /dev/null
    [ $? -eq 1 ] && updated=1
fi
if [ $updated -eq 0 ]; then
    diff -q "$tmp/sqlite3ext.h" "$target/sqlite3ext.h" > /dev/null
    [ $? -eq 1 ] && updated=1
fi


if [ $updated -eq 1 ]; then

    cd "$target"

    git diff --exit-code > /dev/null && git diff --cached --exit-code > /dev/null && dirty=0 || dirty=1

    [ $dirty -eq 1 ] && echo "stashing changes" && git stash -q

    echo "updating files"
    cp "$tmp/sqlite3mc_amalgamation.c" "$target/sqlite3.c"
    cp "$tmp/sqlite3mc_amalgamation.h" "$target/sqlite3.h"
    cp "$tmp/sqlite3ext.h" "$target/sqlite3ext.h"

    git add "$target/sqlite3.c" "$target/sqlite3.h" "$target/sqlite3ext.h"

    echo "committing SQLite3 update"
    git commit -m "Update SQLite amalgamation code to $name: https://github.com/utelle/SQLite3MultipleCiphers"

    [ $dirty -eq 1 ] && echo "restoring changes" && git stash pop -q
else
    echo "no changes to SQLite3 detected"
fi


# ls -la "$tmp"

# cleanup
echo "cleanup"
rm -rf "$tmp"/*
rmdir "$tmp"

echo "done"
