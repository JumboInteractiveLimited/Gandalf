#!/usr/bin/env bash
set -e
cvrdir=coverage
mkdir -p "$cvrdir"

declare -a pkgs=("" "pathing" "check" "examples/goserver")

for pkg in "${pkgs[@]}"; do
    safepkg="$(sed 's/\//_/g' <<< "${pkg:-root}")"
    cvrout="$cvrdir/report.$safepkg.out"
    go test -cover -coverprofile "$cvrout" -v -benchtime 10ms -bench . -benchmem "github.com/JumboInteractiveLimited/Gandalf/$pkg"
    go tool cover -html "$cvrout" -o "$cvrdir/gandalf.$safepkg.html"
    xdg-open "$cvrdir/gandalf.$safepkg.html" || true
done

gometalinter --enable-all --line-length=140 ./...
