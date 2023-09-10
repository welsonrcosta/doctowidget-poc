package assets

import (
	"doctogadget/internal/util"
	"embed"
	_ "embed"
	"fmt"

	"github.com/gotk3/gotk3/gdk"
	"github.com/gotk3/gotk3/gtk"
)

//go:embed style.css
var Style string

//go:embed *.svg
var images embed.FS

//go:embed *.glade
var builders embed.FS

func LoadImage(image string, w int, h int) *gtk.Image {
	//buff, err := gdk.PixbufNewFromFileAtScale(image, w, h, false)
	buff := PixBuffAtScale(image, w, h)
	img, err := gtk.ImageNewFromPixbuf(buff)
	util.CheckError(err)
	img.SetSizeRequest(w, h)
	return img
}

func PixBuffAtScale(image string, w int, h int) *gdk.Pixbuf {
	data, err := images.ReadFile(image)
	util.CheckError(err)

	buff, err := gdk.PixbufNewFromBytesOnly(data)
	util.CheckError(err)
	buff, err = buff.ScaleSimple(w, h, gdk.INTERP_BILINEAR)
	util.CheckError(err)
	return buff
}

func GetUIBuilder(name string) (*gtk.Builder, error) {
	data, err := builders.ReadFile(name)
	if err != nil {
		return nil, err
	}

	builder, err := gtk.BuilderNewFromString(string(data))
	if err != nil {
		return nil, fmt.Errorf("could not load builder: %s", name)
	}

	return builder, nil
}
