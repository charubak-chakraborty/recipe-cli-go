SHELL:=/bin/bash
# build script
build-recipe-tool:
	cp ${filename} ./input.json
	docker build -t recipe --build-arg filename='./input.json' .
# run script
run-recipe-tool:
	docker run --env searchby=${searchby} --env postcode=${postcode} --env deliverytime='${deliverytime}' recipe
# test script
test-recipe-tool:
	cp ${filename} ./input.json
	go test -v ./...
