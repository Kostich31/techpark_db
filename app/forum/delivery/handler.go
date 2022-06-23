package delivery

import (
	"net/http"

	"github.com/Kostich31/techpark_db/app/domain"
	"github.com/Kostich31/techpark_db/app/tools"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	useCase domain.ForumUseCase
}

func NewHandler(useCase domain.ForumUseCase) *Handler {
	return &Handler{useCase: useCase}
}

func (handler *Handler) CreateForum(ctx echo.Context) error {
	var newForum domain.Forum

	if err := ctx.Bind(&newForum); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(&newForum); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	forum, err := handler.useCase.CreateForum(newForum)
	if err != nil {
		if err.Message == domain.NoUser {
			return ctx.JSON(http.StatusNotFound, err)
		}
		if err.Message == domain.ConflictData {
			return ctx.JSON(http.StatusConflict, forum)
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, forum)
}

func (handler *Handler) GetForumDetails(ctx echo.Context) error {
	slug := ctx.Param("slug")

	forum, err := handler.useCase.GetDetailsForum(slug)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, forum)
}

func (handler *Handler) CreateThread(ctx echo.Context) error {
	var newThread domain.Thread

	if err := ctx.Bind(&newThread); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	if err := ctx.Validate(&newThread); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	newThread.Forum = ctx.Param("slug")

	thread, err := handler.useCase.CreateThread(newThread)
	if err != nil {
		if err.Message == domain.NoUser {
			return ctx.JSON(http.StatusNotFound, err)
		}
		if err.Message == domain.ConflictData {
			return ctx.JSON(http.StatusConflict, thread)
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, thread)
}

func (handler *Handler) GetUsersForum(ctx echo.Context) error {
	slug := ctx.Param("slug")
	filter := tools.ParseQueryFilterUser(ctx)

	users, err := handler.useCase.GetUsersForum(slug, filter)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, users)
}

func (handler *Handler) GetForumThreads(ctx echo.Context) error {
	slug := ctx.Param("slug")
	filter := tools.ParseQueryFilterThread(ctx)

	users, err := handler.useCase.GetForumThreads(slug, filter)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, users)
}
