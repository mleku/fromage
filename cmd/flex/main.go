package main

import (
	"context"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
	"github.com/mleku/fromage"
	"lol.mleku.dev/chk"
)

// Import aliases from fromage package
type (
	C = fromage.C
	D = fromage.D
	W = fromage.W
)

func main() {
	th := fromage.NewThemeWithMode(
		context.Background(),
		fromage.NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		fromage.ThemeModeDark,
	)
	w := fromage.NewWindow(th)
	w.Option(app.Size(
		unit.Dp(640), unit.Dp(1280)),
		app.Title("Flex Demo"),
	)
	w.Run(loop(w.Window, th))
}

func loop(w *app.Window, th *fromage.Theme) func() {
	return func() {
		var ops op.Ops
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				chk.E(e.Err)
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				th.Pool.Reset() // Reset pool at the beginning of each frame
				mainUI(gtx, th)
				e.Frame(gtx.Ops)
			}
		}
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())

	// Main horizontal layout with two vertical flex boxes
	th.VFlex().
		SpaceEvenly().
		Flexed(1, func(g C) D {
			// Left flex box - vertical flex with borders
			return th.BorderOutline().Widget(func(g C) D {
				return th.FillSurface(func(g C) D {
					return th.VFlex().
						Flexed(1, func(g C) D {
							return th.BorderPrimary().Widget(func(g C) D {
								return th.FillPrimary(func(g C) D {
									return th.HFlex().
										SpaceAround().
										Flexed(1, func(g C) D {
											return th.BorderPrimary().Widget(func(g C) D {
												return th.FillPrimary(func(g C) D {
													return th.Caption("Box 1").Color(th.Colors.OnPrimary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSecondary().Widget(func(g C) D {
												return th.FillSecondary(func(g C) D {
													return th.Caption("Box 2").Color(th.Colors.OnSecondary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Rigid(func(g C) D {
											return th.BorderSurface().Widget(func(g C) D {
												return th.NewFill(th.Colors.Tertiary(), func(g C) D {
													return th.Caption("Box 3").Color(th.Colors.OnTertiary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderError().Widget(func(g C) D {
												return th.FillError(func(g C) D {
													return th.Caption("Box 4").Color(th.Colors.OnError()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderOutline().Widget(func(g C) D {
												return th.FillBackground(func(g C) D {
													return th.Caption("Box 5").Color(th.Colors.OnBackground()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Layout(g)

								}).Layout(g)
							}).Layout(g)
						}).
						Flexed(1, func(g C) D {
							return th.BorderSecondary().Widget(func(g C) D {
								return th.FillSecondary(func(g C) D {
									return th.HFlex().
										SpaceBetween().
										Flexed(1, func(g C) D {
											return th.BorderPrimary().Widget(func(g C) D {
												return th.FillPrimary(func(g C) D {
													return th.Caption("Box 1").Color(th.Colors.OnPrimary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSecondary().Widget(func(g C) D {
												return th.FillSecondary(func(g C) D {
													return th.Caption("Box 2").Color(th.Colors.OnSecondary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Rigid(func(g C) D {
											return th.BorderSurface().Widget(func(g C) D {
												return th.NewFill(th.Colors.Tertiary(), func(g C) D {
													return th.Caption("Box 3").Color(th.Colors.OnTertiary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderError().Widget(func(g C) D {
												return th.FillError(func(g C) D {
													return th.Caption("Box 4").Color(th.Colors.OnError()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderOutline().Widget(func(g C) D {
												return th.FillBackground(func(g C) D {
													return th.Caption("Box 5").Color(th.Colors.OnBackground()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Layout(g)

								}).Layout(g)
							}).Layout(g)
						}).
						Flexed(1, func(g C) D {
							return th.BorderSurface().Widget(func(g C) D {
								return th.NewFill(th.Colors.Tertiary(), func(g C) D {
									return th.HFlex().
										SpaceEnd().
										Flexed(1, func(g C) D {
											return th.BorderPrimary().Widget(func(g C) D {
												return th.FillPrimary(func(g C) D {
													return th.Caption("Box 1").Color(th.Colors.OnPrimary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSecondary().Widget(func(g C) D {
												return th.FillSecondary(func(g C) D {
													return th.Caption("Box 2").Color(th.Colors.OnSecondary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Rigid(func(g C) D {
											return th.BorderSurface().Widget(func(g C) D {
												return th.NewFill(th.Colors.Tertiary(), func(g C) D {
													return th.Caption("Box 3").Color(th.Colors.OnTertiary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderError().Widget(func(g C) D {
												return th.FillError(func(g C) D {
													return th.Caption("Box 4").Color(th.Colors.OnError()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderOutline().Widget(func(g C) D {
												return th.FillBackground(func(g C) D {
													return th.Caption("Box 5").Color(th.Colors.OnBackground()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Layout(g)

								}).Layout(g)
							}).Layout(g)
						}).
						Flexed(1, func(g C) D {
							return th.BorderError().Widget(func(g C) D {
								return th.FillError(func(g C) D {
									return th.HFlex().
										SpaceEvenly().
										Flexed(1, func(g C) D {
											return th.BorderPrimary().Widget(func(g C) D {
												return th.FillPrimary(func(g C) D {
													return th.Caption("Box 1").Color(th.Colors.OnPrimary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSecondary().Widget(func(g C) D {
												return th.FillSecondary(func(g C) D {
													return th.Caption("Box 2").Color(th.Colors.OnSecondary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Rigid(func(g C) D {
											return th.BorderSurface().Widget(func(g C) D {
												return th.NewFill(th.Colors.Tertiary(), func(g C) D {
													return th.Caption("Box 3").Color(th.Colors.OnTertiary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderError().Widget(func(g C) D {
												return th.FillError(func(g C) D {
													return th.Caption("Box 4").Color(th.Colors.OnError()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderOutline().Widget(func(g C) D {
												return th.FillBackground(func(g C) D {
													return th.Caption("Box 5").Color(th.Colors.OnBackground()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Layout(g)

								}).Layout(g)
							}).Layout(g)
						}).
						Flexed(1, func(g C) D {
							return th.BorderOutline().Widget(func(g C) D {
								return th.FillBackground(func(g C) D {
									return th.HFlex().SpaceStart().
										Flexed(1, func(g C) D {
											return th.BorderPrimary().Widget(func(g C) D {
												return th.FillPrimary(func(g C) D {
													return th.Caption("Box 1").Color(th.Colors.OnPrimary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSecondary().Widget(func(g C) D {
												return th.FillSecondary(func(g C) D {
													return th.Caption("Box 2").Color(th.Colors.OnSecondary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Rigid(func(g C) D {
											return th.BorderSurface().Widget(func(g C) D {
												return th.NewFill(th.Colors.Tertiary(), func(g C) D {
													return th.Caption("Box 3").Color(th.Colors.OnTertiary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderError().Widget(func(g C) D {
												return th.FillError(func(g C) D {
													return th.Caption("Box 4").Color(th.Colors.OnError()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderOutline().Widget(func(g C) D {
												return th.FillBackground(func(g C) D {
													return th.Caption("Box 5").Color(th.Colors.OnBackground()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Layout(g)

								}).Layout(g)
							}).Layout(g)
						}).
						Layout(g)
				}).CornerRadius(8).Layout(g)
			}).Layout(g)
		}).
		Flexed(1, func(g C) D {
			// Right flex box - vertical flex with borders
			return th.BorderOutline().Widget(func(g C) D {
				return th.FillSurface(func(g C) D {
					return th.HFlex().
						Flexed(1, func(g C) D {
							return th.BorderPrimary().Widget(func(g C) D {
								return th.FillPrimary(func(g C) D {
									return th.VFlex().
										Flexed(1, func(g C) D {
											return th.BorderPrimary().Widget(func(g C) D {
												return th.FillPrimary(func(g C) D {
													return th.Caption("Box 1").Color(th.Colors.OnPrimary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSecondary().Widget(func(g C) D {
												return th.FillSecondary(func(g C) D {
													return th.Caption("Box 2").Color(th.Colors.OnSecondary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSurface().Widget(func(g C) D {
												return th.NewFill(th.Colors.Tertiary(), func(g C) D {
													return th.Caption("Box 3").Color(th.Colors.OnTertiary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderError().Widget(func(g C) D {
												return th.FillError(func(g C) D {
													return th.Caption("Box 4").Color(th.Colors.OnError()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderOutline().Widget(func(g C) D {
												return th.FillBackground(func(g C) D {
													return th.Caption("Box 5").Color(th.Colors.OnBackground()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Layout(g)
								}).Layout(g)
							}).Layout(g)
						}).
						Flexed(1, func(g C) D {
							return th.BorderSecondary().Widget(func(g C) D {
								return th.FillSecondary(func(g C) D {
									return th.VFlex().
										Flexed(1, func(g C) D {
											return th.BorderPrimary().Widget(func(g C) D {
												return th.FillPrimary(func(g C) D {
													return th.Caption("Box 1").Color(th.Colors.OnPrimary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSecondary().Widget(func(g C) D {
												return th.FillSecondary(func(g C) D {
													return th.Caption("Box 2").Color(th.Colors.OnSecondary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSurface().Widget(func(g C) D {
												return th.NewFill(th.Colors.Tertiary(), func(g C) D {
													return th.Caption("Box 3").Color(th.Colors.OnTertiary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderError().Widget(func(g C) D {
												return th.FillError(func(g C) D {
													return th.Caption("Box 4").Color(th.Colors.OnError()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderOutline().Widget(func(g C) D {
												return th.FillBackground(func(g C) D {
													return th.Caption("Box 5").Color(th.Colors.OnBackground()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Layout(g)
								}).Layout(g)
							}).Layout(g)
						}).
						Flexed(1, func(g C) D {
							return th.BorderSurface().Widget(func(g C) D {
								return th.NewFill(th.Colors.Tertiary(), func(g C) D {
									return th.VFlex().
										Flexed(1, func(g C) D {
											return th.BorderPrimary().Widget(func(g C) D {
												return th.FillPrimary(func(g C) D {
													return th.Caption("Box 1").Color(th.Colors.OnPrimary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSecondary().Widget(func(g C) D {
												return th.FillSecondary(func(g C) D {
													return th.Caption("Box 2").Color(th.Colors.OnSecondary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSurface().Widget(func(g C) D {
												return th.NewFill(th.Colors.Tertiary(), func(g C) D {
													return th.Caption("Box 3").Color(th.Colors.OnTertiary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderError().Widget(func(g C) D {
												return th.FillError(func(g C) D {
													return th.Caption("Box 4").Color(th.Colors.OnError()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderOutline().Widget(func(g C) D {
												return th.FillBackground(func(g C) D {
													return th.Caption("Box 5").Color(th.Colors.OnBackground()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Layout(g)
								}).Layout(g)
							}).Layout(g)
						}).
						Flexed(1, func(g C) D {
							return th.BorderError().Widget(func(g C) D {
								return th.FillError(func(g C) D {
									return th.VFlex().
										Flexed(1, func(g C) D {
											return th.BorderPrimary().Widget(func(g C) D {
												return th.FillPrimary(func(g C) D {
													return th.Caption("Box 1").Color(th.Colors.OnPrimary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSecondary().Widget(func(g C) D {
												return th.FillSecondary(func(g C) D {
													return th.Caption("Box 2").Color(th.Colors.OnSecondary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSurface().Widget(func(g C) D {
												return th.NewFill(th.Colors.Tertiary(), func(g C) D {
													return th.Caption("Box 3").Color(th.Colors.OnTertiary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderError().Widget(func(g C) D {
												return th.FillError(func(g C) D {
													return th.Caption("Box 4").Color(th.Colors.OnError()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderOutline().Widget(func(g C) D {
												return th.FillBackground(func(g C) D {
													return th.Caption("Box 5").Color(th.Colors.OnBackground()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Layout(g)
								}).Layout(g)
							}).Layout(g)
						}).
						Flexed(1, func(g C) D {
							return th.BorderOutline().Widget(func(g C) D {
								return th.FillBackground(func(g C) D {
									return th.VFlex().
										Flexed(1, func(g C) D {
											return th.BorderPrimary().Widget(func(g C) D {
												return th.FillPrimary(func(g C) D {
													return th.Caption("Box 1").Color(th.Colors.OnPrimary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSecondary().Widget(func(g C) D {
												return th.FillSecondary(func(g C) D {
													return th.Caption("Box 2").Color(th.Colors.OnSecondary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderSurface().Widget(func(g C) D {
												return th.NewFill(th.Colors.Tertiary(), func(g C) D {
													return th.Caption("Box 3").Color(th.Colors.OnTertiary()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderError().Widget(func(g C) D {
												return th.FillError(func(g C) D {
													return th.Caption("Box 4").Color(th.Colors.OnError()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Flexed(1, func(g C) D {
											return th.BorderOutline().Widget(func(g C) D {
												return th.FillBackground(func(g C) D {
													return th.Caption("Box 5").Color(th.Colors.OnBackground()).Alignment(text.Middle).Layout(g)
												}).Layout(g)
											}).Layout(g)
										}).
										Layout(g)
								}).Layout(g)
							}).Layout(g)
						}).
						Layout(g)
				}).CornerRadius(8).Layout(g)
			}).Layout(g)
		}).
		Layout(gtx)
}
