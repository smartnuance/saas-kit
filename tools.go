//go:build tools

//go:generate bash -c "protoc --go_out=. --go_opt=module=github.com/smartnuance/saas-kit --go-grpc_out=. --go-grpc_opt=module=github.com/smartnuance/saas-kit ./proto/*"

package tools

import (
	_ "github.com/ahmetb/govvv"
	_ "github.com/golang/mock/mockgen"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway"
	_ "github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-openapiv2"
	_ "github.com/volatiletech/sqlboiler/v4"
	_ "github.com/volatiletech/sqlboiler/v4/drivers/sqlboiler-psql"
	_ "google.golang.org/grpc/cmd/protoc-gen-go-grpc"
	_ "google.golang.org/protobuf/cmd/protoc-gen-go"
)
