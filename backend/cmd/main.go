package main

import (
	"fmt"

	"battle-game/internal/battle"
	"battle-game/internal/engine"
)

func main() {
	fmt.Println("===================================")
	fmt.Println("      Battle Engine v0.1")
	fmt.Println("===================================")

	b := battle.NewDemoBattle()

	e := engine.New(b)

	e.Run()

	fmt.Println()

	fmt.Println("============= Events =============")

	for _, event := range b.Events {
		fmt.Println(event)
	}

	fmt.Println()

	fmt.Println("============= Result =============")

	if b.Left.AliveCount() > 0 {
		fmt.Println("Winner: Left Team")
	} else {
		fmt.Println("Winner: Right Team")
	}
}
