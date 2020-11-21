package server

import (
	"fmt"
	"github.com/IngvarListard/not-so-simple-calculator/internal/server/api"
	"github.com/IngvarListard/not-so-simple-calculator/internal/server/errors"
	"github.com/IngvarListard/not-so-simple-calculator/internal/store"
	"github.com/IngvarListard/not-so-simple-calculator/internal/store/sqlstore"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func NewServer(config *Config) (*Server, error) {
	db, err := newDB(config.DBPath)
	if err != nil {
		return nil, fmt.Errorf("can't establish database connection: %w", err)
	}

	s := &Server{Engine: gin.Default(), store: sqlstore.New(db)}
	s.Use(ErrorHandler())
	s.registerCoreAPI()

	return s, nil
}

type Server struct {
	*gin.Engine
	store store.Store
}

func (s *Server) registerCoreAPI() {
	apiGroup := s.Group("/api")
	apiGroup.POST("solve_expression", api.SolveExpression())
}

func ErrorHandler() gin.HandlerFunc {
	return errorHandler(gin.ErrorTypeAny)
}

func errorHandler(errType gin.ErrorType) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)

		if len(detectedErrors) > 0 {
			err := detectedErrors[0].Err
			var parsedError *errors.AppError
			switch v := err.(type) {
			case *errors.AppError:
				parsedError = v
			case validator.ValidationErrors:
				parsedError = &errors.AppError{
					Code:    http.StatusBadRequest,
					Message: fmt.Sprintf("query validation error: %v", v.Error()),
				}
			default:
				parsedError = &errors.AppError{
					Code:    http.StatusInternalServerError,
					Message: "Internal Server Error",
				}
			}
			c.IndentedJSON(parsedError.Code, gin.H{"error": parsedError})
			c.Abort()
			return
		}
	}
}
