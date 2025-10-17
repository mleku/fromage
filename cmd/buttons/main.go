package main

import (
	"context"
	"fmt"
	"image"
	"image/color"

	"gio.mleku.dev/app"
	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/gesture"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/clip"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
	"gio.mleku.dev/widget"
	"gio.tools/icons"
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

// RightClickPopup represents a popup that appears on right-click
type RightClickPopup struct {
	visible     bool
	position    image.Point
	closeButton *fromage.ButtonLayout
	scrimClick  *widget.Clickable
	theme       *fromage.Theme
}

// Application state struct to hold persistent widgets
type AppState struct {
	switchWidget      *fromage.Bool
	colorSelector     *fromage.ColorSelector
	checkbox          *fromage.Checkbox
	modalStack        *fromage.ModalStack
	verticalRadio     *fromage.RadioButtonGroup
	horizontalRadio   *fromage.RadioButtonGroup
	intSlider         *fromage.Int
	rightClickGesture gesture.Click
	popup             *RightClickPopup
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

	w := fromage.NewWindow(th)

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
		verticalRadio: w.VerticalRadioGroup().
			AddButton("Option A", true).
			AddButton("Option B", false).
			AddButton("Option C", false).
			SetOnChange(func(index int, label string) {
				log.I.F("[HOOK] Vertical radio selected: %d - %s", index, label)
			}),
		horizontalRadio: w.HorizontalRadioGroup().
			AddButton("Red", false).
			AddButton("Green", true).
			AddButton("Blue", false).
			SetOnChange(func(index int, label string) {
				log.I.F("[HOOK] Horizontal radio selected: %d - %s", index, label)
			}),
		intSlider: th.NewInt().
			SetRange(0, 100).
			SetValue(50).
			SetHook(func(value int) {
				log.I.F("[HOOK] Int slider changed to: %d", value)
			}),
		popup: &RightClickPopup{
			visible:    false,
			theme:      th,
			scrimClick: &widget.Clickable{},
			closeButton: th.NewButtonLayout().
				Background(th.Colors.Error()).
				CornerRadius(0.5).
				Widget(func(g C) D {
					return th.Caption("×").
						Color(th.Colors.OnError()).
						Alignment(text.Middle).
						Layout(g)
				}),
		},
	}

	// Initialize the color selector with the current surface tint
	currentSurfaceTint := th.Colors.GetSurfaceTint()
	appState.colorSelector.SetColor(currentSurfaceTint)
	w.Option(app.Size(
		unit.Dp(1200), unit.Dp(1200)),
		app.Title("Kitchensink - Theme Demo"),
	)
	w.Run(loop(w.Window, th))
}

func loop(w *app.Window, th *fromage.Theme) func() {
	return func() {
		var ops op.Ops
		// Create a fromage window wrapper
		fromageWindow := &fromage.Window{Window: w, Theme: th}
		for {
			switch e := w.Event().(type) {
			case app.DestroyEvent:
				chk.E(e.Err)
				return
			case app.FrameEvent:
				gtx := app.NewContext(&ops, e)
				th.Pool.Reset() // Reset pool at the beginning of each frame
				mainUI(gtx, th, fromageWindow)
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

	// Push the modal to the stack
	appState.modalStack.Push(modalContent, func() {
		appState.modalStack.Pop()
	})
}

func mainUI(gtx layout.Context, th *fromage.Theme, w *fromage.Window) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())

	// Handle right-click gestures for showing popup
	for {
		ev, ok := appState.rightClickGesture.Update(gtx.Source)
		if !ok {
			break
		}

		if ev.Kind == gesture.KindClick {
			// Show popup at the click position
			clickPos := image.Pt(int(ev.Position.X), int(ev.Position.Y))
			appState.popup.ShowPopup(clickPos, gtx.Constraints.Max)
		}
	}

	// Register right-click gesture area for the entire screen
	area := image.Rectangle{Max: gtx.Constraints.Max}
	defer clip.Rect(area).Push(gtx.Ops).Pop()
	appState.rightClickGesture.Add(gtx.Ops)

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
			// Two neat rows of buttons at the top
			return th.VFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// First row: Main button styles
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
						Rigid(func(g C) D {
							// Disabled button example
							btn := th.PrimaryButton(func(g C) D {
								return th.Body2("Disabled").
									Color(th.Colors.OnPrimary()).
									Alignment(text.Middle).
									Layout(g)
							},
							).Disabled(true) // This button is disabled

							return btn.Layout(g)
						}).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Second row: Shape and style buttons
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
							// Icon-only button
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
											return th.Body2("Text").
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
						Layout(g)
				}).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Three-column layout: Radio buttons (vertical), Radio buttons (horizontal), Switch & Checkbox
			return th.HFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					// First column: Vertical radio buttons
					return w.Inset(0.5, func(g C) D {
						return th.VFlex().
							SpaceEvenly().
							Rigid(func(g C) D {
								return th.Caption("Radio Buttons - Vertical").
									Color(th.Colors.OnBackground()).
									Alignment(text.Middle).
									Layout(g)
							}).
							Rigid(func(g C) D {
								return appState.verticalRadio.Layout(g)
							}).
							Layout(g)
					}).Fn(g)
				}).
				Rigid(func(g C) D {
					// Second column: Horizontal radio buttons
					return w.Inset(0.5, func(g C) D {
						return th.VFlex().
							SpaceEvenly().
							Rigid(func(g C) D {
								return th.Caption("Radio Buttons - Horizontal").
									Color(th.Colors.OnBackground()).
									Alignment(text.Middle).
									Layout(g)
							}).
							Rigid(func(g C) D {
								return appState.horizontalRadio.Layout(g)
							}).
							Layout(g)
					}).Fn(g)
				}).
				Rigid(func(g C) D {
					// Third column: Switch and Checkbox
					return th.VFlex().
						SpaceEvenly().
						Rigid(func(g C) D {
							// Switch widget
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
							// Checkbox
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
						Layout(g)
				}).
				Layout(g)
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
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Int Slider showcase
			return th.VFlex().
				SpaceEvenly().
				Rigid(func(g C) D {
					return th.Caption("Integer Slider Example").
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Display current value
					currentValue := appState.intSlider.Value()
					return th.Body2(fmt.Sprintf("Current Value: %d", currentValue)).
						Color(th.Colors.OnBackground()).
						Alignment(text.Middle).
						Layout(g)
				}).
				Rigid(func(g C) D {
					// Int slider
					return appState.intSlider.Layout(g, th)
				}).
				Layout(g)
		}).
		Layout(gtx)

	// Layout the modal stack on top of everything
	if !appState.modalStack.IsEmpty() {
		appState.modalStack.Layout(gtx)
	}

	// Layout the popup on top of everything
	appState.popup.Layout(gtx)
}

