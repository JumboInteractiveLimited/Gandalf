#!/bin/sh
set -e
cvrdir=coverage
cvrout="$cvrdir/report.out"
cvrhtm="$cvrdir/report.html"
mkdir -p "$cvrdir"

go test -gandalf.colour -cover -coverprofile "$cvrout" -v -bench . -benchmem github.com/JumboInteractiveLimited/Gandalf/
go tool cover -html "$cvrout" -o "$cvrdir/gandalf.html"
xdg-open "$cvrdir/gandalf.html" || true

go test -cover -coverprofile "$cvrout" -v -bench . -benchmem github.com/JumboInteractiveLimited/Gandalf/pathing
go tool cover -html "$cvrout" -o "$cvrdir/gandalf.pathing.html"
xdg-open "$cvrdir/gandalf.pathing.html" || true

go test -cover -coverprofile "$cvrout" -v -bench . -benchmem github.com/JumboInteractiveLimited/Gandalf/check
go tool cover -html "$cvrout" -o "$cvrdir/gandalf.pathing.check.html"
xdg-open "$cvrdir/gandalf.pathing.check.html" || true

gometalinter --enable-all github.com/JumboInteractiveLimited/Gandalf/...
