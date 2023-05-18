package config

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/microcosm-cc/bluemonday"
)

var Sanitbypolicy = func(bp *bluemonday.Policy) func(string) string {
	return (func(input string) string {
		return (bp).Sanitize(input)
	})
}(bluemonday.StrictPolicy())

type envmap map[string]string

var (
	localmap envmap
)

func Getnonempty(key string) string {
	v, ok := localmap[key]
	if !ok || v == "" || key == "" {
		log.Panicln(fmt.Errorf("the env variable %v isn't declared or empty; Check the spelling and its presence in .env file", key))
	}
	return v
}

func init() {
	localmap = getenvmap()
	localmap["ROOT"] = (getroot())
	localmap["psqlconn"] = (getpsqlconn())
}

func getenvmap() envmap {
	err := godotenv.Load(".env")
	if err != nil {
		log.Fatalln(err)
	}

	neededvars := [...]string{
		"HOST", "PORT",
		"dbhost", "dbport", "dbuser", "dbpassword", "dbname",
		"JWTSECRET",
	}

	m := envmap{}
	for _, v := range neededvars {
		val := os.Getenv(v)
		if val == "" {
			log.Fatalln(fmt.Errorf("failed init on env var called %s; not enough environment variables declared;", v))
		}
		m[v] = (val)
	}

	return m
}

func getroot() string {
	return Getnonempty("HOST") + ":" + Getnonempty("PORT")
}

func getpsqlconn() string {
	return fmt.Sprintf(
		"host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
		Getnonempty("dbhost"), Getnonempty("dbport"), Getnonempty("dbuser"), Getnonempty("dbpassword"), Getnonempty("dbname"),
	)
}
