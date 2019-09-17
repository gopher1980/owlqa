#!/usr/bin/env bash


for filename in ./tests/test*.json; do
    URL=http://127.0.0.1:9090/owl?s=$(openssl rand -base64 12)
    RESULT=./output/result_$(basename "$filename" .json).json
    echo $URL  " >> "   $filename " > " $RESULT
    curl -s -X POST -H "Content-Type: application/json" -d @$filename $URL | jq '.' > $RESULT &
done