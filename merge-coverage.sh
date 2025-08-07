#!/bin/bash
# merge-coverage.sh - Merges multiple Go coverage files
# Usage: ./merge-coverage.sh output.out file1.out file2.out [file3.out ...]

set -e

if [ "$#" -lt 2 ]; then
  echo "Usage: $0 output.out file1.out [file2.out ...]"
  exit 1
fi

OUTPUT=$1
shift

# Determine the mode from the first input file
MODE="atomic"
for f in "$@"; do
  if [ -f "$f" ]; then
    FIRST_LINE=$(head -n 1 "$f")
    if [[ "$FIRST_LINE" == mode:* ]]; then
      MODE=$(echo "$FIRST_LINE" | cut -d' ' -f2)
      echo "Using coverage mode: $MODE (from $f)"
      break
    fi
  fi
done

# Create output file with proper mode
echo "mode: $MODE" > "$OUTPUT"

# Append the content from all files (excluding mode line)
for f in "$@"; do
  if [ -f "$f" ]; then
    echo "Processing $f..."
    # Skip the mode line (first line) and append to output
    tail -n +2 "$f" >> "$OUTPUT"
  else
    echo "Warning: File $f not found, skipping"
  fi
done

echo "Successfully merged coverage files into $OUTPUT"
