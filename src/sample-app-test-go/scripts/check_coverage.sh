#!/bin/bash

# Define the minimum coverage required
MIN_COVERAGE=80

# Extract the coverage percentage from the coverage.out file
COVERAGE=$(go tool cover -func=coverage.out | grep total | awk '{print $3}' | sed 's/%//')

# Check if coverage is at least the minimum required
if (( $(echo "$COVERAGE >= $MIN_COVERAGE" | bc -l) )); then
  echo "Test coverage is sufficient: ${COVERAGE}%"
else
  echo "Test coverage is insufficient: ${COVERAGE}% (minimum required is ${MIN_COVERAGE}%)"
  exit 1
fi
