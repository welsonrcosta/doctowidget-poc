package main

import (
	"doctogadget/app"
	"doctogadget/internal/nativemessage"
	"fmt"
	"log"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	in := make(chan interface{})
	out := make(chan nativemessage.QtToZDMessage)

	wg.Add(1)
	go func() {
		for m := range out {
			fmt.Printf("%v", m)
		}
		wg.Done()
	}()

	wg.Add(1)
	go setupNativeMessageReader(in)

	wg.Add(1)
	go startGtkApp(in, out, &wg)
	wg.Wait()
}

func setupNativeMessageReader(in chan interface{}) {
	reader := nativemessage.NewNativeMessageReader()
	for {
		m, err := reader.ReadMessage()

		if err != nil {
			log.Println("Error in read string", err)
			continue
		}
		// TODO native message marshal
		in <- m
	}
}

func startGtkApp(in chan interface{}, out chan nativemessage.QtToZDMessage, wg *sync.WaitGroup) {
	app := app.NewApp(in, out)
	app.Run()
	wg.Done()
	close(in)
	close(out)
}
