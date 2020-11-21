package calcserv

import (
	"fmt"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/api"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/errors"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/store"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/store/sqlstore"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	_ "github.com/mattn/go-sqlite3"
	"net/http"
)

func NewServer(config *server.Config) (*Server, error) {
	db, err := server.NewDB(config.DBPath)
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
	store store.Interface
}

func (s *Server) Store() store.Interface {
	return s.store
}

func (s *Server) registerCoreAPI() {
	apiGroup := s.Group("/api")
	apiGroup.POST("solve_expression", api.SolveExpression(s))
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
