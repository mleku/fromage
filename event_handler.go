package fromage

import (
	"fmt"

	"lol.mleku.dev/log"

	"gio.mleku.dev/gesture"
	"gio.mleku.dev/io/event"
	"gio.mleku.dev/io/pointer"
	"gio.mleku.dev/layout"
	"gio.mleku.dev/op"
)

// EventHandler provides consistent event processing for pointer and scroll events
type EventHandler struct {
	// Track button state to show which button was released
	pressedButtons pointer.Buttons
	// Gesture scroll to capture scroll events
	scroll gesture.Scroll
	// Event logging function
	logEvent func(string)
	// Callback functions for different event types
	onClick   func(pointer.Event)
	onScroll  func(float32)
	onHover   func(bool)
	onDrag    func(pointer.Event)
	onPress   func(pointer.Event)
	onRelease func(pointer.Event)
}

// NewEventHandler creates a new event handler with the given event logging function
func NewEventHandler(logEvent func(string)) *EventHandler {
	return &EventHandler{
		logEvent: logEvent,
	}
}

// SetOnClick sets the click callback
func (eh *EventHandler) SetOnClick(callback func(pointer.Event)) *EventHandler {
	eh.onClick = callback
	return eh
}

// SetOnScroll sets the scroll callback
func (eh *EventHandler) SetOnScroll(callback func(float32)) *EventHandler {
	eh.onScroll = callback
	return eh
}

// SetOnHover sets the hover callback
func (eh *EventHandler) SetOnHover(callback func(bool)) *EventHandler {
	eh.onHover = callback
	return eh
}

// SetOnDrag sets the drag callback
func (eh *EventHandler) SetOnDrag(callback func(pointer.Event)) *EventHandler {
	eh.onDrag = callback
	return eh
}

// SetOnPress sets the press callback
func (eh *EventHandler) SetOnPress(callback func(pointer.Event)) *EventHandler {
	eh.onPress = callback
	return eh
}

// SetOnRelease sets the release callback
func (eh *EventHandler) SetOnRelease(callback func(pointer.Event)) *EventHandler {
	eh.onRelease = callback
	return eh
}

// AddToOps registers the event handler for events in the given ops
func (eh *EventHandler) AddToOps(ops *op.Ops) {
	// Register for pointer events using the proper Gio event system
	event.Op(ops, eh)
}

// ProcessEvents processes all events from the given context
func (eh *EventHandler) ProcessEvents(gtx layout.Context) {
	// Process scroll events first and re-emit them
	eh.processScrollEvents(gtx)

	// Process pointer events from the source
	eh.processPointerEvents(gtx)
}

// processScrollEvents processes scroll events from pointer events
func (eh *EventHandler) processScrollEvents(gtx layout.Context) {
	// Process scroll events from pointer events
	for {
		event, ok := gtx.Event(pointer.Filter{
			Kinds: pointer.Scroll,
		})
		if !ok {
			break
		}
		if pointerEvent, ok := event.(pointer.Event); ok {
			if pointerEvent.Kind == pointer.Scroll {
				log.I.F("Found SCROLL event: Y=%.1f", pointerEvent.Scroll.Y)
				eh.logEvent(fmt.Sprintf("SCROLL: Direction=%s, Y=%.1f",
					getScrollDirection(pointerEvent.Scroll.Y), pointerEvent.Scroll.Y))

				// Call scroll callback if set
				if eh.onScroll != nil {
					eh.onScroll(pointerEvent.Scroll.Y)
				}
			}
		}
	}
}

// getScrollDirection returns a string describing the scroll direction
func getScrollDirection(scrollY float32) string {
	if scrollY > 0 {
		return "Down"
	} else if scrollY < 0 {
		return "Up"
	}
	return "None"
}

