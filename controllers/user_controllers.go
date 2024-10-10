package controllers

import (
	"context"
	"log"
	"net/http"
	"time"

	"avidlogic/database" // Import the new database package
	"avidlogic/models"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
)

// ErrorResponse defines the structure for error messages.
type ErrorResponse struct {
	Error string `json:"error"`
}

// Input struct for creating a user
type CreateUserInput struct {
	Username string `json:"username" binding:"required"`
	Email    string `json:"email" binding:"required"`
	Password string `json:"password" binding:"required"`
}

// HashPassword hashes the password using bcrypt
func HashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), 14)
	return string(bytes), err
}

// CreateUser handles the creation of a new user
// @Summary Create a new user
// @Description Create a new user with username, email, and password
// @Tags Users
// @Accept json
// @Produce json
// @Param user body CreateUserInput true "User Data"
// @Success 200 {object} models.User
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /users [post]
func CreateUser(c *gin.Context) {
	var input CreateUserInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	// Check if user already exists
	var existingUser models.User
	err := database.DB.QueryRow(context.Background(), "SELECT id FROM users WHERE email=$1", input.Email).Scan(&existingUser.ID)
	if err == nil {
		c.JSON(http.StatusConflict, ErrorResponse{Error: "User with this email already exists"})
		return
	}

	passwordHash, err := HashPassword(input.Password)
	if err != nil {
		log.Fatal("Failed to hash password:", err)
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to hash password"})
		return
	}

	newUser := models.User{
		ID:           uuid.New(),
		Username:     input.Username,
		Email:        input.Email,
		PasswordHash: passwordHash,
		CreatedAt:    time.Now(),
		UpdatedAt:    time.Now(),
	}

	query := `INSERT INTO users (id, username, email, password_hash, created_at, updated_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`
	_, err = database.DB.Exec(context.Background(), query,
		newUser.ID, newUser.Username, newUser.Email, newUser.PasswordHash, newUser.CreatedAt, newUser.UpdatedAt)

	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: "Failed to create user"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "User created successfully", "user": newUser})
}
