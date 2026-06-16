#!/usr/bin/env sh

set -e

echo "Checking '$*'"

before=$(mktemp)
after=$(mktemp)
trap 'rm -f "$before" "$after"' EXIT

git status -s >"$before"

"$@"

git status -s >"$after"
if ! diff -q "$before" "$after" >/dev/null 2>&1; then
	echo "There are changed files after running '$*'"
	diff -u "$before" "$after" || true
	exit 1
fi
