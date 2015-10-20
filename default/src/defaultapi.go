package defaultapp

import (
	gin "github.com/gamescores/gin"
	http "net/http"
)

func init() {
	r := gin.New()

	api := r.Group("/api")
	api.GET("/", func(c *gin.Context) {
		c.String(200, "HELLO")
	})

	http.Handle("/", r);
}
