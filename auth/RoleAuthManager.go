package auth

import (
	"net/http"
	"strings"

	cdata "github.com/pip-services3-gox/pip-services3-commons-gox/data"
	cerr "github.com/pip-services3-gox/pip-services3-commons-gox/errors"
	services "github.com/pip-services3-gox/pip-services3-rpc-gox/services"
)

type RoleAuthManager struct {
}

func (c *RoleAuthManager) UserInRoles(roles []string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {

		user, ok := req.Context().Value(AuthUser).(cdata.AnyValueMap)
		if !ok {
			services.HttpResponseSender.SendError(
				res, req,
				cerr.NewUnauthorizedError("", "NOT_SIGNED",
					"User must be signed in to perform this operation").WithStatus(401))
		} else {
			authorized := false
			userRoles := user.GetAsArray("roles")

			if userRoles == nil {
				services.HttpResponseSender.SendError(
					res, req,
					cerr.NewUnauthorizedError("", "NOT_SIGNED",
						"User must be signed in to perform this operation").WithStatus(401))
				return
			}

			for _, role := range roles {
				for _, userRole := range userRoles.Value() {
					r, ok := userRole.(string)
					if ok && role == r {
						authorized = true
					}
				}
			}

			if !authorized {
				services.HttpResponseSender.SendError(
					res, req,
					cerr.NewUnauthorizedError(
						"", "NOT_IN_ROLE",
						"User must be "+strings.Join(roles, " or ")+" to perform this operation").WithDetails(string(AuthRoles), roles).WithStatus(403))
			} else {
				next.ServeHTTP(res, req)
			}
		}
	}
}

func (c *RoleAuthManager) UserInRole(role string) func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return c.UserInRoles([]string{role})
}

func (c *RoleAuthManager) Admin() func(res http.ResponseWriter, req *http.Request, next http.HandlerFunc) {
	return c.UserInRole(string(AuthAdmin))
}
