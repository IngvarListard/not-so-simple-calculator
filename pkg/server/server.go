package server

import (
	"fmt"
	"github.com/IngvarListard/not-so-simple-calculator/pkg/server/api"
	"github.com/IngvarListard/not-so-simple-calculator/pkg/server/errors"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"net/http"
)

func NewServer() Server {
	s := Server{gin.Default()}
	s.Use(ErrorHandler())
	s.registerCoreAPI()
	return s
}

type Server struct {
	*gin.Engine
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
					Message: fmt.Sprintf("query params validation error: %v", v.Error()),
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
