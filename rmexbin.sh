#!/bin/sh
find examples -type f | grep -v "\.go$" | xargs rm
