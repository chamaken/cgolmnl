#!/bin/sh

GOOS=${1:-linux}
GOARCH=${2:-amd64}
GOOSARCH="${GOOS}_${GOARCH}"

cwd=`pwd`
cmd=${1:-"build"}

case "$cmd" in
    build)
	for dir in examples/*; do
	    cd $dir
	    for g in *.go; do
		if [ $g = "types_${GOOS}.go" -o $g = "ztypes_${GOOSARCH}.go" ]; then
		    continue
		fi
		echo "building ${g}..."
		go build $g ztypes_${GOOSARCH}.go
	    done
	    cd $cwd
	done
	;;
    clean)
	find examples -type f | grep -v "\.go$" | xargs rm
	;;
esac

