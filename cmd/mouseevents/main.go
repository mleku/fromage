// SPDX-License-Identifier: Unlicense OR MIT

// GLFW doesn't build on OpenBSD and FreeBSD.
//go:build !openbsd && !freebsd && !android && !ios && !js
// +build !openbsd,!freebsd,!android,!ios,!js

package main

import (
	"context"
	"fmt"
	"image"
	"log"
	"runtime"
	"time"

	"gioui.org/f32"
	"gioui.org/font/gofont"
	"gioui.org/gpu"
	"gioui.org/io/event"
	"gioui.org/io/input"
	"gioui.org/io/pointer"
	"gioui.org/layout"
	"gioui.org/op"
	"gioui.org/op/clip"
	"gioui.org/op/paint"
	"gioui.org/text"
	"gioui.org/unit"
	"github.com/go-gl/gl/v3.3-core/gl"
	"github.com/go-gl/glfw/v3.3/glfw"
	"github.com/mleku/fromage"
)

// desktopGL is true when the (core, desktop) OpenGL should
// be used, false for OpenGL ES.
const desktopGL = runtime.GOOS == "darwin"

// Import aliases from fromage package
type (
	C = fromage.C
	D = fromage.D
	W = fromage.W
)

// MouseEvent represents a mouse event
type MouseEvent struct {
	Type      string
	Position  image.Point
	Timestamp string
}

// Application state
type AppState struct {
	lastEvent  MouseEvent
	theme      *fromage.Theme
	pointerTag interface{}
}

var appState *AppState

func main() {
	// Required by the OpenGL threading model.
	runtime.LockOSThread()

	err := glfw.Init()
	if err != nil {
		log.Fatal(err)
	}
	defer glfw.Terminate()

	// Gio assumes a sRGB backbuffer.
	glfw.WindowHint(glfw.SRGBCapable, glfw.True)
	glfw.WindowHint(glfw.ScaleToMonitor, glfw.True)
	glfw.WindowHint(glfw.CocoaRetinaFramebuffer, glfw.True)
	if desktopGL {
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 3)
		glfw.WindowHint(glfw.OpenGLProfile, glfw.OpenGLCoreProfile)
		glfw.WindowHint(glfw.OpenGLForwardCompatible, glfw.True)
	} else {
		glfw.WindowHint(glfw.ContextCreationAPI, glfw.EGLContextAPI)
		glfw.WindowHint(glfw.ClientAPI, glfw.OpenGLESAPI)
		glfw.WindowHint(glfw.ContextVersionMajor, 3)
		glfw.WindowHint(glfw.ContextVersionMinor, 0)
	}

	window, err := glfw.CreateWindow(400, 300, "Mouse Events Demo", nil, nil)
	if err != nil {
		log.Fatal(err)
	}

	window.MakeContextCurrent()

	err = gl.Init()
	if err != nil {
		log.Fatalf("gl.Init failed: %v", err)
	}
	// Enable sRGB.
	gl.Enable(gl.FRAMEBUFFER_SRGB)
	// Set up default VBA, required for the forward-compatible core profile.
	var defVBA uint32
	gl.GenVertexArrays(1, &defVBA)
	gl.BindVertexArray(defVBA)

	th := fromage.NewThemeWithMode(
		context.Background(),
		fromage.NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		fromage.ThemeModeDark,
	)

	// Initialize application state
	appState = &AppState{
		theme:      th,
		pointerTag: &struct{}{}, // Unique tag for pointer events
		lastEvent: MouseEvent{
			Type:      "No events yet",
			Position:  image.Pt(0, 0),
			Timestamp: "",
		},
	}

	var queue input.Router
	var ops op.Ops
	gpuCtx, err := gpu.New(gpu.OpenGL{ES: false, Shared: true})
	if err != nil {
		log.Fatal(err)
	}
	defer gpuCtx.Release()

	registerCallbacks(window, &queue)
	for !window.ShouldClose() {
		glfw.PollEvents()
		scale, _ := window.GetContentScale()
		width, height := window.GetFramebufferSize()
		sz := image.Point{X: width, Y: height}
		ops.Reset()
		gtx := layout.Context{
			Ops:    &ops,
			Now:    time.Now(),
			Source: queue.Source(),
			Metric: unit.Metric{
				PxPerDp: scale,
				PxPerSp: scale,
			},
			Constraints: layout.Exact(sz),
		}
		mainUI(gtx, th)
		gpuCtx.Frame(gtx.Ops, gpu.OpenGLRenderTarget{}, sz)
		queue.Frame(gtx.Ops)
		window.SwapBuffers()
	}
}

