package auth

import (
	"net/http"

	"github.com/gorilla/mux"
	cdata "github.com/pip-services3-go/pip-services3-commons-go/data"
	cerr "github.com/pip-services3-go/pip-services3-commons-go/errors"
	services "github.com/pip-services3-go/pip-services3-rpc-go/services"
)

type OwnerAuthManager struct {
}

func (c *OwnerAuthManager) Owner(idParam string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	if idParam == "" {
		idParam = "user_id"
	}
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

		_, ok := req.Context().Value("user").(cdata.AnyValueMap)

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

			reqUserId, ok := req.Context().Value("user_id").(string)
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
		idParam = "user_id"
	}
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

		user, ok := req.Context().Value("user").(cdata.AnyValueMap)

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
			roles := user.GetAsArray("roles")
			admin := false
			for _, role := range roles.Value() {
				r, ok := role.(string)
				if ok && r == "admin" {
					admin = true
					break
				}
			}

			reqUserId, ok := req.Context().Value("user_id").(string)
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
