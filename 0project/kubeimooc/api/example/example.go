package example

import "github.com/gin-gonic/gin"

type ExampleApi struct {
}

func (*ExampleApi) ExampleTest(c *gin.Context) {

	c.JSON(200, gin.H{
		"message": "pong",
	})
}
