package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/lalathealter/artkeeper/config"
	"github.com/lalathealter/artkeeper/router"
	_ "github.com/lib/pq"
)

func main() {
	fmt.Println("Server starting on port ", config.Getnonempty("PORT"))
	log.Fatal(http.ListenAndServe(config.Getnonempty("ROOT"), router.Use()))
}
