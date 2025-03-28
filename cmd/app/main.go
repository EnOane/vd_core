package main

import (
	"core/internal/di"
	"core/internal/infr/broker"
)

func main() {
	di.MakeDIContainer()

	broker.MustConnect()
}
