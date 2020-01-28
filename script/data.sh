#!/bin/bash

current="$1"

declare -A section_title=( 
  ["index"]="Inventory" 
  ["inboundconfig"]="Finished products"
  ["outboundconfig"]="Outbound"
  ["stock"]="Stock"
)

declare -A section_info=( 
  ["index"]="Configure inventory types" 
  ["inboundconfig"]="Configure inbound types"
  ["outboundconfig"]="Configure outbound types"
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
    { url: config.html, title: ${section_title['inboundconfig']}, selected: $(is_selected inboundconfig) },
    { url: outbound.htm, title: ${section_title['outboundconfig']}, selected: $(is_selected outboundconfig) },
    { url: stock.htm, title: ${section_title['stock']}, selected: $(is_selected stock) }
  ]

title: ${section_title[$current]}
description: ${section_info[$current]}
year: 2020
---
EOF
