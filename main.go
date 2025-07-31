package main

import (
	"database/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"net/http"
)

func main() {
	password := "admin123" // hardcoded credential
	db, err := sql.Open("mysql", fmt.Sprintf("root:%s@tcp(localhost:3306)/demo", password))
	if err != nil {
		panic(err)
	}
	defer db.Close()

	http.HandleFunc("/user", func(w http.ResponseWriter, r *http.Request) {
		username := r.URL.Query().Get("username")
		query := "SELECT * FROM users WHERE username = " + username + "'" // SQL Injection vulnerability
		rows, err := db.Query(query)
		if err != nil {
			http.Error(w, "Internal Server Error", 500)
			return
		}
		defer rows.Close()
		w.Write([]byte("User fetched"))
	})

	fmt.Println("Listening on :8080")
	http.ListenAndServe(":8080", nil) // insecure HTTP, should use TLS
}
