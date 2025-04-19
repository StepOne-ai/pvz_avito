package routes

import (
	"fmt"
	"net/http"
	"regexp"
	"strconv"
	"time"

	"github.com/StepOne-ai/pvz_avito/internal/db"
	"github.com/StepOne-ai/pvz_avito/internal/middleware"
	"github.com/StepOne-ai/pvz_avito/internal/models"
	"github.com/gin-gonic/gin"
)

func DummyLogin(c *gin.Context) {
	var req struct {
		Role string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
		return
	}

	user := models.User{
		ID:       "123",
		Email:    "dummy@mail.ru",
		Role:     "dummy",
		Password: "dummy_password",
	}

	tokenString, err := middleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: "Failed to generate token"})
		return
	}
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	c.SetCookie("role", user.Role, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func Register(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
		Role     string `json:"role"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
		return
	}

	if !isValidEmail(req.Email) {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid email"})
		return
	}

	id := fmt.Sprintf("user-%d", time.Now().Unix())
	err := db.CreateUser(id, req.Email, req.Password, req.Role)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, models.User{ID: id, Email: req.Email, Role: req.Role, Password: req.Password})
}

func isValidEmail(email string) bool {
	emailRegex := regexp.MustCompile(`^[a-zA-Z0-9._%+\-]+@[a-zA-Z0-9.\-]+\.[a-zA-Z]{2,}$`)
	return emailRegex.MatchString(email)
}

func Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
		return
	}

	if db.CheckCredentials(req.Email, req.Password) != nil {
		c.JSON(http.StatusUnauthorized, models.Error{Message: "Invalid credentials"})
		return
	}
	user, err := db.GetUserByEmail(req.Email)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "User does not exist"})
		return
	}

	tokenString, err := middleware.GenerateToken(user)
	if err != nil {
		c.JSON(http.StatusInternalServerError, models.Error{Message: "Failed to generate token"})
		return
	}
	c.SetCookie("token", tokenString, 3600, "/", "localhost", false, true)
	c.SetCookie("role", user.Role, 3600, "/", "localhost", false, true)

	c.JSON(http.StatusOK, gin.H{"token": tokenString})
}

func PVZ_post(c *gin.Context) {
	var req struct {
		City string `json:"city"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
		return
	}
	// City check
	if !(req.City == "Москва" || req.City == "Казань" || req.City == "Санкт-Петербург") {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid city"})
		return
	}

	var pvz models.PVZ
	pvz.ID = fmt.Sprintf("pvz-%d", time.Now().Unix())
	pvz.RegistrationDate = time.Now()
	pvz.City = req.City

	err := db.CreatePVZ(pvz.ID, pvz.City, pvz.RegistrationDate.Format(time.RFC3339))
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, req)
}

func PVZ_get(c *gin.Context) {
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
}

func PVZ_close_last_reception(c *gin.Context) {
	pvzId := c.Param("pvzId")
	err := db.CloseLastReception(pvzId)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Reception closed"})
}

func PVZ_delete_last_product(c *gin.Context) {
	pvzId := c.Param("pvzId")
	err := db.DeleteLastProduct(pvzId)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}

func Receptions(c *gin.Context) {
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
}

func Products(c *gin.Context) {
	var req struct {
		Type        string `json:"type"`
		PvzId       string `json:"pvzId"`
		ReceptionId string `json:"receptionId"`
	}
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: "Invalid request"})
		return
	}
	product := models.Product{
		ID:          fmt.Sprintf("product-%d", time.Now().Unix()),
		DateTime:    time.Now(),
		Type:        req.Type,
		ReceptionId: req.ReceptionId,
	}

	err := db.CreateProduct(product.ID, product.DateTime.Format(time.RFC3339), product.Type, product.ReceptionId)
	if err != nil {
		c.JSON(http.StatusBadRequest, models.Error{Message: err.Error()})
		return
	}
	c.JSON(http.StatusCreated, product)
}
