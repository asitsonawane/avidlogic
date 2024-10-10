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

	// Step 1: Validate the PAT
	validPat, err := ValidatePAT(input.PAT)
	if err != nil || !validPat {
		c.JSON(400, ErrorResponse{Error: "Invalid GitHub PAT"})
		return
	}

	// Step 2: Validate the GitHub user or organization
	if input.ProjectType == "personal" {
		// Validate the GitHub user
		validUser, err := ValidateUser(input.Username)
		if err != nil || !validUser {
			c.JSON(400, ErrorResponse{Error: "GitHub user not found"})
			return
		}
	} else if input.ProjectType == "org" {
		// Validate the organization
		validOrg, err := ValidateOrg(input.PAT, input.Username)
		if err != nil || !validOrg {
			c.JSON(400, ErrorResponse{Error: "GitHub organization not found or no access"})
			return
		}
	}

	// Step 3: Validate PAT for each repository
	repoNames := strings.Split(input.RepoNames, ",")
	for _, repo := range repoNames {
		repo = strings.TrimSpace(repo)
		validRepo, err := ValidateRepoAccess(input.PAT, input.Username, repo)
		if err != nil || !validRepo {
			c.JSON(400, ErrorResponse{Error: "No access to repository: " + repo})
			return
		}
	}

	// Step 4: Add the project to the database if all checks pass
	userID, _ := c.Get("userID")
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
	_, err = database.DB.Exec(context.Background(), query, newProject.UserID, newProject.ProjectType, newProject.Username, newProject.PAT, newProject.RepoNames, newProject.CreatedAt)
	if err != nil {
		c.JSON(500, ErrorResponse{Error: "Failed to add project"})
		return
	}

	c.JSON(200, SuccessResponse{Message: "Project added successfully"})
}

// ValidatePAT checks if the provided PAT is valid by calling the /user GitHub API
func ValidatePAT(pat string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/user", nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "token "+pat)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false, err
	}

	defer resp.Body.Close()
	return true, nil
}

// ValidateUser checks if the GitHub user exists
func ValidateUser(username string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/users/"+username, nil)
	if err != nil {
		return false, err
	}

	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false, err
	}

	defer resp.Body.Close()
	return true, nil
}

// ValidateOrg checks if the GitHub organization exists
func ValidateOrg(pat, org string) (bool, error) {
	client := &http.Client{}
	req, err := http.NewRequest("GET", "https://api.github.com/orgs/"+org, nil)
	if err != nil {
		return false, err
	}

	req.Header.Set("Authorization", "token "+pat)
	resp, err := client.Do(req)
	if err != nil || resp.StatusCode != 200 {
		return false, err
	}

	defer resp.Body.Close()
	return true, nil
}

// ValidateRepoAccess checks if the PAT can access the given repository
func ValidateRepoAccess(pat, owner, repo string) (bool, error) {
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
