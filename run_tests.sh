#!/bin/bash

# @copilot-script run_tests.sh 
# @description A helper script to run all tests in the HAWS project and provide a summary
# @author HAWS Team
# @date May 30, 2025
# 
# Usage:
#   ./run_tests.sh           - Run all tests in the project
#   ./run_tests.sh verbose   - Run all tests with detailed output
#   ./run_tests.sh coverage  - Run tests with code coverage reporting
#   ./run_tests.sh v-cov     - Run verbose tests with code coverage
#
# Exit codes:
#   0 - All tests passed or were skipped
#   1 - One or more tests failed

# Process command line arguments
VERBOSE=false
COVERAGE=false

if [ "$1" == "verbose" ] || [ "$1" == "v-cov" ]; then
  VERBOSE=true
fi

if [ "$1" == "coverage" ] || [ "$1" == "v-cov" ]; then
  COVERAGE=true
fi

# @copilot set up test command based on options
# @description Configure the go test command based on command line arguments
TEST_CMD="go test"
[ "$VERBOSE" == "true" ] && TEST_CMD="$TEST_CMD -v"
[ "$COVERAGE" == "true" ] && TEST_CMD="$TEST_CMD -coverprofile=coverage.out"
TEST_CMD="$TEST_CMD ./..."

echo "Running HAWS test suite..."
[ "$VERBOSE" == "true" ] && echo "Verbose mode: enabled"
[ "$COVERAGE" == "true" ] && echo "Coverage reporting: enabled"

# @copilot runs all tests and captures both stdout and stderr
# @description Use go test to run all tests in the project and capture the output
result=$($TEST_CMD 2>&1)
exit_code=$?

# Print the test output
echo "$result"

# @copilot count test results using grep
# @description Count the number of passing, failing, and skipped tests using grep
passed=$(echo "$result" | grep -- "--- PASS:" | wc -l | xargs)
failed=$(echo "$result" | grep -- "--- FAIL:" | wc -l | xargs)
skipped=$(echo "$result" | grep -- "--- SKIP:" | wc -l | xargs)

echo ""
echo "Test Results Summary:"
echo "✅ Passed:  $passed"
echo "❌ Failed:  $failed"
echo "⏭️  Skipped: $skipped"

# @copilot check test results and display coverage if requested
# @description Check if any tests failed and show code coverage if requested
if [ $exit_code -eq 0 ]; then
  echo "✅ All tests passed or were skipped!"
  
  # If coverage was requested, display the coverage report
  if [ "$COVERAGE" == "true" ]; then
    echo ""
    echo "Code Coverage Report:"
    go tool cover -func=coverage.out
    
    # Generate HTML report if requested
    if [ "$VERBOSE" == "true" ]; then
      echo ""
      echo "Generating HTML coverage report (coverage.html)..."
      go tool cover -html=coverage.out -o coverage.html
      echo "Open coverage.html in your browser to view detailed coverage report"
    fi
  fi
  
  exit 0
else
  echo "❌ Some tests failed!"
  exit 1
fi
