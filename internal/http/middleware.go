package http

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"

	"github.com/getkin/kin-openapi/openapi3"
	"github.com/getkin/kin-openapi/openapi3filter"
	"github.com/getkin/kin-openapi/routers"
	"github.com/getkin/kin-openapi/routers/gorillamux"
	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

func Middleware() (mux.MiddlewareFunc, error) {
	openapi3.SchemaErrorDetailsDisabled = true
	openapi3.DefineStringFormatCallback("uuid", func(value string) error {
		_, err := uuid.Parse(value)
		return err
	})

	openapi3.DefineStringFormatCallback("json", func(value string) error {
		var out any
		return json.Unmarshal([]byte(value), &out)
	})

	ctx := context.Background()
	loader := &openapi3.Loader{Context: ctx, IsExternalRefsAllowed: true}
	doc, err := loader.LoadFromFile("./api/openapi.yaml")
	if err != nil {
		return nil, err
	}
	if err = doc.Validate(ctx, openapi3.EnableSchemaFormatValidation()); err != nil {
		return nil, err
	}
	router, err := gorillamux.NewRouter(doc)
	if err != nil {
		return nil, err
	}

	options := openapi3filter.Options{
		AuthenticationFunc: openapi3filter.NoopAuthenticationFunc,
	}

	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route, pathParams, err := router.FindRoute(r)
			if errors.Is(err, routers.ErrPathNotFound) {
				h.ServeHTTP(w, r)
				return
			}

			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				return
			}

			requestValidationInput := &openapi3filter.RequestValidationInput{
				Request:    r,
				PathParams: pathParams,
				Route:      route,
				Options:    &options,
			}
			if err := openapi3filter.ValidateRequest(r.Context(), requestValidationInput); err != nil {
				_, _ = w.Write([]byte(err.Error()))
				w.WriteHeader(http.StatusBadRequest)

				return
			}

			h.ServeHTTP(w, r)
		})
	}, nil
}
