#!/bin/bash

current="$1"

declare -A section_title=( 
  ["index"]="Inventory" 
  ["config"]="Finished products"
  ["outbound"]="Outbound"
  ["stock"]="Stock"
)

declare -A section_info=( 
  ["index"]="Configure inventory types" 
  ["config"]="Configure inbound types"
  ["outbound"]="Configure outbound types"
  ["stock"]="View/provision stock"
)

is_selected () {
    if [[ "$current" = "$1" ]] ; then
        echo "true"
    else
        echo "false"
    fi
}

cat << EOF
---
sections:
  [
    { url: index.html, title: ${section_title['index']}, selected: $(is_selected index) },
    { url: config.html, title: ${section_title['config']}, selected: $(is_selected config) },
    { url: outbound.html, title: ${section_title['outbound']}, selected: $(is_selected outbound) },
    { url: stock.html, title: ${section_title['stock']}, selected: $(is_selected stock) }
  ]

title: ${section_title[$current]}
description: ${section_info[$current]}
year: 2020
---
EOF
