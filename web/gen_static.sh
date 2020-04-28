#!/bin/bash

echo -e \
"title: Stock management\n"\
"year: 2024" \
| mustache template/index.mustache > src/index.html
