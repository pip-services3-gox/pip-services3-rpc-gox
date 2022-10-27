package auth

import (
	"net/http"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	services "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

type BasicAuthManager struct {
}

func (c *BasicAuthManager) Anybody() func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		next.ServeHTTP(res, req)
	}
}

func (c *BasicAuthManager) Signed() func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
		_, ok := req.Context().Value(User).(cdata.AnyValueMap)
		if !ok {
			services.HttpResponseSender.SendError(
				res, req,
				cerr.NewUnauthorizedError("",
					"NOT_SIGNED",
					"User must be signed in to perform this operation",
				).WithStatus(401),
			)
		} else {
			next.ServeHTTP(res, req)
		}
	}
}
