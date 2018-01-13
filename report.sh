#!/usr/bin/env bash
set -e
cvrdir=coverage
mkdir -p "$cvrdir"

declare -a pkgs=("" "pathing" "check")

for pkg in "${pkgs[@]}"; do
    cvrout="$cvrdir/report.${pkg:-root}.out"
    go test -gandalf.colour -cover -coverprofile "$cvrout" -v -bench . -benchmem "github.com/JumboInteractiveLimited/Gandalf/$pkg"
    go tool cover -html "$cvrout" -o "$cvrdir/gandalf.${pkg:-root}.html"
    xdg-open "$cvrdir/gandalf.${pkg:-root}.html" || true
done

gometalinter --enable-all --line-length=140 ./...
