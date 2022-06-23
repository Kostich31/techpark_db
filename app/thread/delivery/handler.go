package threadhandler

import (
	"net/http"

	"github.com/Kostich31/techpark_db/app/domain"
	"github.com/Kostich31/techpark_db/app/tools"
	"github.com/labstack/echo/v4"
)

type Handler struct {
	UseCase domain.ThreadUseCase
}

func NewHandler(usecase domain.ThreadUseCase) *Handler {
	return &Handler{UseCase: usecase}
}

func (handler *Handler) CreatePosts(ctx echo.Context) error {
	var newPosts []domain.Post

	if err := ctx.Bind(&newPosts); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	slugOrId := ctx.Param("slug_or_id")
	posts, err := handler.UseCase.CreatePosts(slugOrId, newPosts)
	if err != nil {
		if err.Message == domain.NoUser {
			return ctx.JSON(http.StatusNotFound, err)
		}
		if err.Message == domain.ConflictData {
			return ctx.JSON(http.StatusConflict, err)
		}
		if err.Message == domain.BadParentPost {
			return ctx.JSON(http.StatusConflict, err)
		}
		return ctx.JSON(http.StatusInternalServerError, err)
	}

	return ctx.JSON(http.StatusCreated, posts)
}

func (handler *Handler) Vote(ctx echo.Context) error {
	var newVoice domain.Vote

	if err := ctx.Bind(&newVoice); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	slugOrId := ctx.Param("slug_or_id")

	thread, err := handler.UseCase.CreateVote(slugOrId, newVoice)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, thread)
}

func (handler *Handler) Details(ctx echo.Context) error {
	slugOrId := ctx.Param("slug_or_id")
	thread, err := handler.UseCase.GetThreadDetails(slugOrId)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, thread)
}

func (handler *Handler) GetPosts(ctx echo.Context) error {
	filter := tools.ParseQueryFilterPost(ctx)
	slugOrId := ctx.Param("slug_or_id")

	posts, err := handler.UseCase.GetPosts(slugOrId, filter)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, posts)
}

func (handler *Handler) UpdateThread(ctx echo.Context) error {
	var newThread domain.Thread

	if err := ctx.Bind(&newThread); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}
	slugOrId := ctx.Param("slug_or_id")

	thread, err := handler.UseCase.UpdateThread(slugOrId, newThread)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, thread)
}

func (handler *Handler) GetOnePost(ctx echo.Context) error {
	id := ctx.Param("id")
	filter := tools.ParseQueryFilterOnePost(ctx)

	post, err := handler.UseCase.GetPost(id, filter)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, post)
}

func (handler *Handler) UpdatePost(ctx echo.Context) error {
	var postInfo domain.Post

	if err := ctx.Bind(&postInfo); err != nil {
		return ctx.JSON(http.StatusBadRequest, err.Error())
	}

	id := ctx.Param("id")
	post, err := handler.UseCase.UpdatePost(id, postInfo)
	if err != nil {
		return ctx.JSON(http.StatusNotFound, err)
	}

	return ctx.JSON(http.StatusOK, post)
}
