#!/bin/bash
# Run acceptance tests for the Ory Terraform provider
#
# Usage:
#   ./scripts/test.sh              # Run all acceptance tests
#   ./scripts/test.sh identity     # Run only identity tests
#   ./scripts/test.sh oauth2       # Run only OAuth2 client tests
#   ./scripts/test.sh organization # Run only organization tests

set -e

# Load .env if it exists
if [ -f .env ]; then
    echo "Loading environment from .env..."
    set -a
    source .env
    set +a
fi

# Also check parent directory for shared .env (for mono-repo setups)
if [ -f ../.env ]; then
    echo "Loading environment from ../.env..."
    set -a
    source ../.env
    set +a
fi

# Validate required environment variables
missing_vars=()
[ -z "$ORY_PROJECT_API_KEY" ] && missing_vars+=("ORY_PROJECT_API_KEY")
[ -z "$ORY_PROJECT_ID" ] && missing_vars+=("ORY_PROJECT_ID")
[ -z "$ORY_PROJECT_SLUG" ] && missing_vars+=("ORY_PROJECT_SLUG")

if [ ${#missing_vars[@]} -ne 0 ]; then
    echo "Error: Missing required environment variables:"
    for var in "${missing_vars[@]}"; do
        echo "  - $var"
    done
    echo ""
    echo "Please copy .env.example to .env and fill in your credentials."
    exit 1
fi

# Determine which tests to run
TEST_PATH="./..."
case "${1:-all}" in
    identity)
        TEST_PATH="./internal/resources/identity/..."
        ;;
    oauth2|oauth2client)
        TEST_PATH="./internal/resources/oauth2client/..."
        ;;
    organization|org)
        TEST_PATH="./internal/resources/organization/..."
        ;;
    all)
        TEST_PATH="./..."
        ;;
    *)
        echo "Unknown test target: $1"
        echo "Usage: $0 [identity|oauth2|organization|all]"
        exit 1
        ;;
esac

echo "Running acceptance tests for: $TEST_PATH"
echo ""

# Run tests with TF_ACC=1 to enable acceptance tests
# Use -p 1 to run packages serially and avoid conflicts when hitting the same Ory project
TF_ACC=1 go test "$TEST_PATH" -v -timeout 30m -p 1
