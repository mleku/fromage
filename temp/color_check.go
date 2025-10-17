package main

import (
	"context"
	"fmt"

	"gio.mleku.dev/font/gofont"
	"gio.mleku.dev/text"
	"gio.mleku.dev/unit"
	"github.com/mleku/fromage"
)

func main() {
	th := fromage.NewThemeWithMode(
		context.Background(),
		fromage.NewColors,
		text.NewShaper(text.WithCollection(gofont.Collection())),
		unit.Dp(16),
		fromage.ThemeModeLight,
	)

	fmt.Printf("Background: %+v\n", th.Colors.Background())
	fmt.Printf("Surface: %+v\n", th.Colors.Surface())
	fmt.Printf("OnBackground: %+v\n", th.Colors.OnBackground())
	fmt.Printf("Primary: %+v\n", th.Colors.Primary())
	fmt.Printf("Secondary: %+v\n", th.Colors.Secondary())
}
