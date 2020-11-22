package calcserv

import (
	"fmt"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/api"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/errors"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/logging"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/store"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/store/sqlstore"
	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func New(config *calcapp.Config) (*Server, error) {
	logger, err := logging.New(config)
	if err != nil {
		return nil, fmt.Errorf("creating logger error: %w", err)
	}

	db, err := calcapp.NewDB(config.DBPath)
	if err != nil {
		return nil, fmt.Errorf("can't establish database connection: %w", err)
	}

	if config.LogsToFile != nil {
		gin.DefaultWriter = logger.Writer()
	}

	if !config.AccessLog {
		f, err := os.OpenFile(os.DevNull, os.O_WRONLY|os.O_APPEND, 0660)
		if err != nil {
			return nil, fmt.Errorf("error opening %s: %w", os.DevNull, err)
		}

		gin.DefaultWriter = f
	}

	s := &Server{Engine: gin.Default(), store: sqlstore.New(db), logger: logger}
	s.Use(ErrorHandler(logger))
	s.registerCoreAPI()

	return s, nil
}

type Server struct {
	*gin.Engine
	store  store.Interface
	logger *logrus.Logger
}

func (s *Server) Store() store.Interface {
	return s.store
}

func (s *Server) Logger() *logrus.Logger {
	return s.logger
}

func (s *Server) registerCoreAPI() {
	apiGroup := s.Group("/api")
	{
		apiGroup.POST("solve_expression", api.SolveExpression(s))
		apiGroup.GET("get_all_history", api.GetAllHistory(s.store))
		apiGroup.POST("get_history_by_time_range", api.GetHistoryByTimeRange(s.store))
	}
}

func ErrorHandler(logger *logrus.Logger) gin.HandlerFunc {
	return errorHandler(gin.ErrorTypeAny, logger)
}

func errorHandler(errType gin.ErrorType, logger *logrus.Logger) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Next()
		detectedErrors := c.Errors.ByType(errType)
		if len(detectedErrors) == 0 {
			return
		}

		err := detectedErrors[len(detectedErrors)-1].Err
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
			logrus.Fatalf("Internal Server Error occurred: %v", v)
		}
		c.IndentedJSON(parsedError.Code, gin.H{"error": parsedError})
		logger.Debugf("query handling error: %v", parsedError.Error())
		c.Abort()
		return
	}
}
