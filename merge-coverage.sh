#!/bin/bash
# merge-coverage.sh - Merges multiple Go coverage files properly
# Usage: ./merge-coverage.sh output.out file1.out file2.out [file3.out ...]

set -e

if [ "$#" -lt 2 ]; then
  echo "Usage: $0 output.out file1.out [file2.out ...]"
  exit 1
fi

OUTPUT=$1
shift

# Check if go is available
if ! command -v go &> /dev/null; then
    echo "Error: go command not found"
    exit 1
fi

# Determine the mode from the first input file
MODE="set"
for f in "$@"; do
  if [ -f "$f" ] && [ -s "$f" ]; then
    FIRST_LINE=$(head -n 1 "$f")
    if [[ "$FIRST_LINE" == mode:* ]]; then
      MODE=$(echo "$FIRST_LINE" | cut -d' ' -f2)
      echo "Using coverage mode: $MODE (from $f)"
      break
    fi
  fi
done

# Create temporary directory for processing
TEMP_DIR=$(mktemp -d)
trap "rm -rf $TEMP_DIR" EXIT

# Create output file with proper mode
echo "mode: $MODE" > "$OUTPUT"

# Collect all coverage entries in a temporary file
TEMP_COMBINED="$TEMP_DIR/combined.tmp"
touch "$TEMP_COMBINED"

# Process each input file
for f in "$@"; do
  if [ -f "$f" ] && [ -s "$f" ]; then
    echo "Processing $f..."
    # Check if file has actual coverage data (more than just mode line)
    if [ $(wc -l < "$f") -gt 1 ]; then
      # Skip the mode line and append coverage data
      tail -n +2 "$f" >> "$TEMP_COMBINED"
    else
      echo "Warning: $f contains no coverage data"
    fi
  else
    echo "Warning: File $f not found or empty, skipping"
  fi
done

# Sort and deduplicate coverage entries
# Coverage lines have format: filepath:startLine.startCol,endLine.endCol numStmts count
# We need to handle potential duplicates by using the maximum count for each location
if [ -s "$TEMP_COMBINED" ]; then
  # Sort by filepath and line numbers, then process duplicates
  sort "$TEMP_COMBINED" | awk '
  {
    # Extract the location part (everything before the last two spaces)
    # Format: filepath:startLine.startCol,endLine.endCol numStmts count
    n = split($0, parts, " ")
    if (n >= 3) {
      location = ""
      for (i = 1; i <= n-2; i++) {
        if (i > 1) location = location " "
        location = location parts[i]
      }
      stmts = parts[n-1]
      count = parts[n]
      
      # Keep the maximum count for each location
      if (location in max_count) {
        if (count > max_count[location]) {
          max_count[location] = count
          stmt_count[location] = stmts
        }
      } else {
        max_count[location] = count
        stmt_count[location] = stmts
      }
    }
  }
  END {
    for (loc in max_count) {
      print loc " " stmt_count[loc] " " max_count[loc]
    }
  }
  ' >> "$OUTPUT"
else
  echo "Warning: No coverage data found in any input files"
fi

echo "Successfully merged coverage files into $OUTPUT"

# Validate the output file
if [ -f "$OUTPUT" ]; then
  line_count=$(wc -l < "$OUTPUT")
  echo "Merged coverage file contains $line_count lines"
  if [ $line_count -gt 1 ]; then
    echo "Coverage merge completed successfully"
  else
    echo "Warning: Merged coverage file contains only header line"
  fi
else
  echo "Error: Failed to create merged coverage file"
  exit 1
fi
