package main

import (
	"os"
	"os/signal"
	"syscall"
	"todo-app/pkg/handler"
	"todo-app/pkg/repository"
	"todo-app/pkg/service"

	"github.com/gin-gonic/gin"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
)

// Использование swagger
// http://localhost:8000/swagger/index.html#/

// @title Todo App API
// @version 1.0
// @description API Server for TodoList Application

// @host localhost:8000
// @BasePath /

// @securityDefinitions.apikey ApiKeyAuth
// @in header
// @name Authorization

func main() {
	logrus.SetFormatter(new(logrus.JSONFormatter)) // Установка логирования в формат JSON

	context := &gin.Context{} // Контекст

	if err := initConfig(); err != nil { //Инициализируем конфигурации
		logrus.Fatalf("error initializing configs: %s", err.Error())
		return
	}

	if err := godotenv.Load(); err != nil { //Загрузка переменного окружения (для передачи пароля из файла .env)
		logrus.Fatalf("error loading env variables: %s", err.Error())
		return
	}

	db, err := repository.NewPostgresDB(repository.Config{ //Инициализация БД
		Host:     viper.GetString("db.host"), // Читаем данные из файла config.yml по ключу
		Port:     viper.GetString("db.port"),
		Username: viper.GetString("db.username"),
		Password: os.Getenv("DB_PASSWORD"), // Читаем пароль из файла .env по ключу
		DBName:   viper.GetString("db.dbname"),
		SSLMode:  viper.GetString("db.sslmode"),
	})
	if err != nil {
		logrus.Fatalf("failed to initialize db: %s", err.Error())
		return
	}

	redisClient, err := repository.NewRedisCache(context, repository.ConfigRedis{ // Подключение к серверу Redis
		Addr:     viper.GetString("redis.addr"),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
	})

	if err != nil {
		logrus.Fatalf("failed to initialize Redis: %s", err.Error())
		return
	}

	repos := repository.NewRepository(db, context, redisClient) // Создание зависимостей
	services := service.NewService(repos)
	handlers := handler.NewHandler(services)

	rsv := handlers.InitRoutes()
	go func() {
		if err := rsv.Run(viper.GetString("port")); err != nil {
			logrus.Fatalf("Error run web serv")
			return
		}
	}()

	logrus.Print("TodoApp Started")

	quit := make(chan os.Signal, 1)
	signal.Notify(quit, syscall.SIGTERM, syscall.SIGINT)
	<-quit

	logrus.Print("TodoApp Stoped")

	if err := db.Close(); err != nil {
		logrus.Errorf("error occured on db connection close: %s", err.Error())
	}
}

func initConfig() error { //Инициализация конфигураций
	//viper.AddConfigPath("configs")
	//viper.SetConfigName("config")
	viper.SetConfigFile("config.yml")
	return viper.ReadInConfig()
}
