
.PHONY: proto
proto:
	protoc --proto_path=proto:. --micro_out=services --go_out=services services/proto/hicourt.proto
