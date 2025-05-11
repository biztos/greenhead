#!/bin/bash
#
# jq-pretty-build.sh -- prettify any JSON found in the build dir
# 
# This is (very!) useful when debugging stuff with the dump-dir option.
find build -type f -name "*.json" ! -name "*-pretty.json" -exec sh -c 'for f; do jq . "$f" > "${f%.json}-pretty.json"; done' sh {} +

