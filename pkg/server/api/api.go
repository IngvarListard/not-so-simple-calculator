package api

import (
	"fmt"
	"github.com/IngvarListard/not-so-simple-calculator/pkg/calc"
	apperr "github.com/IngvarListard/not-so-simple-calculator/pkg/server/errors"
	"github.com/gin-gonic/gin"
	"net/http"
)

func SolveExpression() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		jsonBody := struct {
			Expression string `json:"expression" form:"expression" binding:"required"`
		}{}
		err := ctx.BindJSON(&jsonBody)
		if err != nil {
			_ = ctx.Error(apperr.New(http.StatusBadRequest, fmt.Sprintf("binding query parameters error: %s", err.Error())))
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
	}
}
