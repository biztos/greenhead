#!/bin/bash
#
# Hooray for robots!  ChatGPT wrote the first draft here.
#
# Usage: ./licenses.sh <source_dir> <output_file>


set -e

if [[ $# -ne 2 ]]; then
    echo "Usage: $0 <source_dir> <output_file>"
    exit 1
fi

SRC_DIR="$1"
OUTPUT="$2"

# Define file patterns to include
LICENSE_FILES="LICENSE LICENSE.txt LICENSE.md COPYING COPYING.txt COPYING.md"
NOTICE_FILES="NOTICE NOTICE.txt NOTICE.md"

# Start fresh
echo "# Greenhead – Full Licenses" > "$OUTPUT"
echo >> "$OUTPUT"
echo "## Greenhead itself uses the MIT License." >> "$OUTPUT"
echo >> "$OUTPUT"
cat LICENSE | sed 's/^/    /' >> "$OUTPUT"

echo -e "\n\n## Third-Party Licenses\n" >> "$OUTPUT"

# Loop over each subdirectory in source
find "$SRC_DIR" -type d | while read -r dir; do
    found_any=
    output_section=""

    for file in $LICENSE_FILES $NOTICE_FILES; do
        filepath="$dir/$file"
        if [[ -f "$filepath" ]]; then
            if [[ -z "$found_any" ]]; then
                output_section+="\n### ${dir#$SRC_DIR/}\n"
                found_any=1
            fi
            output_section+="\n$file\n\n"
            output_section+="$(cat "$filepath" | sed 's/^/    /')\n"
        fi
    done

    if [[ -n "$found_any" ]]; then
        echo -e "$output_section" >> "$OUTPUT"
    fi
done
