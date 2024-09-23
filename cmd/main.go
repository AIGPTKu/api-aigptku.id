package main

import (
	"net/http"

	"github.com/AIGPTku/api-aigptku.id/app/controller/rest"
	"github.com/AIGPTku/api-aigptku.id/lib/utils/database"
	"github.com/AIGPTku/api-aigptku.id/lib/utils/middleware"
	"github.com/goccy/go-json"
	"github.com/gofiber/fiber/v2"
	"github.com/gofiber/fiber/v2/middleware/cors"
	"github.com/spf13/viper"
)

func loadEnv() {
	viper.SetConfigFile(".env")
	viper.SetConfigType("env")
	viper.AutomaticEnv()
	err := viper.ReadInConfig()
	if err != nil {
		panic(err)
	}
}

func main() {
	loadEnv()

	app := fiber.New(fiber.Config{
		JSONEncoder: json.Marshal,
		JSONDecoder: json.Unmarshal,
		Concurrency: viper.GetInt("MAX_CONCURRENT_CONN"),
	})

	dbMysqlMaster := database.NewMysqlConn(database.MysqlConfig{
		Host: viper.GetString("MYSQL_MASTER_HOST"),
		Port: viper.GetInt("MYSQL_MASTER_PORT"),
		Username: viper.GetString("MYSQL_MASTER_USER"),
		Password: viper.GetString("MYSQL_MASTER_PASS"),
		DatabaseName: viper.GetString("MYSQL_MASTER_DB"),
	})

	dbMongoMaster := database.NewMongoConn(viper.GetString("MONGO_MASTER_URI"), viper.GetString("MONGO_MASTER_DB"))

	app.Use(cors.New(cors.Config{
		AllowOrigins: viper.GetString("CORS_ALLOW_ORIGINS"),
	}))

	app.Use(middleware.ExecutionTimeMiddleware(dbMysqlMaster, viper.GetString("ENVIRONMENT"), viper.GetString("SERVICE")))

	app.Static("/files", "assets/")

	restConfig := &rest.InitRestHandler{}

	restConfig.Mysql.Master = dbMysqlMaster
	restConfig.Mysql.Trx = nil

	restConfig.Mongo.Master = dbMongoMaster
	restConfig.Mongo.Trx = nil

	restConfig.Api.Client = &http.Client{}
	restConfig.Api.GPTApiKey = viper.GetString("GPT_API_KEY")
	restConfig.Api.GeminiApiKey = viper.GetString("GEMINI_API_KEY")

	rest.NewRestHandler(app, restConfig).RegisterRoute()

	err := app.Listen(":"+viper.GetString("SERVICE_PORT"))
	if err != nil {
		panic(err)
	}
}