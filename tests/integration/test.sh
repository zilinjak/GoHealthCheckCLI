#!/bin/bash
# This script is used to run integration tests for the project.
# We run tests in repeat to make sure that they are stable.

for i in {1..100}; do
    echo "Running test iteration $i"
    go clean -testcache
    go test -race .
    if [ $? -ne 0 ]; then
        echo "Integration tests failed on iteration $i"
        exit 1
    fi
done