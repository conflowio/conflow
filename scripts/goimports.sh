#!/bin/bash

set -euo pipefail

PROJECT_DIR="$( cd "$( dirname "${BASH_SOURCE[0]}" )"/.. >/dev/null 2>&1 && pwd )"

remove_lines_file() {
  if sed --version 2>&1 | grep -q "GNU sed"; then
    sed -i '
      /^import/,/)/ {
        /^$/ d
      }
    ' "${1}"
  else
    sed -i .orig '
      /^import/,/)/ {
        /^$/ d
      }
    ' "${1}"
    rm -f -- "${1}.orig"
  fi
}

cd "${PROJECT_DIR}"

GOBIN="$PROJECT_DIR/bin" go install golang.org/x/tools/cmd/goimports

# shellcheck disable=SC2016
for path in $(go list -f '{{ $dir := .Dir }}{{ range .GoFiles }}{{ printf "%s/%s\n" $dir . }}{{ end }}{{ range .CgoFiles }}{{ printf "%s/%s\n" $dir . }}{{ end }}{{ range .TestGoFiles }}{{ printf "%s/%s\n" $dir . }}{{ end }}{{ range .XTestGoFiles }}{{ printf "%s/%s\n" $dir . }}{{ end }}' ./...); do
  if [[ "${path}" == *.cf.go ]] || [[ "${path}" == */fake_*.go ]]; then
    continue
  fi

  if [[ "${path}" == */interpreter_template.go ]]; then
    continue
  fi

  remove_lines_file "${path}"

  "${PROJECT_DIR}"/bin/goimports -local github.com/conflowio/conflow -w "${path}"
done
