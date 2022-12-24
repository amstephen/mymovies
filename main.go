package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
)

//APIKEY variable where you have to paste your API key
const APIKEY = "https://api.themoviedb.org/3/movie/550?api_key=6ee506ec9fb480af76ba4862253e92e8"

//APIURL for API
const APIURL = "https://api.themoviedb.org/3/movie/popular?api_key=" + APIKEY + "&page=1"

//Change your database configuration below

//DBNAME variable for database name
const DBNAME = "testing"

//DBUSERNAME variable for database username
const DBUSERNAME = "root"

//DBPASSWORD variable for database password
const DBPASSWORD = "test123"

//DBHOST variable for database address
const DBHOST = "127.0.0.1"

//DBPORT variable for database default port number
const DBPORT = "3306"

//DATABASEURL variable
const DATABASEURL = DBUSERNAME + ":" + DBPASSWORD + "@tcp(" + DBHOST + ":" + DBPORT + ")/" + DBNAME + "?charset=utf8&parseTime=True&loc=Local"

type movie struct {
	MovieID     int     `gorm:"column:movie_id;primary_key" json:"id"`
	Title       string  `json:"title"`
	ReleaseDate string  `json:"release_date"`
	Language    string  `json:"original_language"`
	Adult       bool    `json:"adult"`
	Image       string  `json:"poster_path"`
	Overview    string  `gorm:"type:varchar(1000)" json:"overview"`
	VoteAverage float32 `json:"vote_average"`
}

type movieList struct {
	List []movie `json:"results"`
}

func (list movieList) save() {
	db, dberr := gorm.Open("mysql", DATABASEURL)
	
	if dberr != nil {
		log.Fatal(dberr)
	}
	defer db.Close()
	db.Debug().DropTableIfExists(&movie{})
	db.AutoMigrate(&movie{})
	for _, row := range list.List {
		db.Debug().Create(&row)
	}
}

func tmdbImplementation(w http.ResponseWriter, r *http.Request) {
	//For receiving API call
	client := http.Client{}
	movieRequest, httperr := http.NewRequest(http.MethodGet, APIURL, nil)
	if httperr != nil {
		log.Fatal(httperr)
	}
	movieResponse, geterr := client.Do(movieRequest)
	if geterr != nil {
		log.Fatal(geterr)
	}
	movieBody, readerr := io.ReadAll(movieResponse.Body)
	if readerr != nil {
		log.Fatal(readerr)
	}
	list := movieList{}
	jsonerr := json.Unmarshal(movieBody, &list)
	if jsonerr != nil {
		log.Fatal(jsonerr)
	}
	list.save()
	//For sending API call
	w.Header().Set("Access-Control-Allow-Origin", "*") //This heading is necessary for cross origin data transfer
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(list)

}

func main() {
	router := mux.NewRouter()
	router.HandleFunc("/", tmdbImplementation)
	fmt.Println("Listening of port 8080")
	http.ListenAndServe(":8080", router)
}