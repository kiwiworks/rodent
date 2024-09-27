package api

import (
	"context"
	"fmt"
	"runtime"
	"strings"

	"github.com/danielgtaylor/huma/v2"

	"github.com/kiwiworks/rodent/system/opt"
	"github.com/kiwiworks/rodent/web/http"
)

type (
	Handler struct {
		Options Options
		Mount   func(api huma.API, config Config)
	}
	Options struct {
		Method          http.Method
		Path            string
		RegisterOas3    bool
		OperationId     string
		ContentType     string
		Tags            []string
		Protected       bool
		AuthProviders   []string
		OAuth2Providers map[string][]string
	}
)

func operationIdFromCaller(skip int) string {
	pc, _, _, ok := runtime.Caller(skip)
	if !ok {
		panic(fmt.Sprintf("call runtime.Caller() failed, this should not happen ever"))
	}
	callerFunc := runtime.FuncForPC(pc)
	callerName := callerFunc.Name()
	callerParts := strings.Split(callerName, "/")
	lastPart := callerParts[len(callerParts)-1]
	lastPartParts := strings.Split(lastPart, ".")
	packageName := lastPartParts[0]
	functionName := lastPartParts[1]
	return fmt.Sprintf("%s%s", packageName, functionName)
}

func NewHandler[Request any, Response any](
	method http.Method,
	path string,
	impl func(ctx context.Context, request *Request) (*Response, error),
	opts ...opt.Option[Options],
) *Handler {
	//todo this is wonky and should not be written like this
	operationId := operationIdFromCaller(3)
	if operationId == "reflectValue" {
		operationId = operationIdFromCaller(2)
	}
	fmt.Println(operationId)
	options := Options{
		Method:        method,
		Path:          path,
		RegisterOas3:  true,
		OperationId:   operationId,
		ContentType:   "application/json; charset=utf-8",
		Tags:          []string{},
		AuthProviders: []string{},
	}
	opt.Apply(&options, opts...)
	return &Handler{
		Options: options,
		Mount: func(api huma.API, config Config) {
			op := huma.Operation{
				Method:      options.Method.String(),
				Path:        options.Path,
				OperationID: options.OperationId,
				Tags:        options.Tags,
				Security:    []map[string][]string{},
			}
			if options.Protected {
				op.Security = append(op.Security, map[string][]string{
					"protected": {"read", "write"},
				})
			}
			huma.Register(api, op, func(ctx context.Context, i *Request) (*Response, error) {
				response, err := impl(ctx, i)
				return response, config.ErrorConverter(err)
			})
		},
	}
}
