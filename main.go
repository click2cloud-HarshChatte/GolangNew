//Package Main
package main

//Imported Necessary Files
import (
	"database/sql"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
	"net/http"
	"os"
)

// Define User
type Connection struct {
	Status bool `json:"status"`
	Users []User
}
type User struct {
	Id   int    `json:"id"`
	Name string `json:"name"`
	Time string `json:"timezone"`
}



//Defined variable for database and error
var db *sql.DB
var err error


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
	 Conn:= Connection{}
	//Conn := Connection{}
	var users  []User

	db = db_connect()
	defer db.Close()
	fmt.Print(db)
	if db != nil {
        Conn.Status=true
		result, err := db.Query(`SELECT * from users4`)
		defer result.Close()
		fmt.Print(result)
		fmt.Print(err)
		if err != nil {
			Conn.Status = false
		} else {
			for result.Next() {
				var user User
				err := result.Scan(&user.Id, &user.Name, &user.Time)
				if err != nil {
					Conn.Status = false
					panic(err.Error())
				}
				users = append(users, user)
				Conn.Users=append(users,user)
			}

		}

	}


	json.NewEncoder(w).Encode(Conn)

}



func db_connect() *sql.DB {
	dbHost := os.Getenv("DB_HOST")
	dbUser := os.Getenv("DB_USER")
	dbPass := os.Getenv("DB_PASS")
	dbName := os.Getenv("DB_NAME")
	dbPort := os.Getenv("DB_PORT")
	// DB_HOST

	psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable", dbHost, dbPort, dbUser, dbPass, dbName)

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
    Conn:=Connection{}
	u:=new(User)
	db = db_connect()
	defer db.Close()
	if db != nil {
		Conn.Status=true
		stmt:=`INSERT INTO users4 (Name,Time) VALUES ($1,$2)`

		if err != nil {
          Conn.Status=false
			panic(err.Error())
		}
		decoder := json.NewDecoder(r.Body)
		if err := decoder.Decode(&u); err != nil {
			json.NewEncoder(w).Encode(err)
		}


		_, err = db.Exec(stmt,u.Name,u.Time)

		if err != nil {

			panic(err.Error())
		}
            Conn.Status=true
		json.NewEncoder(w).Encode(Conn)
		json.NewEncoder(w).Encode("New User Created Successfully")

	} else {
		panic(err.Error())
	}

}
func main() {

	router := mux.NewRouter()


	createtable()

	router.HandleFunc("/users/New", createUser).Methods("POST")
	router.HandleFunc("/users", GetAllUsers).Methods("GET")
	http.ListenAndServe(":8080", router)

}
