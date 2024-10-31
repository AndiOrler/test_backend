package main

import (
	"fmt"
	"log"
	"net/http"
	"test_backend/database"
	"test_backend/models"
	"test_backend/utils/env"

	"github.com/labstack/echo/v4"
)

var (
	dbUser     string
	dbPassword string
	dbName     string
	dbUseSSL   bool
	dbHost     string
	dbPort     int
	serverPort int
	dbLogLevel int
)

func main() {
	// no timestamp on logs
	log.SetFlags(0)

	log.Println()
	log.Println("🚀 Starting Test Backend App")
	log.Println("----------------------------")
	log.Println()

	// loading env
	loadEnvs()

	// connecting to db
	session, err := initDatabase()

	if err != nil {
		log.Println("Initializing db failed")
		return
	}

	defer session.Close()

	err = migrateDB(session)

	if err != nil {
		log.Println("DB migration failed")
		return
	}

	// userRepo := user.NewUserRepo(session.Connection)
	// userService := user.NewUserService(userRepo)
	// userRepo.Get//userRepo.Get

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, foo!")
	})

	port := fmt.Sprintf(":%v", serverPort)
	e.Logger.Fatal(e.Start(port))
}

func initDatabase() (*database.Session, error) {
	log.Println("Connecting to postgresdb")

	session, err := database.NewSession(&database.Parameters{
		Host:                   dbHost,
		Port:                   dbPort,
		User:                   dbUser,
		Password:               dbPassword,
		Name:                   dbName,
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
		LogLevel:               dbLogLevel,
		UseSSL:                 dbUseSSL,
	})

	if err != nil {
		return nil, err
	}

	log.Println("💿️ Connecting to database successful")

	return session, err
}

func loadEnvs() {
	log.Println("Reading env variables from .env")
	_ = env.LoadLocalFile()

	dbHost, _ = env.GetStr("DB_HOST")
	dbPort, _ = env.GetInt("DB_P0RT")
	dbUser, _ = env.GetStr("DB_USER")
	dbPassword, _ = env.GetStr("DB_PASSWORD")
	dbName, _ = env.GetStr("DB_NAME")
	dbUseSSL, _ = env.GetBool("DB_USE_SSL")
	serverPort, _ = env.GetInt("PORT")
	dbLogLevel, _ = env.GetInt("DB_LOG_LEVEL")
}

func migrateDB(session *database.Session) error {
	err := session.Connection.AutoMigrate(&models.User{})

	if err != nil {
		log.Println("Migration failed")
	}

	log.Println("Migration successful")

	return err
}
