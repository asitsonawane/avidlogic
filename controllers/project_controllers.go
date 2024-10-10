package controllers

import (
	"avidlogic/database"
	"avidlogic/models" // Now using models.UserProject
	"context"
	"time"

	"github.com/gin-gonic/gin"
)

// Input struct for adding a project
type AddProjectInput struct {
	ProjectType string `json:"project_type" binding:"required"` // 'personal' or 'org'
	Username    string `json:"username" binding:"required"`
	PAT         string `json:"pat" binding:"required"`
	RepoNames   string `json:"repo_names" binding:"required"` // Comma-separated repo names
}

// AddProject adds a new project (GitHub repositories) to the user
// @Summary Add a new project
// @Description Adds a new project to the user's account (personal or organizational repos)
// @Tags Projects
// @Accept  json
// @Produce  json
// @Param project body AddProjectInput true "Project Details"
// @Success 200 {object} SuccessResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Security BearerAuth
// @Router /projects [post]
func AddProject(c *gin.Context) {
	var input AddProjectInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(400, ErrorResponse{Error: "Invalid input"})
		return
	}

	userID, _ := c.Get("userID")

	// Create the new project using the models.UserProject struct
	newProject := models.UserProject{
		UserID:      userID.(string),
		ProjectType: input.ProjectType,
		Username:    input.Username,
		PAT:         input.PAT,
		RepoNames:   input.RepoNames,
		CreatedAt:   time.Now(),
	}

	query := `INSERT INTO user_projects (user_id, project_type, username, pat, repo_names, created_at) 
              VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := database.DB.Exec(context.Background(), query, newProject.UserID, newProject.ProjectType, newProject.Username, newProject.PAT, newProject.RepoNames, newProject.CreatedAt)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: "Failed to add project"})
		return
	}

	c.JSON(200, SuccessResponse{Message: "Project added successfully"})
}
