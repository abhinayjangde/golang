package handlers

import (
	"database/sql"
	"fmt"
	"net/http"
)

func List(db *sql.DB) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {

		result, err := db.Exec("SELECT * FROM listing")
		if err != nil {
			http.Error(w, "Failed to retrieve listings", http.StatusInternalServerError)
			return
		}

		fmt.Println(result.RowsAffected())

		w.Write([]byte("get all listings"))
	}
}
