build_info_flag = -ldflags "-X main.buildVersion=$$(cat cmd/client/version) -X main.buildDate=$$(date +'%d/%m/%Y') -X main.defaultConfigPath=./config.json"
client_app = cmd/client/main.go

gen_proto:
	rm api/proto/*.go || true
# TODO don't know why internal path to protoeditor is required
	protoc -I C:/Users/mixa1/AppData/Local/JetBrains/GoLand2022.2/protoeditor -I . --go_out=. --go_opt=paths=source_relative --go-grpc_out=. --go-grpc_opt=paths=source_relative api/proto/*.proto

client_build_windows:
	GOOS=windows GOARCH=amd64 go build -o bin/client/gophkeeper.exe $(build_info_flag) $(client_app)

client_build_osx:
	GOOS=darwin GOARCH=arm64 go build -o bin/client/gophkeeper_osx $(build_info_flag) $(client_app)

client_build_linux:
	GOOS=linux GOARCH=amd64 go build -o bin/client/gophkeeper_linux $(build_info_flag) $(client_app)

cert:
    ./cert/gen.sh