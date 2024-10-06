package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
)

// Project struct represents a project with its fields.
type Project struct {
	ID          int    `json:"id"`
	Title       string `json:"title"`
	Description string `json:"description"`
	Status      string `json:"status"` // "notstarted", "inprogress", "completed"
}

var projectList []Project //variable to hold project lists
var nextProjectID int = 1 //counter to create IDs

// CREATE - Create a new project
func createProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var newProject Project
	json.NewDecoder(r.Body).Decode(&newProject)
	newProject.ID = nextProjectID
	nextProjectID++
	newProject.Status = "notstarted"              // Default status
	projectList = append(projectList, newProject) //add project

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(newProject)
}

// READ - Get all projects
func getProjects(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(projectList)
}

// READ - Get a project by ID
func getProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	for _, project := range projectList {
		if project.ID == id {
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(project)
			return
		}
	}
	http.Error(w, "Project not found", http.StatusNotFound)
}

// UPDATE - Update an existing project by ID
func updateProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPut {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	for i, project := range projectList {
		if project.ID == id {
			if err := json.NewDecoder(r.Body).Decode(&projectList[i]); err != nil {
				http.Error(w, "Failed to decode request body", http.StatusBadRequest)
				return
			}
			projectList[i].ID = id // Retain original ID
			w.Header().Set("Content-Type", "application/json")
			json.NewEncoder(w).Encode(projectList[i])
			return
		}
	}
	http.Error(w, "Project not found", http.StatusNotFound)
}

// DELETE - Delete a project by ID
func deleteProject(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodDelete {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	id, err := extractID(r.URL.Path)
	if err != nil {
		http.Error(w, "Invalid project ID", http.StatusBadRequest)
		return
	}

	for i, project := range projectList {
		if project.ID == id {
			projectList = append(projectList[:i], projectList[i+1:]...)
			w.WriteHeader(http.StatusNoContent)
			return
		}
	}
	http.Error(w, "Project not found", http.StatusNotFound)
}

// function to extract the project ID from the URL
func extractID(path string) (int, error) {
	parts := strings.Split(path, "/")
	if len(parts) < 3 {
		return 0, fmt.Errorf("invalid path")
	}
	return strconv.Atoi(parts[2])
}

func main() {
	// Route handling for project API
	http.HandleFunc("/projects", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProjects(w, r)
		case http.MethodPost:
			createProject(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	http.HandleFunc("/projects/", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			getProject(w, r)
		case http.MethodPut:
			updateProject(w, r)
		case http.MethodDelete:
			deleteProject(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	// Start the HTTP server
	fmt.Println("Project Tracking API is running on port 8080...")
	log.Fatal(http.ListenAndServe(":8080", nil))
}
