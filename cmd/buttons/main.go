package main

import (
	"context"
	"fmt"
	"image/color"

	"gio.tools/icons"
	"gioui.org/app"
	"gioui.org/font/gofont"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/mleku/fromage"
	"lol.mleku.dev/chk"
	"lol.mleku.dev/log"
)

// Import aliases from fromage package
type (
	C = fromage.C
	D = fromage.D
	W = fromage.W
)

// Icon instances using gio.tools/icons
var (
	starIcon     = icons.ActionGrade
	heartIcon    = icons.ActionFavorite
	settingsIcon = icons.ActionSettings
)

// Application state struct to hold persistent widgets
type AppState struct {
	switchWidget  *fromage.Bool
	colorSelector *fromage.ColorSelector
	checkbox      *fromage.Checkbox
	modalStack    *fromage.ModalStack
	slideInCheck  *fromage.Checkbox
}

var appState *AppState

func main() {
	th := fromage.NewThemeWithMode(
		context.Background(),
		fromage.NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		fromage.ThemeModeDark,
	)

	// Initialize application state with persistent widgets
	appState = &AppState{
		switchWidget: th.Switch(false).SetOnChange(func(b bool) {
			log.I.F("[HOOK] Switch toggled to: %v", b)
		}),
		colorSelector: th.NewColorSelector().SetOnChange(func(c color.NRGBA) {
			log.I.F("[HOOK] Color changed to: R=%d G=%d B=%d", c.R, c.G, c.B)
			// Update surface tint
			th.Colors.SetSurfaceTint(c)
		}),
		checkbox: th.NewCheckbox(false).SetOnChange(func(b bool) {
			log.I.F("[HOOK] Checkbox toggled to: %v", b)
		}),
		modalStack: th.NewModalStack().ScrimDarkness(0.7), // 70% opacity scrim
		slideInCheck: th.NewCheckbox(false).SetOnChange(func(b bool) {
			log.I.F("[HOOK] Slide-in checkbox toggled to: %v", b)
		}),
	}

	// Initialize the color selector with the current surface tint
	currentSurfaceTint := th.Colors.GetSurfaceTint()
	appState.colorSelector.SetColor(currentSurfaceTint)

	w := fromage.NewWindow(th)
	w.Option(app.Size(
		unit.Dp(800), unit.Dp(800)),
		app.Title("Kitchensink - Theme Demo"),
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

// showModal creates and displays a modal with generated text content
func showModal(th *fromage.Theme) {
	// Generate some sample text content
	sampleText := `This is a modal dialog with customizable scrim darkness.

You can put any widgets you want in here - buttons, text, images, forms, etc.

The scrim behind this modal is set to 70% opacity, making the background content visible but dimmed.

Click anywhere outside this content area to close the modal.

This modal demonstrates:
• Customizable scrim darkness
• Click-outside-to-close functionality  
• Centered content layout
• Full-screen overlay`

	// Create the modal content
	modalContent := func(g C) D {
		return th.NewCard(
			func(g C) D {
				return th.VFlex().
					SpaceEvenly().
					Rigid(func(g C) D {
						return th.H3("Modal Dialog").
							Color(th.Colors.OnSurface()).
							Alignment(text.Middle).
							Layout(g)
					}).
					Rigid(func(g C) D {
						return th.Body1(sampleText).
							Color(th.Colors.OnSurface()).
							Alignment(text.Start).
							Layout(g)
					}).
					Rigid(func(g C) D {
						// Close button
						btn := th.SecondaryButton(func(g C) D {
							return th.Body2("Close").
								Color(th.Colors.OnSecondary()).
								Alignment(text.Middle).
								Layout(g)
						})
						if btn.Clicked(g) {
							appState.modalStack.Pop()
						}
						return btn.Layout(g)
					}).
					Layout(g)
			},
		).CornerRadius(8).Padding(unit.Dp(16)).Layout(g)
	}

	// Push the modal to the stack with slide-in option
	slideIn := appState.slideInCheck.GetValue()
	appState.modalStack.PushWithSlide(modalContent, func() {
		appState.modalStack.Pop()
	}, slideIn)
}

func mainUI(gtx layout.Context, th *fromage.Theme) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())

	th.CenteredColumn().
		Rigid(func(g C) D {
			// Title with primary color fill
			return th.H3("Interactive Button Demo").Alignment(text.Middle).Layout(g)
		}).
		Rigid(func(g C) D {
			themeText := "Current Theme: Light"
			if th.IsDark() {
				themeText = "Current Theme: Dark"
			}
			// Theme info with surface color fill
			return th.FillSurface(
				func(g C) D {
					return th.Body1(themeText).Alignment(text.Middle).Layout(g)
				},
			).CornerRadius(4).Layout(g)
		}).
		Rigid(func(g C) D {
			// Main interactive button with theme toggle and icon
			button := th.PrimaryButton(func(g C) D {
				return th.HFlex().
					SpaceEvenly().
					AlignMiddle().
					Rigid(func(g C) D {
						return settingsIcon.Layout(g, th.Colors.Primary())
					}).
					Rigid(func(g C) D {
						return th.Body1("Toggle Theme").
							Color(th.Colors.OnPrimary()).
							Alignment(text.Middle).
							Layout(g)
					}).
					Layout(g)
			},
			)

			// Check for clicks BEFORE layout (this is the key fix!)
			if button.Clicked(g) {
				log.I.F("Toggle theme button clicked")
				th.ToggleTheme()
				// Update switch widget colors to match new theme
				appState.switchWidget.UpdateThemeColors(g.Now)
			}

			return button.Layout(g)
		}).
		Rigid(func(g C) D {
			// Button style showcase with hover effects
			return th.HFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// Secondary button with star icon
					btn := th.SecondaryButton(
						func(g C) D {
							return th.HFlex().
								SpaceEvenly().
								AlignMiddle().
								Rigid(func(g C) D {
									return starIcon.Layout(g, th.Colors.OnSecondary())
								}).
								Rigid(func(g C) D {
									return th.Body2("Secondary").
										Color(th.Colors.OnSecondary()).
										Alignment(text.Middle).
										Layout(g)
								}).
								Layout(g)
						},
					)
					if btn.Clicked(g) {
						log.I.F("Secondary button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Surface button with heart icon
					btn := th.SurfaceButton(func(g C) D {
						return th.HFlex().
							SpaceEvenly().
							AlignMiddle().
							Rigid(func(g C) D {
								return heartIcon.Layout(g, th.Colors.OnSurface())
							}).
							Rigid(func(g C) D {
								return th.Body2("Surface").
									Color(th.Colors.OnSurface()).
									Alignment(text.Middle).
									Layout(g)
							}).
							Layout(g)
					},
					)
					if btn.Clicked(g) {
						log.I.F("Surface button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Error button with warning icon
					btn := th.ErrorButton(func(g C) D {
						return th.HFlex().
							SpaceEvenly().
							AlignMiddle().
							Rigid(func(g C) D {
								return settingsIcon.Layout(g, th.Colors.OnError())
							}).
							Rigid(func(g C) D {
								return th.Body2("Error").
									Color(th.Colors.OnError()).
									Alignment(text.Middle).
									Layout(g)
							}).
							Layout(g)
					},
					)
					if btn.Clicked(g) {
						log.I.F("Error button clicked")
					}
					return btn.Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Shape showcase with different corner styles
			return th.HFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// Rounded button
					btn := th.RoundedButton(
						func(g C) D {
							return th.Caption("Rounded").
								Color(th.Colors.OnPrimary()).
								Alignment(text.Middle).
								Layout(g)
						},
					)
					if btn.Clicked(g) {
						log.I.F("rounded button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Pill button
					btn := th.PillButton(func(g C) D {
						return th.Caption("Pill Shape").
							Color(th.Colors.OnPrimary()).
							Alignment(text.Middle).
							Layout(g)
					},
					)
					if btn.Clicked(g) {
						log.I.F("pill button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Icon-only button using the icon widget
					btn := th.NewButtonLayout().
						Background(th.Colors.Tertiary()).
						CornerRadius(0.5). // 50% of text size
						Widget(func(g C) D {
							return starIcon.Layout(g, th.Colors.OnTertiary())
						})
					if btn.Clicked(g) {
						log.I.F("icon-only button clicked")
					}
					return btn.Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Icon size showcase
			return th.HFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// Small icon button
					btn := th.NewButtonLayout().
						Background(th.Colors.Primary()).
						CornerRadius(0.25). // 25% of text size
						Widget(func(g C) D {
							return starIcon.Layout(g, th.Colors.OnPrimary())
						})
					if btn.Clicked(g) {
						log.I.F("small icon button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Medium icon button
					btn := th.NewButtonLayout().
						Background(th.Colors.Secondary()).
						CornerRadius(0.25). // 25% of text size
						Widget(func(g C) D {
							return heartIcon.Layout(g, th.Colors.OnSecondary())
						})
					if btn.Clicked(g) {
						log.I.F("medium icon button clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Large icon button
					btn := th.NewButtonLayout().
						Background(th.Colors.Tertiary()).
						CornerRadius(0.25). // 25% of text size
						Widget(func(g C) D {
							return settingsIcon.Layout(g, th.Colors.OnTertiary())
						})
					if btn.Clicked(g) {
						log.I.F("large icon button clicked")
					}
					return btn.Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Text buttons with icons showcase
			return th.HFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// Text button with icon
					btn := th.NewButtonLayout().
						Widget(func(g C) D {
							return th.HFlex().
								SpaceEvenly().
								AlignMiddle().
								Rigid(func(g C) D {
									return starIcon.Layout(g, th.Colors.OnBackground())
								}).
								Rigid(func(g C) D {
									return th.Body2("Text Button").
										Color(th.Colors.OnBackground()).
										Alignment(text.Middle).
										Layout(g)
								}).
								Layout(g)
						})
					if btn.Clicked(g) {
						log.I.F("text button with icon clicked")
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Custom styled button with icon
					btn := th.NewButtonLayout().
						Background(th.Colors.Tertiary()).
						CornerRadius(0.3). // 30% of text size
						Corners(fromage.CornerNW | fromage.CornerNE).
						Widget(func(g C) D {
							return th.HFlex().
								SpaceEvenly().
								AlignMiddle().
								Rigid(func(g C) D {
									return heartIcon.Layout(g, th.Colors.OnTertiary())
								}).
								Rigid(func(g C) D {
									return th.Caption("Custom Style").
										Color(th.Colors.OnTertiary()).
										Alignment(text.Middle).
										Layout(g)
								}).
								Layout(g)
						})
					if btn.Clicked(g) {
						log.I.F("custom button with icon clicked")
					}
					return btn.Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Disabled button example
			btn := th.PrimaryButton(func(g C) D {
				return th.Body2("Disabled Button").
					Color(th.Colors.OnPrimary()).
					Alignment(text.Middle).
					Layout(g)
			},
			).Disabled(true) // This button is disabled

			return btn.Layout(g)
		}).
		Rigid(func(g C) D {
			// Color selector for surface tint
			return th.VFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					return th.Caption("Surface Tint Color").
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).
				Rigid(func(g C) D {
					return appState.colorSelector.Layout(g, th)
				}).
				Rigid(func(g C) D {
					// Display current surface tint value
					currentTint := th.Colors.GetSurfaceTint()
					return th.Caption(fmt.Sprintf("Current: %s", fromage.ColorToHex(currentTint))).
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Switch widget showcase
			return th.VFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// Let the bool widget handle its own clicks
					return appState.switchWidget.Layout(g)
				}).
				Rigid(func(g C) D {
					return th.Caption("Switch").
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Checkbox showcase
			return th.VFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					return th.Caption("Checkbox Example").
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Single checkbox
					checkbox := appState.checkbox.Label("Enable Feature")
					return checkbox.Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Modal showcase
			return th.VFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					return th.Caption("Modal Example").
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Button to show modal
					btn := th.PrimaryButton(func(g C) D {
						return th.Body2("Show Modal").
							Color(th.Colors.OnPrimary()).
							Alignment(text.Middle).
							Layout(g)
					})
					if btn.Clicked(g) {
						log.I.F("Show modal button clicked")
						showModal(th)
					}
					return btn.Layout(g)
				}).
				Rigid(func(g C) D {
					// Slide-in checkbox
					slideInCheckbox := appState.slideInCheck.Label("Slide in")
					return slideInCheckbox.Layout(g)
				}).
				Layout(g)
		}).
		Layout(gtx)

	// Layout the modal stack on top of everything
	if !appState.modalStack.IsEmpty() {
		appState.modalStack.Layout(gtx)
	}
}
