package main

import (
	_ "graduationEnd/middleware"
	"graduationEnd/router"
)

func main() {
	router.StartRouter()
}