// processPointerEvents processes pointer events from the source
func (eh *EventHandler) processPointerEvents(gtx layout.Context) {
	log.I.F("processPointerEvents called")

	// Process pointer events from the source
	pointerCount := 0
	for {
		ev, ok := gtx.Source.Event(pointer.Filter{
			Target: eh,
			Kinds:  pointer.Press | pointer.Release | pointer.Drag | pointer.Move | pointer.Enter | pointer.Leave | pointer.Cancel | pointer.Scroll,
		})
		if !ok {
			// Try to get any event without filter to see what's available
			ev2, ok2 := gtx.Source.Event(pointer.Filter{Target: eh})
			if ok2 {
				log.I.F("Found unfiltered event: %T", ev2)
			}
		}
		if !ok {
			break
		}
		pointerCount++
		if e, ok := ev.(pointer.Event); ok {
			// Track button state and determine which button was released
			var buttonInfo string

			switch e.Kind {
			case pointer.Press:
				eh.pressedButtons |= e.Buttons
				buttonInfo = e.Buttons.String()
				// Call press callback if set
				if eh.onPress != nil {
					eh.onPress(e)
				}
			case pointer.Release:
				// For release, show which button was released (the difference between old and new state)
				releasedButton := eh.pressedButtons &^ e.Buttons
				if releasedButton != 0 {
					buttonInfo = fmt.Sprintf("Released: %s", releasedButton.String())
				} else {
					buttonInfo = "Released: Unknown"
				}
				eh.pressedButtons = e.Buttons
				// Call release callback if set
				if eh.onRelease != nil {
					eh.onRelease(e)
				}
				// Check if this was a click (press followed by release)
				if releasedButton != 0 && eh.onClick != nil {
					eh.onClick(e)
				}
			case pointer.Drag:
				// Call drag callback if set
				if eh.onDrag != nil {
					eh.onDrag(e)
				}
				buttonInfo = e.Buttons.String()
			case pointer.Enter, pointer.Leave:
				// Call hover callback if set
				if eh.onHover != nil {
					eh.onHover(e.Kind == pointer.Enter)
				}
				buttonInfo = e.Buttons.String()
			case pointer.Scroll:
				// Handle scroll events
				log.I.F("Found SCROLL event: Y=%.1f, X=%.1f, Position=(%.1f,%.1f)", e.Scroll.Y, e.Scroll.X, e.Position.X, e.Position.Y)
				eh.logEvent(fmt.Sprintf("SCROLL: Direction=%s, Y=%.1f, X=%.1f",
					getScrollDirection(e.Scroll.Y), e.Scroll.Y, e.Scroll.X))

				// Call scroll callback if set
				if eh.onScroll != nil {
					eh.onScroll(e.Scroll.Y)
				}
				buttonInfo = fmt.Sprintf("Scroll: Y=%.1f, X=%.1f", e.Scroll.Y, e.Scroll.X)
			default:
				buttonInfo = e.Buttons.String()
			}

			log.I.F("Found pointer event: %s at position (%.1f,%.1f), buttons=%s", e.Kind, e.Position.X, e.Position.Y, buttonInfo)

			// Create event description
			eventDesc := fmt.Sprintf("POINTER: Kind=%s, Source=%s, Position=(%.1f,%.1f), Scroll=(%.1f,%.1f), Buttons=%s, Modifiers=%v, Time=%v",
				e.Kind, e.Source, e.Position.X, e.Position.Y, e.Scroll.X, e.Scroll.Y, buttonInfo, e.Modifiers, e.Time)

			eh.logEvent(eventDesc)
		} else {
			log.I.F("Found non-pointer event: %T", ev)
		}
	}
	if pointerCount > 0 {
		log.I.F("Processed %d pointer events", pointerCount)
	}
}

// HandleEvent implements event.Handler to capture all events
func (eh *EventHandler) HandleEvent(ev event.Event) {
	log.I.F("HandleEvent called with: %T", ev)
	switch e := ev.(type) {
	case pointer.Event:
		// Track button state and determine which button was released
		var buttonInfo string

		switch e.Kind {
		case pointer.Press:
			eh.pressedButtons |= e.Buttons
			buttonInfo = e.Buttons.String()
			// Call press callback if set
			if eh.onPress != nil {
				eh.onPress(e)
			}
		case pointer.Release:
			// For release, show which button was released (the difference between old and new state)
			releasedButton := eh.pressedButtons &^ e.Buttons
			if releasedButton != 0 {
				buttonInfo = fmt.Sprintf("Released: %s", releasedButton.String())
			} else {
				buttonInfo = "Released: Unknown"
			}
			eh.pressedButtons = e.Buttons
			// Call release callback if set
			if eh.onRelease != nil {
				eh.onRelease(e)
			}
			// Check if this was a click (press followed by release)
			if releasedButton != 0 && eh.onClick != nil {
				eh.onClick(e)
			}
		case pointer.Drag:
			// Call drag callback if set
			if eh.onDrag != nil {
				eh.onDrag(e)
			}
			buttonInfo = e.Buttons.String()
		case pointer.Enter, pointer.Leave:
			// Call hover callback if set
			if eh.onHover != nil {
				eh.onHover(e.Kind == pointer.Enter)
			}
			buttonInfo = e.Buttons.String()
		default:
			buttonInfo = e.Buttons.String()
		}

		// Create event description
		eventDesc := fmt.Sprintf("POINTER: Kind=%s, Source=%s, Position=(%.1f,%.1f), Scroll=(%.1f,%.1f), Buttons=%s, Modifiers=%v, Time=%v",
			e.Kind, e.Source, e.Position.X, e.Position.Y, e.Scroll.X, e.Scroll.Y, buttonInfo, e.Modifiers, e.Time)

		eh.logEvent(eventDesc)
	default:
		log.I.F("Unknown event type: %T", ev)
	}
}