func mainUI(gtx layout.Context, th *fromage.Theme) {
	// Fill background with theme background color
	paint.Fill(gtx.Ops, th.Colors.Background())

	// Register for pointer events over the entire window area
	r := image.Rectangle{Max: gtx.Constraints.Max}
	area := clip.Rect(r).Push(gtx.Ops)
	event.Op(gtx.Ops, appState.pointerTag)
	area.Pop()

	// Handle pointer events
	for {
		ev, ok := gtx.Event(pointer.Filter{
			Target: appState.pointerTag,
			Kinds:  pointer.Press | pointer.Release | pointer.Move,
		})
		if !ok {
			break
		}
		if e, ok := ev.(pointer.Event); ok {
			clickPos := image.Pt(int(e.Position.X), int(e.Position.Y))

			// Determine which button was pressed
			var buttonType string
			switch e.Kind {
			case pointer.Press:
				switch {
				case e.Buttons == pointer.ButtonPrimary:
					buttonType = "Left Click"
				case e.Buttons == pointer.ButtonSecondary:
					buttonType = "Right Click"
				case e.Buttons == pointer.ButtonTertiary:
					buttonType = "Middle Click"
				default:
					buttonType = fmt.Sprintf("Button %d", e.Buttons)
				}
			case pointer.Release:
				switch {
				case e.Buttons == pointer.ButtonPrimary:
					buttonType = "Left Release"
				case e.Buttons == pointer.ButtonSecondary:
					buttonType = "Right Release"
				case e.Buttons == pointer.ButtonTertiary:
					buttonType = "Middle Release"
				default:
					buttonType = fmt.Sprintf("Button %d Release", e.Buttons)
				}
			case pointer.Move:
				buttonType = "Mouse Move"
			}

			appState.lastEvent = MouseEvent{
				Type:      buttonType,
				Position:  clickPos,
				Timestamp: fmt.Sprintf("%v", e.Time),
			}

			log.Printf("Mouse event: %s at (%d, %d)", appState.lastEvent.Type, clickPos.X, clickPos.Y)
		}
	}

	// Layout the UI
	th.CenteredColumn().
		Rigid(func(g C) D {
			// Title
			return th.H3("Mouse Events Demo").
				Color(th.Colors.OnBackground()).
				Alignment(text.Middle).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Instructions
			return th.Body1("Click anywhere in this window to see mouse events").
				Color(th.Colors.OnSurfaceVariant()).
				Alignment(text.Middle).
				Layout(g)
		}).
		Rigid(func(g C) D {
			// Event display area
			return th.NewCard(
				func(g C) D {
					return th.VFlex().
						SpaceEvenly().
						Rigid(func(g C) D {
							return th.Body2("Last Mouse Event:").
								Color(th.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Body1(fmt.Sprintf("Type: %s", appState.lastEvent.Type)).
								Color(th.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Body1(fmt.Sprintf("Position: (%d, %d)", appState.lastEvent.Position.X, appState.lastEvent.Position.Y)).
								Color(th.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(g)
						}).
						Rigid(func(g C) D {
							return th.Body1(fmt.Sprintf("Time: %s", appState.lastEvent.Timestamp)).
								Color(th.Colors.OnSurface()).
								Alignment(text.Start).
								Layout(g)
						}).
						Layout(g)
				},
			).CornerRadius(8).Padding(unit.Dp(16)).Layout(g)
		}).
		Layout(gtx)
}

func registerCallbacks(window *glfw.Window, q *input.Router) {
	var btns pointer.Buttons
	beginning := time.Now()
	var lastPos f32.Point

	window.SetCursorPosCallback(func(w *glfw.Window, xpos float64, ypos float64) {
		scale := float32(1)
		if runtime.GOOS == "darwin" {
			// macOS cursor positions are not scaled to the underlying framebuffer
			// size when CocoaRetinaFramebuffer is true.
			scale, _ = w.GetContentScale()
		}
		lastPos = f32.Point{X: float32(xpos) * scale, Y: float32(ypos) * scale}
		e := pointer.Event{
			Kind:     pointer.Move,
			Position: lastPos,
			Source:   pointer.Mouse,
			Time:     time.Since(beginning),
			Buttons:  btns,
		}
		q.Queue(e)
	})

	window.SetMouseButtonCallback(func(w *glfw.Window, button glfw.MouseButton, action glfw.Action, mods glfw.ModifierKey) {
		var btn pointer.Buttons
		switch button {
		case glfw.MouseButton1:
			btn = pointer.ButtonPrimary
		case glfw.MouseButton2:
			btn = pointer.ButtonSecondary
		case glfw.MouseButton3:
			btn = pointer.ButtonTertiary
		}
		var typ pointer.Kind
		switch action {
		case glfw.Release:
			typ = pointer.Release
			btns &^= btn
		case glfw.Press:
			typ = pointer.Press
			btns |= btn
		}
		e := pointer.Event{
			Kind:     typ,
			Source:   pointer.Mouse,
			Time:     time.Since(beginning),
			Position: lastPos,
			Buttons:  btns,
		}
		q.Queue(e)
	})
}
