#!/bin/sh

GOOS=${1:-linux}
GOARCH=${2:-amd64}
GOOSARCH="${GOOS}_${GOARCH}"

cwd=`pwd`
cmd=${3:-"build"}

case "$cmd" in
    build)
	find examples -name \*.go | while read pathname; do
	    gopt=""
	    case "$pathname" in
		*types_${GOOS}.go) continue;;
		*ztypes_${GOOSARCH}.go) continue;;
	    esac

	    cd `dirname $pathname`
	    g=`basename $pathname`
	    echo "building ${pathname}..."
	    if [ -f ztypes_${GOOSARCH}.go ]; then
		go build $gopt $g ztypes_${GOOSARCH}.go
	    else
		go build $gopt $g
	    fi
	    cd $cwd
	done
	;;
    clean)
	find examples -type f | grep -v "\.go$" | xargs rm
	;;
esac

