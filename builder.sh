#!/bin/sh
set -e

if [ ! -d .git ]; then
    echo "Error: This script should be run from the root of the repository!" >&2
    exit 1
fi

if [ ! -f .git/hooks/pre-commit ]; then
    if command -v pre-commit >/dev/null 2>&1; then
        echo "Installing pre-commit hook..."
        pre-commit install
        echo "pre-commit hook installed successfully."
    else
        echo "pre-commit not installed!"
        echo "Please install it, then run 'pre-commit install'."
        exit 1
    fi
fi

echo "Repo is setup properly."
