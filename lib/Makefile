.PHONY: gen-pb
gen-pb: 
	find ./proto/domain -name *.proto -print0 | xargs -0 protoc --proto_path=./proto/domain --go_out=./proto/gen/pb --go_opt=paths=source_relative --go-grpc_out=./proto/gen/pb --go-grpc_opt=paths=source_relative

.PHONY: clear-pb
clear-pbs:
	rm -f ./proto/gen/pb/*.go