#!/usr/bin/env bash

set -eo pipefail

function build_build_image {
	docker build \
		--file $(pwd)/make/build.Dockerfile \
		--tag 21stio/ideahub:build-latest \
		make

	docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD

	docker push 21stio/ideahub:build-latest
}

function build {
	wd=/go/src/github.com/21stio/go-ideahub

	docker run --rm -it -v $(pwd):${wd} -w ${wd} 21stio/ideahub:build-latest go build -o main

    docker build \
		--file $(pwd)/make/app.Dockerfile \
		--tag 21stio/ideahub:app-latest \
		.

	rm main

	docker login -u $DOCKER_USERNAME -p $DOCKER_PASSWORD

	docker push 21stio/ideahub:app-latest
}

function render {
	jinja2 --format=auto make/production.yaml .ignore/production.json
}

function download_geodb {
	geoip_database_location=http://geolite.maxmind.com/download/geoip/database/GeoLite2-City.mmdb.gz
	destination=geoip_database

	wget ${geoip_database_location} -O ${destination}.gz
	gzip --decompress --stdout ${destination}.gz > ${destination}.mmdb
	rm ${destination}.gz
}

function postgres_query {
    local query=$1
    local args=$2

    PGPASSWORD=${POSTGRES_PASSWORD} psql ${POSTGRES_URL} -c "${query}" --pset pager ${args}
}

function migrate {
	postgres_query "$(cat make/schema2018-06-08.sql)"
}
