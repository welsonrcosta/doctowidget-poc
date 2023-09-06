package doctowidget

import (
	"doctogadget/internal/assets"
	"errors"
	"fmt"
	"log"
	"os"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type doctowidget struct {
	app        *gtk.Application
	mainWindow *gtk.Window
	resize     *gtk.Button
}

const appId = "com.github.gotk3.gotk3-examples.glade"

func NewDoctowidget() doctowidget {
	app, err := gtk.ApplicationNew(appId, glib.APPLICATION_FLAGS_NONE)
	checkError(err)

	dw := doctowidget{}
	dw.app = app

	app.Connect("activate", func() {

		log.Println("application activate")

		builder, err := gtk.BuilderNewFromString(assets.Ui)
		checkError(err)

		//Connect button signals to doctowidget functions
		signals := map[string]interface{}{
			"on_open_doctolib_button_clicked": dw.onOpenButtonClicked,
			"on_book_button_clicked":          dw.onBookButtonClicked,
			"on_history_button_clicked":       dw.onHistoryButtonClicked,
			"on_list_button_clicked":          dw.onListButtonClicked,
			"on_resize_button_clicked":        dw.onResizeButtonClicked,
		}
		builder.ConnectSignals(signals)

		win, err := getMainWindow(builder)
		checkError(err)
		app.AddWindow(win)
		dw.mainWindow = win

		setupDoctoligLogo(builder)
		resize, err := setupResizeButton(builder)
		checkError(err)
		dw.resize = resize

		win.ShowAll()
		loadCSS()

		builder.Unref()
	})

	return dw
}

func loadCSS() {
	cssProvider, err := gtk.CssProviderNew()
	checkError(err)
	err = cssProvider.LoadFromData(assets.Style)
	checkError(err)
	screen, err := gdk.ScreenGetDefault()
	checkError(err)
	gtk.AddProviderForScreen(screen, cssProvider, gtk.STYLE_PROVIDER_PRIORITY_APPLICATION)
}

func (dw doctowidget) Show() {
	if dw.mainWindow != nil {
		dw.mainWindow.ShowAll()
	}
}
func (dw doctowidget) Hide() {}
func (dw doctowidget) Run() {
	dw.app.Run(os.Args)
}

func (dw *doctowidget) Destroy() {
	dw.app.Unref()
	dw.app = nil
	dw.mainWindow = nil
}

func checkError(err error) {
	if err != nil {
		log.Fatal("An error has occured:", err)
	}
}

func loadImage(image string, w int, h int) *gtk.Image {
	//buff, err := gdk.PixbufNewFromFileAtScale(image, w, h, false)
	buff := pixBuffAtScale(image, w, h)
	img, err := gtk.ImageNewFromPixbuf(buff)
	checkError(err)
	img.SetSizeRequest(w, h)
	return img
}

func pixBuffAtScale(image string, w int, h int) *gdk.Pixbuf {
	data, err := assets.Images.ReadFile(image)
	checkError(err)

	buff, err := gdk.PixbufNewFromBytesOnly(data)
	checkError(err)
	buff, err = buff.ScaleSimple(w, h, gdk.INTERP_BILINEAR)
	checkError(err)
	return buff
}

func getMainWindow(builder *gtk.Builder) (*gtk.Window, error) {
	obj, err := builder.GetObject("main_window")
	checkError(err)

	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func setupDoctoligLogo(builder *gtk.Builder) {
	obj, err := builder.GetObject("doctolib_logo")
	checkError(err)
	if logo, ok := obj.(*gtk.DrawingArea); ok {
		logoW, rh := logo.GetSizeRequest()
		buff := pixBuffAtScale("logo_docto.svg", int(float64(logoW)*0.8), int(float64(rh)*0.8))
		logo.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
			gtk.GdkCairoSetSourcePixBuf(cr, buff, float64(logoW)*0.1, float64(rh)*0.3)
			cr.Paint()
		})
	}
}

func setupResizeButton(builder *gtk.Builder) (*gtk.Button, error) {
	obj, err := builder.GetObject("resize_button")
	checkError(err)
	if btn, ok := obj.(*gtk.Button); ok {
		w, h := btn.GetSizeRequest()
		img := loadImage("dwindle.svg", w, h)
		btn.SetImage(img)
		return btn, nil
	}
	return nil, errors.New("could not configure resize button")
}

func (dw doctowidget) onResizeButtonClicked() {
	fmt.Println("Resize")
}

func (dw doctowidget) onOpenButtonClicked() {
	fmt.Println("Open")
}

func (dw doctowidget) onHistoryButtonClicked() {
	fmt.Println("History")
}

func (dw doctowidget) onListButtonClicked() {
	fmt.Println("List")
}

func (dw doctowidget) onBookButtonClicked() {
	fmt.Println("Book")
}
