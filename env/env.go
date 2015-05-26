package env

import (
	"github.com/joho/godotenv"
	"log"
	"os"
	"strconv"
	"strings"
)

var Env map[string]string

func init() {

	// Init in memory env
	Env = make(map[string]string)

	// Use OS environement
	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")
		Env[pair[0]] = pair[1]
	}

	// Read .env file
	env, err := godotenv.Read()
	if err != nil {
		log.Println(err)
	}

	for k, v := range env {
		Env[k] = v
	}
}

// Return a value from environement
func Get(name string) string {
	v, ok := Env[name]
	if !ok {
		return ""
	}

	return v
}

func GetInt(name string) int {

	v := Get(name)

	nb, err := strconv.Atoi(v)
	if err != nil {
		return 0
	}

	return nb
}
