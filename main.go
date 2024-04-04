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
	environment       string
	useHTTPGateway    bool
	useAuthMiddleware bool

	// web
	serverPort int

	// app
	logLevel    int
	dbLogLevel  int
	apiEndpoint string

	// aws
	awsRegion          string
	awsStorageS3Bucket string

	// db
	dbHost     string
	dbPort     int
	dbUser     string
	dbPassword string
	dbName     string
	dbUseSSL   bool

	// mail
	host        string
	port        int
	username    string
	password    string
	fromName    string
	fromAddress string

	mailInformCounterBookingCreate     bool
	mailInformCounterBookingExtend     bool
	mailInformCounterBookingDelete     bool
	mailInformCounterBookingCheckedOut bool
	mailInformCounterBookingPaid       bool
	mailInformCounterTimeperiodDelete  bool

	counterMailAddress string
	mailBlacklist      string

	// file path
	assetPath string
)

func main() {
	log.SetFlags(0)
	log.Println()
	log.Println("🚀 Starting Test Backend App")
	log.Println("----------------------------")
	log.Println()

	loadEnvs()
	session, err := initDatabase()

	if err != nil {
		log.Println("Initializing db failed")
		return
	}

	defer session.Close()

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, foo!")
	})

	port := fmt.Sprintf(":%v", serverPort)
	e.Logger.Fatal(e.Start(port))
}

func initDatabase() (*database.Session, error) {
	log.Println("Connecting now")

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

	log.Println("💿️ Connecting to database succefull")

	if session.Connection.AutoMigrate(&models.User{}) != nil {
		log.Println("Migration failed")
	}
	log.Println("Migration succesful")

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
