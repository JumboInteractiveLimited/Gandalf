#!/bin/sh
set -euf

build() {
	docker-compose build --force-rm
}

mmock() {
	docker-compose up --remove-orphans -d mock
}

run() {
	docker-compose up --force-recreate --remove-orphans --exit-code-from gandalf gandalf
}

down() {
	docker-compose down -v --remove-orphans
	docker-compose rm -svf
}

cleanup() {
	down
	export COMPOSE_FILE=
	export COMPOSE_PROJECT_NAME=
}

setup() {
	export COMPOSE_FILE=examples/prototype/docker-compose.yml
	export COMPOSE_PROJECT_NAME=gandalf
	trap cleanup EXIT
}

setup && build && mmock && run
