package main

import (
	"go_project_demo/router"
)

func main() {

	r := router.Router()

	r.Run(":9999")
}
