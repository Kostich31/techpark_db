package main

import (
	"fmt"
	"log"

	forumHandler "github.com/Kostich31/techpark_db/app/forum/delivery"
	forumRepository "github.com/Kostich31/techpark_db/app/forum/repository"
	forumUC "github.com/Kostich31/techpark_db/app/forum/usecase"
	serviceHandler "github.com/Kostich31/techpark_db/app/service/delivery"
	serviceRepository "github.com/Kostich31/techpark_db/app/service/repository"
	serviceUC "github.com/Kostich31/techpark_db/app/service/usecase"
	threadHandler "github.com/Kostich31/techpark_db/app/thread/delivery"
	threadRepository "github.com/Kostich31/techpark_db/app/thread/repository"
	threadUC "github.com/Kostich31/techpark_db/app/thread/usecase"
	"github.com/Kostich31/techpark_db/app/tools"
	userHandler "github.com/Kostich31/techpark_db/app/user/delivery"
	userRepository "github.com/Kostich31/techpark_db/app/user/repository"
	userUC "github.com/Kostich31/techpark_db/app/user/usecase"
	validator "github.com/go-playground/validator"
	"github.com/jackc/pgx"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/labstack/echo/v4"
)

var (
	router = echo.New()
)

func main() {
	db, err := GetPostgres()
	if err != nil {
		log.Fatal(err)
	}

	userHandler := userHandler.NewHandler(userUC.NewUseCase(
		userRepository.NewRepository(db)))
	forumHandler := forumHandler.NewHandler(forumUC.NewUseCase(
		forumRepository.NewRepository(db), threadRepository.NewRepository(db)))
	threadHandler := threadHandler.NewHandler(threadUC.NewUseCase(
		threadRepository.NewRepository(db), userRepository.NewRepository(db), forumRepository.NewRepository(db)))
	serviceHandler := serviceHandler.NewHandler(serviceUC.NewUseCase(serviceRepository.NewRepository(db)))

	validator := validator.New()
	router.Validator = tools.NewCustomValidator(validator)

	router.POST("api/user/:nickname/create", userHandler.SignUpUser)
	router.GET("api/user/:nickname/profile", userHandler.GetUser)
	router.POST("api/user/:nickname/profile", userHandler.UpdateUser)
	router.POST("api/forum/create", forumHandler.CreateForum)
	router.GET("api/forum/:slug/details", forumHandler.GetForumDetails)
	router.POST("api/forum/:slug/create", forumHandler.CreateThread)
	router.GET("api/forum/:slug/users", forumHandler.GetUsersForum)
	router.GET("api/forum/:slug/threads", forumHandler.GetForumThreads)
	router.POST("api/thread/:slug_or_id/create", threadHandler.CreatePosts)
	router.POST("api/thread/:slug_or_id/vote", threadHandler.Vote)
	router.GET("api/thread/:slug_or_id/details", threadHandler.Details)
	router.GET("api/thread/:slug_or_id/posts", threadHandler.GetPosts)
	router.POST("api/thread/:slug_or_id/details", threadHandler.UpdateThread)
	router.GET("api/post/:id/details", threadHandler.GetOnePost)
	router.POST("api/post/:id/details", threadHandler.UpdatePost)
	router.GET("api/service/status", serviceHandler.Status)
	router.POST("api/service/clear", serviceHandler.Clear)
	if err := router.Start(":5000"); err != nil {
		log.Fatal(err)
	}
}

func GetPostgres() (*pgx.ConnPool, error) {
	dsn := fmt.Sprintf("user=%s dbname=%s password=%s host=%s port=%s sslmode=disable",
		"kostya", "kostya",
		"kostya", "127.0.0.1",
		"5432")
	db, err := pgx.ParseConnectionString(dsn)
	if err != nil {
		return nil, err
	}
	pool, err := pgx.NewConnPool(pgx.ConnPoolConfig{
		ConnConfig:     db,
		MaxConnections: 100,
		AfterConnect:   nil,
		AcquireTimeout: 0,
	})
	if err != nil {
		log.Fatalf("Error %s occurred during connection to database", err)
	}

	return pool, nil
}
