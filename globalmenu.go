package fromage

import (
	"image"
	"image/color"
	"time"

	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
	"gio.mleku.dev/op/paint"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
	"gio.mleku.dev/widget"
)

// GlobalMenu represents a right-click context menu system
type GlobalMenu struct {
	theme        *Theme
	eventHandler *EventHandler
	items        []*MenuItem
	nextItemID   int
	animating    bool
	visible      bool
	scrimVisible bool
	position     image.Point
	clickPos     image.Point // Store the right-click position
	// Animation state
	showTime     time.Time
	hideTime     time.Time
	isHiding     bool
	shouldRemove bool
}

// MenuItem represents a single item in the context menu
type MenuItem struct {
	id          int
	text        string
	action      func()
	closeButton *ButtonLayout
	clickable   *widget.Clickable
	visible     bool
	// Animation state
	showTime     time.Time
	hideTime     time.Time
	isHiding     bool
	shouldRemove bool
}

// NewGlobalMenu creates a new global menu instance
func (t *Theme) NewGlobalMenu() *GlobalMenu {
	gm := &GlobalMenu{
		theme:        t,
		items:        make([]*MenuItem, 0),
		nextItemID:   1,
		animating:    false,
		visible:      false,
		scrimVisible: false,
	}

	// Create event handler with callbacks
	gm.eventHandler = NewEventHandler(func(event string) {
		// Log menu events if needed
	}).SetOnPress(func(e pointer.Event) {
		if e.Buttons == pointer.ButtonSecondary {
			// Right-click detected
			clickPos := image.Pt(int(e.Position.X), int(e.Position.Y))
			gm.Show(clickPos, image.Pt(800, 600)) // Default viewport size
		} else if e.Buttons == pointer.ButtonPrimary && gm.scrimVisible {
			// Left-click on scrim to close menu
			clickPos := image.Pt(int(e.Position.X), int(e.Position.Y))
			menuWidth := 250
			menuHeight := len(gm.items) * 40
			menuRect := image.Rect(
				gm.position.X,
				gm.position.Y,
				gm.position.X+menuWidth,
				gm.position.Y+menuHeight,
			)

			if !clickPos.In(menuRect) {
				// Click outside menu - hide it
				gm.Hide()
			}
		}
	})

	return gm
}

// AddItem adds a new item to the menu
func (gm *GlobalMenu) AddItem(itemText string, action func()) *GlobalMenu {
	item := &MenuItem{
		id:     gm.nextItemID,
		text:   itemText,
		action: action,
		closeButton: gm.theme.NewButtonLayout().
			Background(gm.theme.Colors.Error()).
			CornerRadius(0.5).
			Widget(func(g C) D {
				return gm.theme.Caption("×").
					Color(gm.theme.Colors.OnError()).
					Alignment(text.Middle).
					Layout(g)
			}),
		clickable:    &widget.Clickable{},
		visible:      true,
		showTime:     time.Now(),
		isHiding:     false,
		shouldRemove: false,
	}

	gm.items = append(gm.items, item)
	gm.nextItemID++
	return gm
}

// Show displays the menu at the specified position
func (gm *GlobalMenu) Show(position image.Point, viewportSize image.Point) {
	gm.clickPos = position // Store the click position
	gm.position = gm.calculateSmartPosition(position, viewportSize)
	gm.visible = true
	gm.scrimVisible = true
	gm.showTime = time.Now()
	gm.isHiding = false
	gm.shouldRemove = false
	gm.animating = true
}

// Hide starts the hide animation for the menu
func (gm *GlobalMenu) Hide() {
	gm.hideTime = time.Now()
	gm.isHiding = true
	gm.scrimVisible = false
}

// IsVisible returns whether the menu is currently visible
func (gm *GlobalMenu) IsVisible() bool {
	return gm.visible && !gm.shouldRemove
}

// GetClickPosition returns the position where the menu was triggered
func (gm *GlobalMenu) GetClickPosition() image.Point {
	return gm.clickPos
}

