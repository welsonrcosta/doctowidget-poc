package widgetwindow

import (
	"doctogadget/internal/assets"
	"doctogadget/internal/util"
	"encoding/json"
	"errors"
	"fmt"
	"log"

	"github.com/gotk3/gotk3/cairo"
	"github.com/gotk3/gotk3/glib"
	"github.com/gotk3/gotk3/gtk"
)

type WidgetWindow struct {
	ready       bool
	mainWindow  *gtk.Window
	resize      *gtk.Button
	in          chan string
	out         chan string
	isMaximized bool
}

func NewDoctowidget(in chan string, out chan string) WidgetWindow {
	dw := WidgetWindow{in: in, out: out, isMaximized: true, ready: true}
	return dw
}

func (dw *WidgetWindow) ActivateDoctowidget(app *gtk.Application) error {
	if !dw.ready {
		return fmt.Errorf("instance was not activated")
	}

	builder, err := assets.GetUIBuilder("ui.glade")
	util.CheckError(err)

	win, err := getMainWindow(builder)
	util.CheckError(err)
	app.AddWindow(win)
	dw.mainWindow = win

	setupDoctoligLogo(builder)
	resize, err := setupResizeButton(builder)
	util.CheckError(err)
	dw.resize = resize

	//Connect button signals to Doctowidget functions
	signals := map[string]interface{}{
		"on_open_doctolib_button_clicked": dw.onOpenButtonClicked,
		"on_book_button_clicked":          dw.onBookButtonClicked,
		"on_history_button_clicked":       dw.onHistoryButtonClicked,
		"on_list_button_clicked":          dw.onListButtonClicked,
		"on_resize_button_clicked":        dw.onResizeButtonClicked,
	}
	builder.ConnectSignals(signals)

	go func() {
		for m := range dw.in {
			log.Printf("received %s\n", m)
			switch m {
			case "show":
				{
					dw.Show()
				}
			case "hide":
				{
					dw.Hide()
				}
			}
		}
	}()

	builder.Unref()

	return nil
}

func (dw WidgetWindow) Show() {
	var f = func() {
		dw.mainWindow.ShowAll()
		dw.mainWindow.SetKeepAbove(true)
	}
	glib.IdleAdd(f)
}
func (dw WidgetWindow) Hide() {
	glib.IdleAdd(dw.mainWindow.Hide)
}

func (dw *WidgetWindow) Destroy() {

}

func getMainWindow(builder *gtk.Builder) (*gtk.Window, error) {
	obj, err := builder.GetObject("main_window")
	util.CheckError(err)

	if win, ok := obj.(*gtk.Window); ok {
		return win, nil
	}
	return nil, errors.New("not a *gtk.Window")
}

func setupDoctoligLogo(builder *gtk.Builder) {
	obj, err := builder.GetObject("doctolib_logo")
	util.CheckError(err)
	if logo, ok := obj.(*gtk.DrawingArea); ok {
		logoW, rh := logo.GetSizeRequest()
		buff := assets.PixBuffAtScale("logo_docto.svg", int(float64(logoW)*0.8), int(float64(rh)*0.8))
		logo.Connect("draw", func(da *gtk.DrawingArea, cr *cairo.Context) {
			gtk.GdkCairoSetSourcePixBuf(cr, buff, float64(logoW)*0.1, float64(rh)*0.3)
			cr.Paint()
		})
	}
}

func setupResizeButton(builder *gtk.Builder) (*gtk.Button, error) {
	obj, err := builder.GetObject("resize_button")
	util.CheckError(err)
	if btn, ok := obj.(*gtk.Button); ok {
		w, h := btn.GetSizeRequest()
		img := assets.LoadImage("dwindle.svg", w, h)
		btn.SetImage(img)
		return btn, nil
	}
	return nil, errors.New("could not configure resize button")
}

func (dw *WidgetWindow) onResizeButtonClicked() {
	fmt.Printf("Resize %v\n", dw.isMaximized)
	w, h := dw.resize.GetSizeRequest()
	dw.isMaximized = !dw.isMaximized
	if dw.isMaximized {
		img := assets.LoadImage("dwindle.svg", w, h)
		dw.resize.SetImage(img)
	} else {
		img := assets.LoadImage("enlarge.svg", w, h)
		dw.resize.SetImage(img)
	}
}

func (dw *WidgetWindow) onOpenButtonClicked() {
	x, y := dw.mainWindow.GetPosition()
	move, err := json.Marshal(struct {
		Command string `json:"command"`
		X       int    `json:"x"`
		Y       int    `json:"y"`
	}{
		Command: "set_position",
		X:       x,
		Y:       y,
	})
	util.CheckError(err)
	dw.out <- string(move)
}

func (dw WidgetWindow) onHistoryButtonClicked() {
	fmt.Println("History")
}

func (dw WidgetWindow) onListButtonClicked() {
	fmt.Println("List")
}

func (dw WidgetWindow) onBookButtonClicked() {
	fmt.Println("Book")
}
