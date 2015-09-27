import (
	gin "github.com/gin-gonic/gin"
  http "net/http"
)

func init() {
  r := gin.New()

  api := r.Group("/api")

  http.Handle("/", r);
}
