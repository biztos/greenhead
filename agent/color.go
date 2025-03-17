package agent

// colors... maybe put in other package but...
import (
	"fmt"
	"strings"

	"github.com/fatih/color"
	"golang.org/x/image/colornames"
)

// ColorPrintFunc returns the colorized print function to use for streaming.
//
// If both input color names are blank, a simple fmt.Print wrapper is used.
func ColorPrintFunc(fg_name, bg_name string) (func(a ...any), error) {

	if fg_name == "" && bg_name == "" {
		// We have to wrap this because color.PrintFunc has other signature.
		return func(a ...any) {
			fmt.Print(a...)
		}, nil
	}

	fg, have := colornames.Map[strings.ToLower(fg_name)] // 147 should be enough for anyone!
	if !have {
		return nil, fmt.Errorf("color not supported: %s", fg_name)
	}
	col := color.RGB(int(fg.R), int(fg.G), int(fg.B))
	if bg_name != "" {
		bg, have := colornames.Map[strings.ToLower(bg_name)]
		if !have {
			return nil, fmt.Errorf("bg color not supported: %s", bg_name)
		}
		col = col.AddBgRGB(int(bg.R), int(bg.G), int(bg.B))
	}

	// Alas, the PrintFunc thing isn't compatible with actual fmt.Print.
	//
	return col.PrintFunc(), nil

}
