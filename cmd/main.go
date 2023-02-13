package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"todo"
	"todo/pkg/handler"
	"todo/pkg/repository"
	"todo/pkg/service"
)

func main() {
	logrus.SetFormatter(&logrus.JSONFormatter{})
	if err := initConfig(); err != nil {
		logrus.Fatalf("error ehile reading config: %v", err)
	}

	if err := godotenv.Load(); err != nil {
		logrus.Fatalf("failed to load env file: %v", err)
	}

	db, err := repository.NewPostgresDB(repository.Config{
		Host:     viper.GetString("db.host"),
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"),
		DBName:   viper.GetString("db.dbName"),
		SSLMode:  viper.GetString("db.sslMode"),
	})

	if err != nil {
		logrus.Fatalf("error connecting db: %v", err)
	}

	repos := repository.NewRepository(db)
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	srv := todo.Server{}
	go func() {
		err = srv.Run(viper.GetString("port"), handlers.InitRouts())
		if err != nil {
			logrus.Errorf("error: %v", err)
		}
	}()

	quit := make(chan os.Signal)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Infof("shutting Down")

	err = srv.Shutdown(context.Background())
	if err != nil {
		logrus.Errorf("failed shutdown server, error: %v", err)
	}

	err = db.Close()
	if err != nil {
		logrus.Errorf("failed closed DB connection, error: %v", err)
	}
}

func initConfig() error {
	viper.AddConfigPath("config")
	viper.SetConfigName("config")
	return viper.ReadInConfig()
}
