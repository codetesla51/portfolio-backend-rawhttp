package handlers

import (
	"encoding/json"
	"log"

	"backend-raw-http/database"
	"backend-raw-http/models"

	"github.com/codetesla51/raw-http/server"
)

type ProjectHandler struct {
	repo *models.ProjectRepository
}

func NewProjectHandler() *ProjectHandler {
	return &ProjectHandler{
		repo: models.NewProjectRepository(database.DB),
	}
}

type APIResponse struct {
	Success bool        `json:"success"`
	Data    interface{} `json:"data,omitempty"`
	Message string      `json:"message,omitempty"`
	Error   string      `json:"error,omitempty"`
}

// GetAllProjects handles GET /api/projects
func (h *ProjectHandler) GetAllProjects(req *server.Request) ([]byte, string) {
	projects, err := h.repo.GetAll()
	if err != nil {
		log.Printf("Error fetching projects: %v", err)
		return jsonResponse(500, APIResponse{
			Success: false,
			Error:   "Failed to fetch projects",
		})
	}

	if projects == nil {
		projects = []models.Project{}
	}

	return jsonResponse(200, APIResponse{
		Success: true,
		Data:    projects,
	})
}

// GetProjectBySlug handles GET /api/projects/:slug
func (h *ProjectHandler) GetProjectBySlug(req *server.Request) ([]byte, string) {
	slug := req.PathParams["slug"]
	if slug == "" {
		return jsonResponse(400, APIResponse{
			Success: false,
			Error:   "Slug is required",
		})
	}

	project, err := h.repo.GetBySlug(slug)
	if err != nil {
		log.Printf("Error fetching project by slug: %v", err)
		return jsonResponse(500, APIResponse{
			Success: false,
			Error:   "Failed to fetch project",
		})
	}

	if project == nil {
		return jsonResponse(404, APIResponse{
			Success: false,
			Error:   "Project not found",
		})
	}

	return jsonResponse(200, APIResponse{
		Success: true,
		Data:    project,
	})
}

func jsonResponse(statusCode int, data interface{}) ([]byte, string) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return server.Serve500("Internal server error")
	}

	statusText := getStatusText(statusCode)
	return server.CreateResponseBytes(
		statusCodeToString(statusCode),
		"application/json",
		statusText,
		jsonData,
	)
}

func statusCodeToString(code int) string {
	switch code {
	case 200:
		return "200"
	case 201:
		return "201"
	case 400:
		return "400"
	case 404:
		return "404"
	case 500:
		return "500"
	default:
		return "200"
	}
}

func getStatusText(code int) string {
	switch code {
	case 200:
		return "OK"
	case 201:
		return "Created"
	case 400:
		return "Bad Request"
	case 404:
		return "Not Found"
	case 500:
		return "Internal Server Error"
	default:
		return "OK"
	}
}
