#!/bin/bash

files=(index config stock outbound)

for file in ${files[*]}
do 
    ./data.sh $file | mustache $file.mustache > $file.html
done
