module grpc-gateway

go 1.19

replace test-service => ./test-service

require (
	github.com/stretchr/testify v1.8.1
	google.golang.org/grpc v1.53.0
	google.golang.org/protobuf v1.28.1
	test-service v0.0.0-00010101000000-000000000000
)

require (
	github.com/davecgh/go-spew v1.1.1 // indirect
	github.com/golang/protobuf v1.5.2 // indirect
	github.com/pmezard/go-difflib v1.0.0 // indirect
	golang.org/x/net v0.6.0 // indirect
	golang.org/x/sys v0.5.0 // indirect
	golang.org/x/text v0.7.0 // indirect
	google.golang.org/genproto v0.0.0-20230209215440-0dfe4f8abfcc // indirect
	gopkg.in/yaml.v3 v3.0.1 // indirect
)
