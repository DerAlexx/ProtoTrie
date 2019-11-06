regenerate-docker:
	docker run --rm -it \
		-v `pwd`:/messages \
		-w /messages \
		obraun/vss-protoactor-jenkins \
		make regenerate

regenerate:
	go install -v github.com/gogo/protobuf/protoc-gen-gogoslick
#	go install -v github.com/AsynkronIT/protoactor-go/actor
	protoc -I=. -I=${GOPATH}/src -I=${GOPATH}/src/github.com/gogo/protobuf/protobuf \
		--gogoslick_out=plugins=grpc:. tree.proto
