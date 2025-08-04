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

# Ensure the output file exists with a mode line
if [ ! -f "$OUTPUT" ]; then
  echo "mode: atomic" > "$OUTPUT"
fi

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
