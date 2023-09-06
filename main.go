package main

import (
	"doctogadget/doctowidget"
)

func main() {
	dw := doctowidget.NewDoctowidget()
	defer dw.Destroy()

	dw.Show()
	dw.Run()
}

// func StartGadget(void (*doctoWidgetCallback)(const char *text),
//                    void (*initPromptCallback)(const char *text),
//                    void (*workingCallback)(const char *text),
//                    void (*retryPromptCallback)(const char *text), int lang = 0) {

// }
