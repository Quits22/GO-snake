package main

//here we are making the webapp
import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"os"

	"github.com/joho/godotenv"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// here we define the router i used mux and also a hello func to make sure everything is in order
func newRouter() *mux.Router {
	r := mux.NewRouter()
	r.HandleFunc("/hello", handler).Methods("GET")

	// Declare the static file directory and point it to the directory with the static webpage
	staticFileDirectory := http.Dir("./assets/")
	//the `stripPrefix` is used to remove the extra /assets/
	staticFileHandler := http.StripPrefix("/assets/", http.FileServer(staticFileDirectory))
	// The scores are shown here
	r.PathPrefix("/assets/").Handler(staticFileHandler).Methods("GET")
	// getting the data from and to /score
	r.HandleFunc("/score", getScoreHandler).Methods("GET")
	r.HandleFunc("/score", createScoreHandler).Methods("POST")
	return r
}

// this is for the server info change it in the go.env
func loadEnv() {
	err := godotenv.Load("go.env")
	if err != nil {
		log.Fatal("Failed to load .env file:", err)
	}
}

// change these in the go.env to connect to your server
func initializeDB() (*sql.DB, error) {
	dbUsername := os.Getenv("DB_USERNAME")
	dbPassword := os.Getenv("DB_PASSWORD")
	dbHost := os.Getenv("DB_HOST")
	dbPort := os.Getenv("DB_PORT")
	dbName := os.Getenv("DB_NAME")

	db, err := sql.Open("mysql", fmt.Sprintf("%s:%s@tcp(%s:%s)/%s", dbUsername, dbPassword, dbHost, dbPort, dbName))
	if err != nil {
		return nil, fmt.Errorf("failed to open database connection: %w", err)
	}
	//just checking for connection
	if err := db.Ping(); err != nil {
		db.Close()
		return nil, fmt.Errorf("failed to ping the database: %w", err)
	}

	return db, nil
}

func main() {
	fmt.Println("Starting server...")

	loadEnv()

	db, err := initializeDB()
	if err != nil {
		log.Fatalf("Fatal error: %s", err)
	}
	defer db.Close()
	//getting the scores from the database
	InitStore(&dbStore{db: db})

	// the port
	r := newRouter()
	serverPort := os.Getenv("SERVER_PORT")
	if serverPort == "" {
		serverPort = "8089" // Default port if not specified
	}

	serverAddr := fmt.Sprintf(":%s", serverPort)
	log.Printf("Serving on port %s", serverPort)
	log.Fatal(http.ListenAndServe(serverAddr, r))
}

func handler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "Hello World!")
}
