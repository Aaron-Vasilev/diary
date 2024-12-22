#!/bin/bash

CHANGED_FILES=$(git diff --name-only | grep templ$)

if [ -z "$CHANGED_FILES" ]; then
  echo "No .templ files have changed."
  exit 0
fi

echo "$CHANGED_FILES" | while read -r file; do
  echo "Ordering file: $file"
  tcs ./public/styles/out.css "$file"
done

echo "Tailwind class order completed."

