package main

import (
	"context"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/config"
	_ "gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/docs"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/internal/adapter"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/internal/handler"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/internal/repository"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/internal/usecase"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/db/dbpostgres"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/db/redis"
	"gitlab.humo.tj/AshurovI/orzu_aggreagtor.git/pkg/logger"
	"log"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"
)

// @title Название API
// @version 1.0
// @description Описание API
// @host localhost:8080
// @BasePath /
func main() {

	l := logger.New()
	c, err := config.New()
	if err != nil {
		l.Fatal(err)
	}
	rdb, err := redis.NewRedisClient(c.Redis)
	if err != nil {
		l.Error(err)
	}
	db, err := dbpostgres.New(c.Postgres)
	if err != nil {
		l.Fatal(err)
	}

	adp := adapter.New(l, c.Orzu, c.Otp)
	repos := repository.NewRepos(rdb, db)
	usc := usecase.New(l, adp, repos)
	hnd := handler.New(usc, l)

	routes := handler.InitRoutes(hnd, c.Srv.Token)

	server := http.Server{
		Handler:      routes,
		Addr:         c.Srv.Host + c.Srv.Port,
		ReadTimeout:  time.Duration(c.Srv.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(c.Srv.WriteTimeout) * time.Second,
	}

	go func() {
		if err := server.ListenAndServe(); err != nil {
			l.Fatal(err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGINT, syscall.SIGTERM)
	s := <-quit
	l.Printf("Shutdown server ... Signal: %v", s)

	ctx, cancel := context.WithDeadline(context.Background(), time.Now().Add(15*time.Second))
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Error while shutting sown the server: %v", err)
	}
}
