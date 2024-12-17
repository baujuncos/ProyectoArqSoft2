package main

import (
	"fmt"
	"microservicios-api/config"
)

func main() {
	fmt.Println("Iniciando servidor de gestión de microservicios...")
	config.StartRoute()
}
