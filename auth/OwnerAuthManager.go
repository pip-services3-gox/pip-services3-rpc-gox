package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	services "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

type OwnerAuthManager struct {
}

func (c *OwnerAuthManager) Owner(idParam string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if idParam == "" {
		idParam = string(AuthUserId)
	}
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

		_, ok := req.Context().Value(AuthUser).(cdata.AnyValueMap)

		if !ok {
			services.HttpResponseSender.SendError(
				res, req,
				cerr.NewUnauthorizedError("",
					"NOT_SIGNED",
					"User must be signed in to perform this operation",
				).WithStatus(401),
			)
		} else {
			userId := req.URL.Query().Get(idParam)
			if userId == "" {
				userId = mux.Vars(req)[idParam]
			}

			reqUserId, ok := req.Context().Value(AuthUserId).(string)
			if !ok || reqUserId != userId {
				services.HttpResponseSender.SendError(
					res, req,
					cerr.NewUnauthorizedError(
						"", "FORBIDDEN",
						"Only data owner can perform this operation",
					).WithStatus(403),
				)
			} else {
				next.ServeHTTP(res, req)
			}
		}
	}
}

func (c *OwnerAuthManager) OwnerOrAdmin(idParam string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if idParam == "" {
		idParam = string(AuthUserId)
	}
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

		user, ok := req.Context().Value(AuthUser).(cdata.AnyValueMap)

		if !ok {
			services.HttpResponseSender.SendError(
				res, req,
				cerr.NewUnauthorizedError("",
					"NOT_SIGNED",
					"User must be signed in to perform this operation",
				).WithStatus(401),
			)
		} else {

			userId := req.URL.Query().Get(idParam)
			if userId == "" {
				userId = mux.Vars(req)[idParam]
			}
			roles := user.GetAsArray(string(AuthRoles))
			admin := false
			for _, role := range roles.Value() {
				r, ok := role.(string)
				if ok && r == string(AuthAdmin) {
					admin = true
					break
				}
			}

			reqUserId, ok := req.Context().Value(AuthUserId).(string)
			if !ok || reqUserId != userId && !admin {
				services.HttpResponseSender.SendError(
					res, req,
					cerr.NewUnauthorizedError("",
						"FORBIDDEN",
						"Only data owner can perform this operation",
					).WithStatus(403),
				)
			} else {
				next.ServeHTTP(res, req)
			}
		}
	}
}
