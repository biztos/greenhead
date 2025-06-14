#!/bin/bash
set -e

# Usage: ./check-updated.sh <file> <file_or_directory>
# Fails if file is older than the target file or newest file in directory
FILE="$1"
TARGET="$2"

if [ ! -f "$FILE" ]; then
    echo "Error: File '$FILE' does not exist"
    exit 1
fi

if [ ! -e "$TARGET" ]; then
    echo "Error: '$TARGET' does not exist"
    exit 1
fi

# Detect stat command format
if stat -c '%Y' /dev/null >/dev/null 2>&1; then
    # GNU stat (Linux)
    STAT_FORMAT="-c %y"
else
    # BSD stat (macOS)
    STAT_FORMAT="-f %Sm"
fi

if [ -f "$TARGET" ]; then
    # Target is a file - simple comparison
    COMPARE_FILE="$TARGET"
elif [ -d "$TARGET" ]; then
    # Target is a directory - find newest file
    if stat -c '%Y' /dev/null >/dev/null 2>&1; then
        # GNU stat (Linux)
        COMPARE_FILE=$(find "$TARGET" -type f -exec stat -c '%Y %n' {} \; 2>/dev/null | sort -n | tail -1 | cut -d' ' -f2-)
    else
        # BSD stat (macOS)
        COMPARE_FILE=$(find "$TARGET" -type f -exec stat -f '%m %N' {} \; 2>/dev/null | sort -n | tail -1 | cut -d' ' -f2-)
    fi
    
    if [ -z "$COMPARE_FILE" ]; then
        echo "No files found in directory '$TARGET'"
        exit 0
    fi
else
    echo "Error: '$TARGET' is neither a file nor a directory"
    exit 1
fi

# Compare timestamps
if [ "$FILE" -ot "$COMPARE_FILE" ]; then
    echo "FAIL: '$FILE' is older than '$COMPARE_FILE'"
    echo "File timestamp:    $(stat $STAT_FORMAT "$FILE")"
    echo "Compare timestamp: $(stat $STAT_FORMAT "$COMPARE_FILE")"
    exit 1
fi

echo "OK: '$FILE' is up to date"