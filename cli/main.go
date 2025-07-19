package main

import "bloggo/internal/app"

func main() {
	application := app.GetInstance()
	application.Bootstrap()
}
