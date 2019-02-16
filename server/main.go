package main

import (
	"fmt"
	"net/http"

	"github.com/gorilla/mux"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/sqlite"
	_ "github.com/mattn/go-sqlite3"
	"github.com/rs/cors"
)

func main() {
	c := loadConfig()

	db, _ := gorm.Open("sqlite3", "./gorm.db")
	defer db.Close()

	db.AutoMigrate(&Player{}, &LadderRecord{}, &Ranking{})

	s := c.serverSetup(db)

	cor := cors.New(cors.Options{
		AllowedOrigins:   []string{"*"},
		AllowedMethods:   []string{"GET", "PUT", "POST", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"Authorization"},
		AllowCredentials: true,

		Debug: true,
	})

	r := mux.NewRouter()
	connectRoutes(r, &s)

	handler := cor.Handler(r)
	fmt.Print("Started running on " + c.domain + ":" + c.port + "\n")
	fmt.Println(http.ListenAndServe(":"+c.port, handler))
}
