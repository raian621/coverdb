.PHONY: codegen
codegen:
	oapi-codegen \
		-generate chi,spec,models \
		-package api -o api/gen.go openapi.yaml

.PHONY: build-client
build-client:
	cd client; npm run build
