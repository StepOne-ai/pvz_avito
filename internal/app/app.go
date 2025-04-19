package main

import (
	"log"

	"github.com/StepOne-ai/pvz_avito/internal/db"
	"github.com/StepOne-ai/pvz_avito/internal/routes"
	"github.com/gin-gonic/gin"
)

func main() {
	if db.InitDB("database.db") != nil {
		log.Fatal("Error accessing db")
	}

	r := gin.Default()

	routes.SetupRoutes(r)

	r.Run(":8080")
}
