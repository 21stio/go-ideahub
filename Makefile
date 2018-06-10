SOURCE_MAKE=source ./make/make.sh

bbi: build_build_image
build_build_image:
	${SOURCE_MAKE} && build_build_image

b: build
build:
	${SOURCE_MAKE} && build

r: render
render:
	${SOURCE_MAKE} && render

m: migrate
migrate:
	${SOURCE_MAKE} && migrate

dgeo: download_geodb
download_geodb:
	${SOURCE_MAKE} && download_geodb