#!/bin/sh

# ztypes_linux_amd64.go
GOOS=${1:-linux}
GOARCH=${2:-amd64}
GOOSARCH="${GOOS}_${GOARCH}"

cwd=`pwd`
GOARCH=$GOARCH find -type f -name types_$GOOS.go | while read infile; do
    dir=`dirname $infile`
    inbase=`basename $infile`
    outbase=ztypes_$GOOSARCH.go
    # pushd $dir
    cd $dir
    go tool cgo -godefs $inbase | gofmt > $outbase
    # popd
    cd $cwd
done
