package communication

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func addEndpoints(server *gin.Engine) *gin.RouterGroup {
	s := server.Group("api/v1")
	s.GET("/ping", ping)
	return s
}

func GetServer() *gin.Engine {
	server := gin.New()
	server.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Length", "Content-Type", "X-User-Platform", "X-User-Agent", "X-App-Version", "X-Access-Token"},
		ExposeHeaders:    []string{"Origin", "Content-Length", "Content-Type"},
		AllowCredentials: true,
	}))

	addEndpoints(server)
	return server
}
