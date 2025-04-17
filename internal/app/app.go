package main

import (
	"fmt"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/StepOne-ai/pvz_avito/internal/db"
	"github.com/StepOne-ai/pvz_avito/internal/models"
	"github.com/gin-gonic/gin"
)

func main() {
	if db.InitDB("database.db") != nil {
		log.Fatal("Error accessing db")
	}

	r := gin.Default()

	r.POST("/dummyLogin", func(c *gin.Context) {
		var req struct {
			Role string `json:"role"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
			return
		}
		token := models.Token(fmt.Sprintf("token-%d", time.Now().Unix()))
		c.JSON(http.StatusOK, token)
	})

	r.POST("/register", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
			Role     string `json:"role"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
			return
		}
		id := fmt.Sprintf("user-%d", time.Now().Unix())
		err := db.CreateUser(id, req.Email, req.Password, req.Role)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, models.User{ID: id, Email: req.Email, Role: req.Role})
	})

	r.POST("/login", func(c *gin.Context) {
		var req struct {
			Email    string `json:"email"`
			Password string `json:"password"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
			return
		}
		user, err := db.GetUserByEmail(req.Email)
		if err != nil || user["password"] != req.Password {
			c.JSON(http.StatusUnauthorized, models.Error{Message: "Invalid credentials"})
			return
		}
		token := models.Token(fmt.Sprintf("token-%d", time.Now().Unix()))
		c.JSON(http.StatusOK, token)
	})

	authMiddleware := func(c *gin.Context) {
		token := c.GetHeader("Authorization")
		if len(token) < 7 || token[:7] != "Bearer " {
			c.JSON(http.StatusForbidden, models.Error{Message: "Access denied"})
			c.Abort()
			return
		}
		c.Next()
	}

	r.POST("/pvz", authMiddleware, func(c *gin.Context) {
		var req models.PVZ
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
			return
		}
		req.ID = fmt.Sprintf("pvz-%d", time.Now().Unix())
		req.RegistrationDate = time.Now()
		err := db.CreatePVZ(req.ID, req.City, req.RegistrationDate.Format(time.RFC3339))
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, req)
	})

	r.GET("/pvz", authMiddleware, func(c *gin.Context) {
		startDate := c.Query("startDate")
		endDate := c.Query("endDate")
		page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
		limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

		pvzs, err := db.GetPVZsFiltered(startDate, endDate, page, limit)
		if err != nil {
			c.JSON(http.StatusInternalServerError, models.Error{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, pvzs)
	})

	r.POST("/pvz/:pvzId/close_last_reception", authMiddleware, func(c *gin.Context) {
		pvzId := c.Param("pvzId")
		err := db.CloseLastReception(pvzId)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Reception closed"})
	})

	r.POST("/pvz/:pvzId/delete_last_product", authMiddleware, func(c *gin.Context) {
		pvzId := c.Param("pvzId")
		err := db.DeleteLastProduct(pvzId)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
			return
		}
		c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
	})

	r.POST("/receptions", authMiddleware, func(c *gin.Context) {
		var req struct {
			PvzId string `json:"pvzId"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
			return
		}
		reception := models.Reception{
			ID:       fmt.Sprintf("reception-%d", time.Now().Unix()),
			DateTime: time.Now(),
			PvzId:    req.PvzId,
			Status:   "in_progress",
		}
		err := db.CreateReception(reception.ID, reception.DateTime.Format(time.RFC3339), reception.PvzId, reception.Status)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, reception)
	})

	r.POST("/products", authMiddleware, func(c *gin.Context) {
		var req struct {
			Type  string `json:"type"`
			PvzId string `json:"pvzId"`
		}
		if err := c.ShouldBindJSON(&req); err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
			return
		}
		product := models.Product{
			ID:          fmt.Sprintf("product-%d", time.Now().Unix()),
			DateTime:    time.Now(),
			Type:        req.Type,
			ReceptionId: "current-reception-id", // Replace with actual reception ID logic
		}
		err := db.CreateProduct(product.ID, product.DateTime.Format(time.RFC3339), product.Type, product.ReceptionId)
		if err != nil {
			c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
			return
		}
		c.JSON(http.StatusCreated, product)
	})

	r.Run(":8080")
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}
