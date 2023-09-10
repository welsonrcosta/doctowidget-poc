package main

import (
	"bufio"
	"doctogadget/app"
	"fmt"
	"log"
	"os"
	"sync"
)

func main() {
	var wg sync.WaitGroup
	in := make(chan string)
	out := make(chan string)

	wg.Add(1)
	go func() {
		for m := range out {
			fmt.Println(m)
		}
		wg.Done()
	}()

	wg.Add(1)
	go setupNativeMessageReader(in)

	wg.Add(1)
	go startGtkApp(in, out, &wg)
	wg.Wait()
}

func setupNativeMessageReader(in chan string) {
	reader := bufio.NewReader(os.Stdin)
	for {
		s, err := reader.ReadString('\n')

		if err != nil {
			log.Println("Error in read string", err)
		}
		// TODO native message marshal
		in <- s[:len(s)-1]
	}
}

func startGtkApp(in chan string, out chan string, wg *sync.WaitGroup) {
	app := app.NewApp(in, out)
	app.Run()
	wg.Done()
	close(in)
	close(out)
}
