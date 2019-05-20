.: generate

generate:
	@echo "Generating schema..."
	@protoc \
		--proto_path=${GOPATH}/src \
		--proto_path=${GOPATH}/src/github.com/google/protobuf/src \
		--proto_path=. \
		--go_out=. \
		--js_out=../../client/src/schema/ \
		--govalidators_out=. \
		*.proto