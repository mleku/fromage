package main

import (
	"context"

	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/op"
	"gioui.org/text"
	"github.com/mleku/fromage"
	"lol.mleku.dev/chk"
)

func main() {
	th := fromage.NewTheme(context.Background(),
		fromage.NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())))
	_ = th
	w := fromage.NewWindow(th)
	w.Run(loop(w.Window))
}

func loop(w *app.Window) func() {
	return func() {
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				chk.E(e.Err)
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				e.Frame(gtx.Ops)
			}
		}

	}
}
