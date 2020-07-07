//Package Main
package main

//Imported Necessary Files
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/gorilla/mux"
)

// Define User
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	//Timezone string `json:"timezone"`
}
type Result struct {
	Status bool `json:"status"`
}

//Defined variable for database and error
var db *sql.DB
var err error

//  type Person struct {
//     Name string `json:"name"`
//  }
// Rest API for getting all users
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	server := Result{}
	result, err := db.Query("SELECT * from users ORDER BY id desc")
	if err != nil {
		server.Status = false
		panic(err.Error())
	}
	defer result.Close()

	for result.Next() {
		var user User
		err := result.Scan(&user.Id, &user.Name)
		if err != nil {
			server.Status = false
			panic(err.Error())
		}
		users = append(users, user)
	}
	json.NewEncoder(w).Encode(users)
}

//Rest API for getting user by name
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	server := Result{}
	result, err := db.Query("SELECT * FROM users WHERE Name = ?")
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&result); err != nil {
		json.NewEncoder(w).Encode(err)
	}
	if err != nil {
		server.Status = false
		panic(err.Error())
	}
	defer result.Close()
	var user User
	for result.Next() {
		err := result.Scan(&user.Id, &user.Name)
		if err != nil {
			server.Status = false
			panic(err.Error())
		}
	}
	json.NewEncoder(w).Encode(user)
}

// Rest API for Entering new user
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result := Result{}
	stmt, err := db.Prepare("INSERT INTO users(Name) VALUES(?)")
	if err != nil {
		result.Status = false
		panic(err.Error())
	}
	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		result.Status = false
		panic(err.Error())

	}
	keyVal := make(map[string]string)
	json.Unmarshal(body, &keyVal)
	Name := keyVal["Name"]
	decoder := json.NewDecoder(r.Body)
	if err := decoder.Decode(&Name); err != nil {
		json.NewEncoder(w).Encode(err)
		fmt.Print(err)
	}

	_, err = stmt.Exec(Name)
	result.Status = true
	if err != nil {
		result.Status = false
		panic(err.Error())
	}

	fmt.Fprintf(w, "New User was created")
}

func main() {
	dbHost := "192.168.0.196"
	dbUser := os.Getenv("DB_USER")
	dbPass := ""
	dbName := "restapi"
	dbPort := os.Getenv("DB_PORT")
	// DB_HOST
	db, err = sql.Open("mysql", dbUser+":"+dbPass+"@tcp("+dbHost+":"+dbPort+")/"+dbName)
	defer db.Close()
	router := mux.NewRouter()
	router.HandleFunc("/users/New", createUser).Methods("POST")
	router.HandleFunc("/users/Name", getUser).Methods("GET")
	router.HandleFunc("/users", GetAllUsers).Methods("GET")
	http.ListenAndServe(":8092", router)

}
