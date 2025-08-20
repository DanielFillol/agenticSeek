package main

import (
	"embed"
	"fmt"
	"io/fs"
	"log"
	"net/http"
	"path/filepath"
	"strings"

	"flexible-agent/internal/agent"
	"flexible-agent/internal/documents"
	"flexible-agent/internal/llm"
	"flexible-agent/internal/schemas"

	"github.com/gin-gonic/gin"
)

//go:embed frontend/build
var frontendFS embed.FS

type Server struct {
	agent *agent.Agent
}

func main() {
	provider := "ollama"
	model := "llama3"

	llmClient, err := llm.NewLLM(provider, model)
	if err != nil {
		log.Fatalf("failed to create LLM client: %v", err)
	}

	server := &Server{
		agent: agent.NewAgent(llmClient),
	}

	r := gin.Default()

	// API routes
	api := r.Group("/api")
	{
		api.GET("/health", server.healthHandler)
		api.POST("/chat", server.chatHandler)
		api.POST("/upload", server.uploadHandler)
	}

	// Serve frontend
	staticFS, err := fs.Sub(frontendFS, "frontend/build")
	if err != nil {
		log.Fatalf("failed to create static file system: %v", err)
	}
	r.StaticFS("/", http.FS(staticFS))

	r.NoRoute(func(c *gin.Context) {
		if !strings.HasPrefix(c.Request.RequestURI, "/api") {
			c.FileFromFS("index.html", http.FS(staticFS))
		}
	})

	r.Run() // listen and serve on 0.0.0.0:8080
}

func (s *Server) healthHandler(c *gin.Context) {
	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}

func (s *Server) chatHandler(c *gin.Context) {
	var req schemas.ChatRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	plan, err := s.agent.Plan(req.Message)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create plan"})
		return
	}

	err = s.agent.Execute(plan)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to execute plan"})
		return
	}

	resp := schemas.ChatResponse{
		Message: fmt.Sprintf("Plan:\n%s", strings.Join(plan, "\n")),
	}
	c.JSON(http.StatusOK, resp)
}

func (s *Server) uploadHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "file upload error"})
		return
	}

	ext := strings.ToLower(filepath.Ext(file.Filename))
	f, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to open uploaded file"})
		return
	}
	defer f.Close()

	var text string
	switch ext {
	case ".pdf":
		text, err = documents.ExtractTextFromPDF(f)
	case ".docx":
		text, err = documents.ExtractTextFromDOCX(f)
	default:
		c.JSON(http.StatusBadRequest, gin.H{"error": "unsupported file type"})
		return
	}

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to extract text from file"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"text": text})
}
