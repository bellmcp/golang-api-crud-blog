package middleware

import (
	"course-go/models"
	"net/http"

	"github.com/casbin/casbin"
	"github.com/gin-gonic/gin"
)

func Authorize() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		// authentication
		user, ok := ctx.Get("sub")
		if !ok { // no sub = not logged in
			ctx.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"}) // do this and abort the next handler (see routes file)
			return
		}

		// authentication
		enforcer := casbin.NewEnforcer("config/acl_model.conf", "config/policy.csv")
		ok = enforcer.Enforce(user.(*models.User), ctx.Request.URL.Path, ctx.Request.Method) // sub (User), obj (Path), act (HTTP Method)

		if !ok {
			ctx.AbortWithStatusJSON(http.StatusForbidden, gin.H{ // do this and abort the next handler (see routes file)
				"error": "you are not allowed to access this resource",
			})
			return
		}

		ctx.Next() // allow to run next handler
	}
}
