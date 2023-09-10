package app

import (
	"doctogadget/internal/assets"
	"doctogadget/internal/util"
	"doctogadget/internal/widgetwindow"
	"log"
	"os"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

const appId = "com.doctolib.widget"

type DoctoWidgetsApp struct {
	app   *gtk.Application
	dw    widgetwindow.WidgetWindow
	ready bool
}

func NewApp(in chan string, out chan string) DoctoWidgetsApp {
	app, err := gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	util.CheckError(err)

	widgetsApp := DoctoWidgetsApp{ready: true}
	widgetsApp.app = app

	widgetsApp.initAppWindows(in, out)

	return widgetsApp
}

func (wa *DoctoWidgetsApp) initAppWindows(in chan string, out chan string) {
	dw := widgetwindow.NewDoctowidget(in, out)
	wa.dw = dw

	wa.app.Connect("activate", func() {
		loadCSSProvider()
		err := dw.ActivateDoctowidget(wa.app)
		util.CheckError(err)
	})
}

func (app DoctoWidgetsApp) Run() {
	if !app.ready {
		log.Fatal("attempted to run unitialized app")
	}
	app.app.Run(os.Args)
}

func loadCSSProvider() {
	cssProvider, err := gtk.CssProviderNew()
	util.CheckError(err)
	err = cssProvider.LoadFromData(assets.Style)
	util.CheckError(err)
	screen, err := gdk.ScreenGetDefault()
	util.CheckError(err)
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}
