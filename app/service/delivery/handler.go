package servicedelivery

import (
	"net/http"

	"github.com/Kostich31/techpark_db/app/domain"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	UseCase domain.ServiceUseCase
}

func NewHandler(useCase domain.ServiceUseCase) *Handler {
	return &Handler{UseCase: useCase}
}

func (handler *Handler) Status(ctx echo.Context) error {
	status, err := handler.UseCase.GetStatus()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusOK, status)
}

func (handler *Handler) Clear(ctx echo.Context) error {
	err := handler.UseCase.Clear()
	if err != nil {
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.NoContent(http.StatusOK)
}
