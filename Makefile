SHELL:=/bin/bash
build-recipe-tool:
	cp ${filename} ./input.json
	docker build -t recipe --build-arg filename='./input.json' .
run-recipe-tool:
	docker run --env searchby=${searchby} --env postcode=${postcode} --env deliverytime='${deliverytime}' recipe
test-recipe-tool:
	cp ${filename} ./input.json
	go test -v ./...
