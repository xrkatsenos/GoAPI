package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"strconv"
	"time"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	"github.com/joho/godotenv"
)

type TempHumValue struct {
	gorm.Model
	Temp     float32
	Humidity float32
}

// func params(w http.ResponseWriter, r *http.Request) {
// 	pathParams := mux.Vars(r)
// 	w.Header().Set("Content-Type", "application/json")

// 	userID := -1
// 	var err error
// 	if val, ok := pathParams["userID"]; ok {
// 		userID, err = strconv.Atoi(val)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(`{"message": "need a number"}`))
// 			return
// 		}
// 	}

// 	commentID := -1
// 	if val, ok := pathParams["commentID"]; ok {
// 		commentID, err = strconv.Atoi(val)
// 		if err != nil {
// 			w.WriteHeader(http.StatusInternalServerError)
// 			w.Write([]byte(`{"message": "need a number"}`))
// 			return
// 		}
// 	}

// 	query := r.URL.Query()
// 	location := query.Get("location")

// 	w.Write([]byte(fmt.Sprintf(`{"userID": %d, "commentID": %d, "location": "%s" }`, userID, commentID, location)))
// }

func store(w http.ResponseWriter, r *http.Request) {

	println(time.Now().Format("2006-01-02 15:04:05"))
	println("Request:", r.URL.RequestURI())
	println("Method:", r.Method)
	println("")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	message := "Values failed to stored"

	db, err := gorm.Open(os.Getenv("DBDRIVER"), os.Getenv("DBNAME"))
	if err != nil {
		message = "failed to connect database"
		w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, message)))
		println("PANIC")
		panic(message)
	}
	defer db.Close()

	temp, err := strconv.ParseFloat(r.FormValue("temp"), 2)
	hum, err2 := strconv.ParseFloat(r.FormValue("humidity"), 2)

	if err == nil && err2 == nil {
		db.Create(&TempHumValue{Temp: float32(temp), Humidity: float32(hum)})
		//println("Values stored")
		message = "Values stored"
	}

	w.Write([]byte(fmt.Sprintf(`{"message": "%s"}`, message)))
}

func getValues(w http.ResponseWriter, r *http.Request) {

	println(time.Now().Format("2006-01-02 15:04:05"))
	println("Request:", r.URL.RequestURI())
	println("Method:", r.Method)
	println("")

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	message := ""

	db, err := gorm.Open(os.Getenv("DBDRIVER"), os.Getenv("DBNAME"))
	if err != nil {
		message = "failed to connect database"
		println("PANIC")
		panic(message)
	}
	defer db.Close()

	vals := []TempHumValue{}
	db.Find(&vals)

	message = fmt.Sprintf(`%d records found`, len(vals))

	responce, err := json.Marshal(vals)
	if err != nil {
		fmt.Println(err)
		message = err.Error()
		return
	}
	//fmt.Println(string(responce))

	w.Write([]byte(fmt.Sprintf(`{"message": "%s","data": %s}`, message, responce)))
}

func init() {
	err := godotenv.Load("./config/.env")
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	db, err := gorm.Open(os.Getenv("DBDRIVER"), os.Getenv("DBNAME"))
	if err != nil {
		panic("failed to connect database")
	}
	defer db.Close()

	// Migrate the schema
	db.AutoMigrate(&TempHumValue{})
}

func main() {
	r := mux.NewRouter()

	api := r.PathPrefix("/api/v1").Subrouter()
	//api.HandleFunc("", get).Methods(http.MethodGet)
	//api.HandleFunc("", post).Methods(http.MethodPost)
	//api.HandleFunc("", put).Methods(http.MethodPut)
	//api.HandleFunc("", delete).Methods(http.MethodDelete)
	//api.HandleFunc("/user/{userID}/comment/{commentID}", params).Methods(http.MethodGet)

	api.HandleFunc("/store", store).Methods(http.MethodPost)
	api.HandleFunc("/read", getValues).Methods(http.MethodGet)

	println("")
	println("================================")
	println("Server started on port 8080")
	println("and is running in the following URL:")
	println("http://localhost:8080/api/v1/")
	println("================================")
	println("")
	println("")

	log.Fatal(http.ListenAndServe(":8080", r))

}
