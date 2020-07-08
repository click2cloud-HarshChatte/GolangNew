//Package Main
package main

//Imported Necessary Files
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

// Define User
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Time string `json:"timezone"`
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
func createtable() {

	db = db_connect()
	fmt.Print(db)
	stmt, err := db.Prepare(`CREATE TABLE users4 (  id SERIAL UNIQUE NOT NULL PRIMARY KEY,  Name TEXT, Time TEXT );`)
	if err != nil {
		fmt.Println(err.Error())
	}
	_, err = stmt.Exec()
	if err != nil {
		fmt.Println(err.Error())
	} else {
		fmt.Println("Table created successfully..")
	}
	fmt.Print(db)
	defer stmt.Close()
	defer db.Close()
}
func GetAllUsers(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var users []User
	server := Result{}
	db = db_connect()
	defer db.Close()
	fmt.Print(db)
	if db != nil {
		result, err := db.Query(`SELECT * from users4`)
		defer result.Close()
		fmt.Print(result)
		fmt.Print(err)
		if err != nil {
			server.Status = false
		} else {
			for result.Next() {
				var user User
				err := result.Scan(&user.Id, &user.Name, &user.Time)
				if err != nil {
					server.Status = false
					panic(err.Error())
				}
				users = append(users, user)
			}
		}
	}

	json.NewEncoder(w).Encode(users)

}

//Rest API for getting user by name
func getUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	server := Result{}
	db = db_connect()
	defer db.Close()

	if db != nil {
			params := mux.Vars(r)
			result, err := db.Query(`SELECT Id FROM users4 WHERE Id = ?`,params["Id"])
		if err != nil {
			server.Status = false
			panic(err.Error())
		}
		var user User
		for result.Next() {
			err := result.Scan(&user.Id, &user.Name, &user.Time)
			if err != nil {
				server.Status = false
				panic(err.Error())
			}
		}
		json.NewEncoder(w).Encode(user)
	}

}

func db_connect() *sql.DB {
	dbHost := "192.168.0.143"
	dbUser := "postgres"
	dbPass := "mysecretpassword"
	dbName := "postgres"
	dbPort := 5432
	// DB_HOST

	psqlInfo := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

	db, err := sql.Open("postgres", psqlInfo)

	if err != nil {
		return nil
	} else {
		return db
	}

}

// Rest API for Entering new user
func createUser(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	result := Result{}
	db = db_connect()
	defer db.Close()
	if db != nil {
		stmt:=`INSERT INTO users4 (Name,Time) VALUES ($1,$2)`

		if err != nil {
			result.Status = false
			panic(err.Error())
		}
		//body, err := ioutil.ReadAll(r.Body)
		//if err != nil {
		//	result.Status = false
		//	panic(err.Error())
		//
		//}
		//keyVal := make(map[string]string)
		//json.Unmarshal(body, &keyVal)
		//Name := keyVal["Name"]
		//Time :=keyVal["Time"]
		//
		_, err = db.Exec(stmt,"Harsh","India")
		result.Status = true
		if err != nil {
			result.Status = false
			panic(err.Error())
		}

		fmt.Fprintf(w, "New User was created")

	} else {
		result.Status = false
	}

}
func main() {

	router := mux.NewRouter()

	// create table func
	createtable()
	// test()

	router.HandleFunc("/users/New", createUser).Methods("POST")
	router.HandleFunc("/users", GetAllUsers).Methods("GET")
	http.ListenAndServe(":8080", router)

}