// Layout renders the global menu
func (gm *GlobalMenu) Layout(gtx layout.Context) {
	if !gm.visible {
		return
	}

	now := time.Now()

	// Calculate animation progress
	var alpha float32 = 1.0
	if !gm.isHiding {
		// Fade in animation
		elapsed := now.Sub(gm.showTime)
		if elapsed < 250*time.Millisecond {
			alpha = float32(elapsed) / float32(250*time.Millisecond)
		}
	} else {
		// Fade out animation
		elapsed := now.Sub(gm.hideTime)
		if elapsed < 250*time.Millisecond {
			alpha = 1.0 - (float32(elapsed) / float32(250*time.Millisecond))
		} else {
			// Animation complete, mark for removal
			gm.shouldRemove = true
			gm.visible = false
			return
		}
	}

	// Handle scrim clicks
	if gm.scrimVisible {
		// Create scrim with animation alpha
		scrimAlpha := uint8(128 * alpha) // 50% opacity * animation alpha
		scrimColor := color.NRGBA{R: 0, G: 0, B: 0, A: scrimAlpha}
		paint.Fill(gtx.Ops, scrimColor)
	}

	// Register event handler for this menu
	gm.eventHandler.AddToOps(gtx.Ops)
	gm.eventHandler.ProcessEvents(gtx)

	// Position the menu
	offset := op.Offset(gm.position).Push(gtx.Ops)
	defer offset.Pop()

	// Layout menu items
	menuWidth := 250                 // Increased width to accommodate longer text
	menuHeight := len(gm.items) * 40 // 40px per item

	// Constrain menu size
	gtx.Constraints.Min.X = menuWidth
	gtx.Constraints.Max.X = menuWidth
	gtx.Constraints.Min.Y = menuHeight
	gtx.Constraints.Max.Y = menuHeight

	// Apply opacity to the entire menu
	opacity := float32(uint8(255*alpha)) / 255.0
	opacityStack := paint.PushOpacity(gtx.Ops, opacity)
	defer opacityStack.Pop()

	// Create menu background
	gm.theme.NewCard(
		func(g C) D {
			return gm.theme.VFlex().
				Rigid(func(gtx C) D {
					// Header
					return gm.theme.HFlex().
						Flexed(1, func(gtx C) D {
							return gm.theme.Caption("Context Menu").
								Color(gm.theme.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(gtx)
						}).
						Rigid(func(gtx C) D {
							closeBtn := gm.theme.NewButtonLayout().
								Background(gm.theme.Colors.Error()).
								CornerRadius(0.5).
								Widget(func(g C) D {
									return gm.theme.Caption("×").
										Color(gm.theme.Colors.OnError()).
										Alignment(text.Middle).
										Layout(g)
								})
							if closeBtn.Clicked(gtx) {
								gm.Hide()
							}
							return closeBtn.Layout(gtx)
						}).
						Layout(gtx)
				}).
				Rigid(func(gtx C) D {
					// Menu items with proper constraints
					menuWidth := 250 // Match the increased width
					gtx.Constraints.Min.X = menuWidth
					gtx.Constraints.Max.X = menuWidth

					flex := gm.theme.VFlex()
					for _, item := range gm.items {
						if item.visible && !item.shouldRemove {
							flex = flex.Rigid(func(gtx C) D {
								// Constrain each menu item to fill the available width and have minimum height
								gtx.Constraints.Min.Y = 40
								gtx.Constraints.Max.Y = 40
								return item.Layout(gtx, gm.theme)
							})
						}
					}
					return flex.Layout(gtx)
				}).
				Layout(g)
		},
	).CornerRadius(8).Padding(unit.Dp(8)).Layout(gtx)

	// Update animation state
	gm.animating = alpha < 1.0
}

// calculateSmartPosition calculates where to position the menu so it faces toward the center
func (gm *GlobalMenu) calculateSmartPosition(clickPos image.Point, viewportSize image.Point) image.Point {
	centerX := viewportSize.X / 2
	centerY := viewportSize.Y / 2

	menuWidth := 250 // Match the increased width
	menuHeight := len(gm.items) * 40

	// Determine which corner should face toward the center
	if clickPos.X < centerX {
		// Click is on left side of center
		if clickPos.Y < centerY {
			// Click is in top-left quadrant, position menu bottom-right of click
			return image.Pt(clickPos.X, clickPos.Y)
		} else {
			// Click is in bottom-left quadrant, position menu top-right of click
			return image.Pt(clickPos.X, clickPos.Y-menuHeight)
		}
	} else {
		// Click is on right side of center
		if clickPos.Y < centerY {
			// Click is in top-right quadrant, position menu bottom-left of click
			return image.Pt(clickPos.X-menuWidth, clickPos.Y)
		} else {
			// Click is in bottom-right quadrant, position menu top-left of click
			return image.Pt(clickPos.X-menuWidth, clickPos.Y-menuHeight)
		}
	}
}

// Layout renders a menu item
func (item *MenuItem) Layout(gtx layout.Context, th *Theme) layout.Dimensions {
	now := time.Now()

	// Calculate animation progress
	var alpha float32 = 1.0
	if !item.isHiding {
		// Fade in animation
		elapsed := now.Sub(item.showTime)
		if elapsed < 250*time.Millisecond {
			alpha = float32(elapsed) / float32(250*time.Millisecond)
		}
	} else {
		// Fade out animation
		elapsed := now.Sub(item.hideTime)
		if elapsed < 250*time.Millisecond {
			alpha = 1.0 - (float32(elapsed) / float32(250*time.Millisecond))
		} else {
			// Animation complete, mark for removal
			item.shouldRemove = true
			item.visible = false
			return layout.Dimensions{}
		}
	}

	// Handle item clicks
	if item.clickable.Clicked(gtx) {
		if item.action != nil {
			item.action()
		}
	}

	// Apply opacity
	opacity := float32(uint8(255*alpha)) / 255.0
	opacityStack := paint.PushOpacity(gtx.Ops, opacity)
	defer opacityStack.Pop()

	// Create item background - let flex container manage sizing
	return th.NewButtonLayout().
		Background(th.Colors.SurfaceVariant()).
		CornerRadius(4).
		Widget(func(g C) D {
			// Constrain text to fit within the menu item
			g.Constraints.Min.X = 0
			return th.Body2(item.text).
				Color(th.Colors.OnSurfaceVariant()).
				Alignment(text.Middle).
				Layout(g)
		}).
		Layout(gtx)
}
