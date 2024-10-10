package controllers

import (
	"avidlogic/database"
	"avidlogic/models"
	"context"
	"net/http"
	"strings"
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

	// Split the repo names (comma-separated) into a slice
	repoNames := strings.Split(input.RepoNames, ",")

	// Validate PAT for each repository
	for _, repo := range repoNames {
		repo = strings.TrimSpace(repo)
		valid, err := ValidatePAT(input.PAT, input.Username, repo)
		if err != nil || !valid {
			c.JSON(400, ErrorResponse{Error: "Invalid PAT or no access to repository: " + repo})
			return
		}
	}

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

// GitHubRepoResponse is the response structure from GitHub API
type GitHubRepoResponse struct {
	FullName string `json:"full_name"`
}

// ValidatePAT checks if the PAT can access the given repository
func ValidatePAT(pat, owner, repo string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/repos/"+owner+"/"+repo, nil)
	if err != nil {
		return false, err
	}

	// Add the Authorization header with the PAT
	req.Header.Set("Authorization", "token "+pat)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false, err
	}

	defer resp.Body.Close()
	return true, nil
}
