#!/bin/bash
function error() { echo $@; exit -1; }

cd $(dirname ${BASH_SOURCE[0]})
type git >/dev/null || error "git not found"
tempdir=$(mktemp -d)
REPO=https://githubfast.com/bellard/quickjs.git
git clone --depth=1 $REPO $tempdir || {
    rm -rf $tempdir
    error "Clone $REPO failed"
}
for name in $(ls -1 .); do
    test -f $tempdir/$name && cp $tempdir/$name .
done
version=$(<$tempdir/VERSION)
sed -i "s/CONFIG_VERSION/$verison/g" quickjs.c
rm -rf $tempdir
cd -
