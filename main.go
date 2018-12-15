package main

import (
	"github.com/gin-gonic/gin"
	auth "k8s.io/api/authentication/v1"
)

// do something here with `userInfo` to get actual values...
func getOrganizationAndProject(userInfo auth.UserInfo) (string, string) {
	return "acme", "test"
}

func main() {
	r := gin.Default()
	r.POST("/validate", validationHandler)
	r.POST("/mutate", mutationHandler)
	r.Run()
}