// ShowPopup shows the popup at the specified position
func (p *RightClickPopup) ShowPopup(position image.Point, screenSize image.Point) {
	p.visible = true
	p.position = p.calculatePopupPosition(position, screenSize)
	log.I.F("Showing popup at position (%d, %d)", p.position.X, p.position.Y)
}

// HidePopup hides the popup
func (p *RightClickPopup) HidePopup() {
	p.visible = false
	log.I.F("Hiding popup")
}

// calculatePopupPosition calculates where to position the popup so the corner faces away from center
func (p *RightClickPopup) calculatePopupPosition(clickPos image.Point, screenSize image.Point) image.Point {
	centerX := screenSize.X / 2
	centerY := screenSize.Y / 2

	popupWidth := 200  // Approximate popup width
	popupHeight := 100 // Approximate popup height

	// Determine which corner should face away from center
	if clickPos.X < centerX {
		// Click is on left side, position popup to the right
		if clickPos.Y < centerY {
			// Click is in top-left, position popup bottom-right of click
			return image.Pt(clickPos.X, clickPos.Y)
		} else {
			// Click is in bottom-left, position popup top-right of click
			return image.Pt(clickPos.X, clickPos.Y-popupHeight)
		}
	} else {
		// Click is on right side, position popup to the left
		if clickPos.Y < centerY {
			// Click is in top-right, position popup bottom-left of click
			return image.Pt(clickPos.X-popupWidth, clickPos.Y)
		} else {
			// Click is in bottom-right, position popup top-left of click
			return image.Pt(clickPos.X-popupWidth, clickPos.Y-popupHeight)
		}
	}
}

// Layout renders the popup if it's visible
func (p *RightClickPopup) Layout(gtx C) D {
	if !p.visible {
		return D{}
	}

	// Handle scrim clicks
	if p.scrimClick.Clicked(gtx) {
		p.HidePopup()
		return D{}
	}

	// Handle close button clicks
	if p.closeButton.Clicked(gtx) {
		p.HidePopup()
		return D{}
	}

	// Create scrim (dimmed background)
	scrimColor := color.NRGBA{R: 0, G: 0, B: 0, A: 128} // 50% opacity black
	paint.Fill(gtx.Ops, scrimColor)

	// Layout scrim clickable area
	p.scrimClick.Layout(gtx, func(gtx C) D {
		return layout.Dimensions{Size: gtx.Constraints.Max}
	})

	// Position the popup
	offset := op.Offset(p.position).Push(gtx.Ops)
	defer offset.Pop()

	// Constrain popup size
	gtx.Constraints.Min.X = 200
	gtx.Constraints.Max.X = 200
	gtx.Constraints.Min.Y = 100
	gtx.Constraints.Max.Y = 100

	// Create popup background
	return p.theme.NewCard(
		func(g C) D {
			return p.theme.VFlex().
				Rigid(func(gtx C) D {
					// Title
					return p.theme.Body2("Right-click Popup").
						Color(p.theme.Colors.OnSurface()).
						Alignment(text.Middle).
						Layout(gtx)
				}).
				Rigid(func(gtx C) D {
					// Content
					return p.theme.Caption("This popup appeared because you right-clicked!").
						Color(p.theme.Colors.OnSurfaceVariant()).
						Alignment(text.Middle).
						Layout(gtx)
				}).
				Rigid(func(gtx C) D {
					// Close button
					return p.closeButton.Layout(gtx)
				}).
				Layout(g)
		},
	).CornerRadius(8).Padding(unit.Dp(12)).Layout(gtx)
}
