package api

import (
	"fmt"
	apperr "github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/errors"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/models"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/server"
	"github.com/IngvarListard/not-so-simple-calculator/internal/calcapp/store"
	"github.com/IngvarListard/not-so-simple-calculator/pkg/calc"
	"github.com/gin-gonic/gin"
	"net/http"
	"time"
)

func SolveExpression(s server.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jsonBody := struct {
			Expression string `json:"expression" form:"expression" binding:"required"`
		}{}
		err := ctx.BindJSON(&jsonBody)
		if err != nil {
			_ = ctx.Error(apperr.New(http.StatusBadRequest, fmt.Sprintf("query validation error: %s", err.Error())))
			return
		}

		p, err := calc.NewParser(jsonBody.Expression)
		if err != nil {
			_ = ctx.Error(apperr.New(http.StatusBadRequest, fmt.Sprintf("expression parsing error: %s", err.Error())))
			return
		}

		solver, err := calc.NewSolver(p)
		if err != nil {
			_ = ctx.Error(apperr.New(http.StatusBadRequest, fmt.Sprintf("expression parsing error: %s", err.Error())))
			return
		}

		result, err := solver.Solve()
		if err != nil {
			_ = ctx.Error(apperr.New(http.StatusBadRequest, fmt.Sprintf("expression parsing error: %s", err.Error())))
			return
		}

		ctx.JSON(http.StatusOK, gin.H{"result": result})

		err = s.Store().History().Create(&models.History{
			EventTime:  time.Now(),
			Expression: jsonBody.Expression,
			Result:     result,
		})
		if err != nil {
			s.Logger().Errorf("History creation error: %v", err)
		}
	}
}

func GetHistoryByTimeRange(s store.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		b := struct {
			StartTime string `json:"start_time" form:"start_time" binding:"required"`
			EndTime   string `json:"end_time" form:"end_time" binding:"required"`
		}{}
		err := ctx.BindJSON(&b)
		if err != nil {
			_ = ctx.Error(apperr.New(http.StatusBadRequest, fmt.Sprintf("query validation error: %s", err.Error())))
			return
		}

		history, err := s.History().GetHistoryByTimeRange(b.StartTime, b.EndTime)
		if err != nil {
			_ = ctx.Error(apperr.New(http.StatusBadRequest, fmt.Sprintf("error getting history: %s", err.Error())))
			return
		}
		ctx.JSON(http.StatusOK, history)
	}
}

func GetAllHistory(s store.Interface) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		history, err := s.History().GetAllHistory()
		if err != nil {
			_ = ctx.Error(apperr.New(http.StatusBadRequest, fmt.Sprintf("error getting history records: %s", err.Error())))
			return
		}

		ctx.JSON(http.StatusOK, history)
	}
}
