package main

import (
	"context"
	"fmt"
	"os"

	tea "charm.land/bubbletea/v2"
)

func run() error {
	ctx, cancel := context.WithCancel(context.Background())
	program := tea.NewProgram(newModel(), tea.WithContext(ctx))
	done := make(chan struct{})
	go func() {
		defer close(done)
		watchLayers(ctx, func(event layerEvent) { program.Send(event) })
	}()
	defer func() {
		cancel()
		<-done
	}()
	_, err := program.Run()
	return err
}

func main() {
	if err := run(); err != nil {
		fmt.Fprintln(os.Stderr, "Falha no monitor de camadas:", err)
		os.Exit(1)
	}
}
