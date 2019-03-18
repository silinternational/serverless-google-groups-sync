#!/usr/bin/env bash

# Exit script with error if any step fails.
set -e

# Run tests sequentially for fixture integrity
testDirs=`find -name '*_test.go' -not -path "*vendor*" -printf '%h\n' | sort -u`
i=0
for testDir in $testDirs; do
    go test $testDir
done