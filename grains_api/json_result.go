package grains_api

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
)

// Config

type Config struct {
	ServiceName string
	APIVersion  string
	TimeFormat  string
}

var cfg = Config{
	ServiceName: "API",
	APIVersion:  "1.0",
	TimeFormat:  time.RFC3339,
}

func Init(c Config) {
	if c.ServiceName != "" {
		cfg.ServiceName = c.ServiceName
	}
	if c.APIVersion != "" {
		cfg.APIVersion = c.APIVersion
	}
	if c.TimeFormat != "" {
		cfg.TimeFormat = c.TimeFormat
	}
}

// Request

// Request holds the context and response info for an API call
type Request struct {
	GinContext   *gin.Context
	ResponseType string
	StartTime    time.Time
	Metadata     map[string]interface{}
	Data         any // JSON-ready data (map, slice, struct, etc.)
}

// NewRequest creates a new request object
func NewRequest(gc *gin.Context, responseType string) *Request {
	return &Request{
		GinContext:   gc,
		ResponseType: responseType,
		StartTime:    time.Now(),
		Metadata:     make(map[string]interface{}),
	}
}

// Response

type Response struct {
	Service      string                 `json:"service"`
	APIVersion   string                 `json:"api_version"`
	ResponseType string                 `json:"response_type"`
	Status       string                 `json:"status"`
	Timestamp    string                 `json:"timestamp"`
	Metadata     map[string]interface{} `json:"metadata,omitempty"`
	Data         any                    `json:"data,omitempty"`
	Message      string                 `json:"message,omitempty"`
}

// Internal write method

func (req *Request) write(statusCode int, status string, message string) {
	// Calculate response time and add to metadata
	if req.Metadata == nil {
		req.Metadata = make(map[string]interface{})
	}
	req.Metadata["response_time_ms"] = time.Since(req.StartTime).Milliseconds()

	resp := Response{
		Service:      cfg.ServiceName,
		APIVersion:   cfg.APIVersion,
		ResponseType: req.ResponseType,
		Status:       status,
		Timestamp:    time.Now().UTC().Format(cfg.TimeFormat),
		Metadata:     req.Metadata,
		Data:         req.Data, // already JSON-ready
		Message:      message,
	}

	req.GinContext.JSON(statusCode, resp)
}

// Public helpers

func (r *Request) SetMeta(key string, value any) {
	if r.Metadata == nil {
		r.Metadata = make(map[string]interface{})
	}
	r.Metadata[key] = value
}

// Success sends a 200 OK response
func (req *Request) Success(statusCode int, data any, message string) {
	req.Data = data
	req.write(statusCode, "ok", message)
}

// SuccessNoData sends a 200 OK response with no data
func (req *Request) SuccessNoData(statusCode int, message string) {
	req.Data = nil
	req.write(statusCode, "ok", message)
}

func (req *Request) NotFound(message string) {
	req.write(http.StatusNotFound, "not found", message)
}

func (req *Request) NoContent(message string) {
	req.write(http.StatusNoContent, "no content", message)
}

// Error sends an error response with given status code
func (req *Request) Error(statusCode int, message string) {
	req.Data = nil
	req.write(statusCode, "error", message)
}

func (req *Request) InternalServerError() {
	req.Error(http.StatusInternalServerError, "internal server error")
}

func (req *Request) DatabaseError() {
	req.Error(http.StatusInternalServerError, "database error")
}

func (req *Request) PayloadError() {
	req.Error(http.StatusBadRequest, "payload error")
}
