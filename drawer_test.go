package fromage

import (
	"testing"
	"time"

	"gio.mleku.dev/unit"
)

func TestDrawerCreation(t *testing.T) {
	// Create a mock theme
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	// Test basic drawer creation
	drawer := w.NewDrawer()
	if drawer == nil {
		t.Fatal("Expected drawer to be created, got nil")
	}

	// Test default values
	if drawer.position != DrawerLeft {
		t.Errorf("Expected default position to be DrawerLeft, got %v", drawer.position)
	}

	if drawer.width != unit.Dp(280) {
		t.Errorf("Expected default width to be 280dp, got %v", drawer.width)
	}

	if drawer.height != unit.Dp(200) {
		t.Errorf("Expected default height to be 200dp, got %v", drawer.height)
	}

	if drawer.isVisible {
		t.Error("Expected drawer to start invisible")
	}
}

func TestDrawerPosition(t *testing.T) {
	// Create a mock theme
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	// Test position setting
	drawer := w.NewDrawer().Position(DrawerRight)
	if drawer.position != DrawerRight {
		t.Errorf("Expected position to be DrawerRight, got %v", drawer.position)
	}

	drawer.Position(DrawerTop)
	if drawer.position != DrawerTop {
		t.Errorf("Expected position to be DrawerTop, got %v", drawer.position)
	}

	drawer.Position(DrawerBottom)
	if drawer.position != DrawerBottom {
		t.Errorf("Expected position to be DrawerBottom, got %v", drawer.position)
	}
}

func TestDrawerVisibility(t *testing.T) {
	// Create a mock theme
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	drawer := w.NewDrawer()

	// Test initial state
	if drawer.IsVisible() {
		t.Error("Expected drawer to start invisible")
	}

	// Test show
	drawer.Show()
	if !drawer.IsVisible() {
		t.Error("Expected drawer to be visible after Show()")
	}

	// Test hide
	drawer.Hide()
	if drawer.IsVisible() {
		t.Error("Expected drawer to be invisible after Hide()")
	}

	// Test toggle
	drawer.Toggle()
	if !drawer.IsVisible() {
		t.Error("Expected drawer to be visible after Toggle()")
	}

	drawer.Toggle()
	if drawer.IsVisible() {
		t.Error("Expected drawer to be invisible after second Toggle()")
	}
}

func TestDrawerWithControlsCreation(t *testing.T) {
	// Create a mock theme
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	// Test drawer with controls creation
	dwc := w.NewDrawerWithControls()
	if dwc == nil {
		t.Fatal("Expected drawer with controls to be created, got nil")
	}

	if dwc.drawer == nil {
		t.Fatal("Expected drawer to be created, got nil")
	}

	if dwc.radioGroup == nil {
		t.Fatal("Expected radio group to be created, got nil")
	}

	if dwc.currentPos != DrawerLeft {
		t.Errorf("Expected default position to be DrawerLeft, got %v", dwc.currentPos)
	}
}

func TestDrawerWithControlsPosition(t *testing.T) {
	// Create a mock theme
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	dwc := w.NewDrawerWithControls()

	// Test position setting
	dwc.SetPosition(DrawerRight)
	if dwc.GetCurrentPosition() != DrawerRight {
		t.Errorf("Expected position to be DrawerRight, got %v", dwc.GetCurrentPosition())
	}

	dwc.SetPosition(DrawerTop)
	if dwc.GetCurrentPosition() != DrawerTop {
		t.Errorf("Expected position to be DrawerTop, got %v", dwc.GetCurrentPosition())
	}

	dwc.SetPosition(DrawerBottom)
	if dwc.GetCurrentPosition() != DrawerBottom {
		t.Errorf("Expected position to be DrawerBottom, got %v", dwc.GetCurrentPosition())
	}
}

func TestDrawerAnimation(t *testing.T) {
	// Create a mock theme
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	drawer := w.NewDrawer()

	// Test animation start
	now := time.Now()
	drawer.startAnimation(now)
	if !drawer.isAnimating {
		t.Error("Expected drawer to be animating after startAnimation")
	}

	if !drawer.animationStarted {
		t.Error("Expected animationStarted to be true")
	}

	if drawer.isFadingOut {
		t.Error("Expected isFadingOut to be false for slide-in animation")
	}

	// Test fade out start
	drawer.startFadeOut()
	if !drawer.isAnimating {
		t.Error("Expected drawer to be animating after startFadeOut")
	}

	if !drawer.isFadingOut {
		t.Error("Expected isFadingOut to be true for fade-out animation")
	}
}

func TestDrawerConvenienceMethods(t *testing.T) {
	// Create a mock theme
	th := &Theme{
		TextSize: unit.Dp(16),
		Pool:     &Pool{},
	}

	// Create a mock window
	w := &Window{Theme: th}

	// Test convenience methods
	leftDrawer := w.LeftDrawer()
	if leftDrawer.position != DrawerLeft {
		t.Errorf("Expected LeftDrawer to have DrawerLeft position, got %v", leftDrawer.position)
	}

	rightDrawer := w.RightDrawer()
	if rightDrawer.position != DrawerRight {
		t.Errorf("Expected RightDrawer to have DrawerRight position, got %v", rightDrawer.position)
	}

	topDrawer := w.TopDrawer()
	if topDrawer.position != DrawerTop {
		t.Errorf("Expected TopDrawer to have DrawerTop position, got %v", topDrawer.position)
	}

	bottomDrawer := w.BottomDrawer()
	if bottomDrawer.position != DrawerBottom {
		t.Errorf("Expected BottomDrawer to have DrawerBottom position, got %v", bottomDrawer.position)
	}
}
