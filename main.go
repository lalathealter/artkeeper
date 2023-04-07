package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lalathealter/artkeeper/config"
	"github.com/lalathealter/artkeeper/controllers"
	_ "github.com/lib/pq"
)

func main() {

	controllers.Initdb()
	fmt.Println()
	fmt.Println("Server starting on port ", config.Getnonempty("PORT"))
	http.HandleFunc("/api/addurl", controllers.Posturlhandler)

	log.Fatal(http.ListenAndServe(config.Getnonempty("ROOT"), nil))
}
