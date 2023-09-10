package main

/*
#include <stdlib.h>
*/
import "C"

import (
	"bufio"
	"doctogadget/doctowidget"
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
	go func() {
		reader := bufio.NewReader(os.Stdin)
		for {
			s, err := reader.ReadString('\n')

			if err != nil {

				//close(messages)
				log.Println("Error in read string", err)
			}
			in <- s[:len(s)-1]
		}
	}()
	wg.Add(1)
	go startGtkApp(in, out, &wg)
	wg.Wait()
}

func startGtkApp(in chan string, out chan string, wg *sync.WaitGroup) {
	dw := doctowidget.NewDoctowidget(in, out)
	defer dw.Destroy()
	dw.Run()
	wg.Done()
	close(in)
	close(out)
}
