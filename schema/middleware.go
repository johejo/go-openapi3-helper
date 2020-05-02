package schema

import (
	"net/http"

	"github.com/getkin/kin-openapi/openapi3filter"
)

// RequestValidator returns a http middleware that validates request.
func RequestValidator(router *openapi3filter.Router, validationErrHook func(w http.ResponseWriter, r *http.Request, err error)) func(next http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			route, pathParams, err := router.FindRoute(r.Method, r.URL)
			if err != nil {
				http.NotFound(w, r)
				return
			}
			input := &openapi3filter.RequestValidationInput{Request: r, PathParams: pathParams, Route: route}
			if err := openapi3filter.ValidateRequest(r.Context(), input); err != nil {
				if validationErrHook != nil {
					validationErrHook(w, r, err)
					return
				}
				w.WriteHeader(http.StatusBadRequest)
				return
			}
			next.ServeHTTP(w, r)
		})
	}
}
