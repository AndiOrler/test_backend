package main

import (
	"log"
	"net/http"
	"os"
	"test_backend/database"
	"test_backend/models"
	"test_backend/utils/env"

	"github.com/labstack/echo/v4"
)

var (
	environment       string
	useHTTPGateway    bool
	useAuthMiddleware bool

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
	initialize()

	port := os.Getenv("PORT")

	e := echo.New()
	e.GET("/", func(c echo.Context) error {
		return c.String(http.StatusOK, "Hello, foo!")
	})

	e.Logger.Fatal(e.Start(":" + port))

}

func initialize() {
	log.Println("Initializing App")
	loadEnvs()
	connectToDB()
}

func connectToDB() (*database.Session, error) {
	log.Println("Connecting now")

	session, err := database.NewSession(&database.Parameters{
		Host:                   dbHost,
		Port:                   dbPort,
		User:                   dbUser,
		Password:               dbPassword,
		Name:                   dbName,
		SkipDefaultTransaction: false,
		PrepareStmt:            true,
		LogLevel:               3,
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

	dbHost, _ = env.GetStr("POSTGRES_HOST")
	dbPort, _ = env.GetInt("POSTGRES_DB_P0RT")
	dbUser, _ = env.GetStr("POSTGRES_USER")
	dbPassword, _ = env.GetStr("POSTGRES_PASSWORD")
	dbName, _ = env.GetStr("POSTGRES_DB")
	dbUseSSL, _ = env.GetBool("DB_USE_SSL")
}
