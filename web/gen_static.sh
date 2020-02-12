#!/bin/bash

files=(index config stock outbound)

for file in ${files[*]}
do 
    ./data.sh $file | mustache template/$file.mustache > src/$file.html
done
