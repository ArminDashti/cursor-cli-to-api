package httpserver

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/ArminDashti/cursor-cli-to-api/cursor-cli-to-api-api/internal/auth"
	"github.com/ArminDashti/cursor-cli-to-api/cursor-cli-to-api-api/internal/config"
	"github.com/ArminDashti/cursor-cli-to-api/cursor-cli-to-api-api/internal/cursorcli"
	"github.com/ArminDashti/cursor-cli-to-api/cursor-cli-to-api-api/internal/openai"
	"github.com/ArminDashti/cursor-cli-to-api/cursor-cli-to-api-api/internal/store"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	cfg    config.Config
	db     *sql.DB
	runner *cursorcli.Runner
	router *gin.Engine
}

func New(cfg config.Config, db *sql.DB) *Server {
	r := gin.New()
	r.Use(gin.Logger(), gin.Recovery())
	r.Use(cors.New(cors.Config{
		AllowOrigins:     []string{cfg.CORSOrigin, "http://127.0.0.1:5201"},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s := &Server{
		cfg: cfg,
		db:  db,
		runner: &cursorcli.Runner{
			Bin:       cfg.CursorAgentBin,
			Mode:      cfg.CursorAgentMode,
			Workspace: cfg.CursorWorkspace,
			Timeout:   time.Duration(cfg.CursorTimeoutSec) * time.Second,
		},
		router: r,
	}
	s.routes()
	return s
}

func (s *Server) Router() *gin.Engine { return s.router }

func (s *Server) routes() {
	s.router.GET("/health", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"status": "ok"})
	})

	v1 := s.router.Group("/v1")
	v1.Use(auth.Middleware(s.cfg.AuthUsername, s.cfg.AuthPassword))
	{
		v1.GET("/models", s.listModels)
		v1.POST("/chat/completions", s.chatCompletions)
		v1.POST("/responses", s.createResponse)
		v1.GET("/responses/:id", s.getResponse)
	}
}

func (s *Server) listModels(c *gin.Context) {
	models, err := s.runner.ListModels(c.Request.Context())
	if err != nil {
		models = []string{"auto"}
	}
	data := make([]gin.H, 0, len(models))
	for _, m := range models {
		data = append(data, gin.H{
			"id":       m,
			"object":   "model",
			"created":  time.Now().Unix(),
			"owned_by": "cursor",
		})
	}
	c.JSON(http.StatusOK, gin.H{"object": "list", "data": data})
}

func (s *Server) chatCompletions(c *gin.Context) {
	var req openai.ChatCompletionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"message": err.Error(),
			"type":    "invalid_request_error",
		}})
		return
	}
	if len(req.Messages) == 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"message": "messages is required",
			"type":    "invalid_request_error",
		}})
		return
	}
	model := req.Model
	if model == "" {
		model = "auto"
	}
	prompt := openai.PromptFromChatMessages(req.Messages)
	if prompt == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"message": "empty prompt",
			"type":    "invalid_request_error",
		}})
		return
	}

	if req.Stream {
		s.streamChat(c, model, prompt)
		return
	}

	res, err := s.runner.Run(c.Request.Context(), prompt, model)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"message": err.Error(),
			"type":    "api_error",
		}})
		return
	}
	c.JSON(http.StatusOK, openai.BuildChatCompletion(model, res.Text))
}

func (s *Server) streamChat(c *gin.Context, model, prompt string) {
	id := openai.NewChatID()
	c.Header("Content-Type", "text/event-stream")
	c.Header("Cache-Control", "no-cache")
	c.Header("Connection", "keep-alive")
	c.Status(http.StatusOK)
	flusher, ok := c.Writer.(http.Flusher)
	if !ok {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{"message": "streaming unsupported"}})
		return
	}

	write := func(v any) {
		b, _ := json.Marshal(v)
		fmt.Fprintf(c.Writer, "data: %s\n\n", b)
		flusher.Flush()
	}

	// Role opener like OpenAI streams
	write(map[string]any{
		"id":      id,
		"object":  "chat.completion.chunk",
		"created": time.Now().Unix(),
		"model":   model,
		"choices": []map[string]any{{
			"index": 0,
			"delta": map[string]any{"role": "assistant"},
			"finish_reason": nil,
		}},
	})

	_, err := s.runner.Stream(c.Request.Context(), prompt, model, func(delta string) error {
		write(openai.ChatStreamChunk(id, model, delta, nil))
		return nil
	})
	if err != nil {
		write(gin.H{"error": gin.H{"message": err.Error(), "type": "api_error"}})
		fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
		flusher.Flush()
		return
	}
	finish := "stop"
	write(openai.ChatStreamChunk(id, model, "", &finish))
	fmt.Fprintf(c.Writer, "data: [DONE]\n\n")
	flusher.Flush()
}

func (s *Server) createResponse(c *gin.Context) {
	var req openai.ResponsesRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"message": err.Error(),
			"type":    "invalid_request_error",
		}})
		return
	}
	prompt, err := openai.PromptFromResponsesInput(req.Input)
	if err != nil || strings.TrimSpace(prompt) == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": gin.H{
			"message": "input is required",
			"type":    "invalid_request_error",
		}})
		return
	}
	model := req.Model
	if model == "" {
		model = "auto"
	}

	res, err := s.runner.Run(c.Request.Context(), prompt, model)
	if err != nil {
		c.JSON(http.StatusBadGateway, gin.H{"error": gin.H{
			"message": err.Error(),
			"type":    "api_error",
		}})
		return
	}
	obj := openai.BuildResponse(model, res.Text)
	if err := store.SaveResponse(c.Request.Context(), s.db, obj.ID, obj.Model, obj.Status, obj); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"message": err.Error(),
			"type":    "api_error",
		}})
		return
	}
	c.JSON(http.StatusOK, obj)
}

func (s *Server) getResponse(c *gin.Context) {
	id := c.Param("id")
	raw, err := store.GetResponse(c.Request.Context(), s.db, id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": gin.H{
			"message": err.Error(),
			"type":    "api_error",
		}})
		return
	}
	if raw == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": gin.H{
			"message": "No such response",
			"type":    "invalid_request_error",
			"code":    "response_not_found",
		}})
		return
	}
	c.Data(http.StatusOK, "application/json", raw)
}
