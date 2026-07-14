package main

import (
	"battle-game/internal/handler"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	r.POST("/battle/start", handler.StartBattle)

	r.Run(":8080")
}
