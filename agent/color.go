package agent

// colors... maybe put in other package but...
import (
	"fmt"
	"image/color"
	"io"
	"math"
	"sort"
	"strings"

	pcolor "github.com/fatih/color"
	"golang.org/x/image/colornames"
)

var ErrColorNotSupported = fmt.Errorf("color not supported")

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
	col, err := PrintColor(fg_name, bg_name)
	if err != nil {
		return nil, err
	}

	// Alas, the PrintFunc thing isn't compatible with actual fmt.Print.
	//
	return col.PrintFunc(), nil

}

// PrintColor returns the color which will provide a PrintFunc et al.
func PrintColor(fg_name, bg_name string) (*pcolor.Color, error) {

	fg, have := colornames.Map[strings.ToLower(fg_name)] // 147 should be enough for anyone!
	if !have {
		return nil, fmt.Errorf("%w: %s", ErrColorNotSupported, fg_name)
	}
	col := pcolor.RGB(int(fg.R), int(fg.G), int(fg.B))
	if bg_name != "" {
		bg, have := colornames.Map[strings.ToLower(bg_name)]
		if !have {
			return nil, fmt.Errorf("%w: %s", ErrColorNotSupported, bg_name)
		}
		col = col.AddBgRGB(int(bg.R), int(bg.G), int(bg.B))
	}
	return col, nil
}

// PrintColorPairSample prints a sample of the standard and pair colors for
// fg and bg to w, preceded by prefix.
func PrintColorPairSample(w io.Writer, fg, bg, prefix string) error {

	col, err := PrintColor(fg, bg)
	if err != nil {
		return err
	}
	col.Fprint(w, prefix, colorMsg(fg, bg))

	pair_fg, _ := FindComplementaryColor(fg) // known-good from above.
	pair_bg, _ := FindComplementaryColor(bg)
	pair_col, err := PrintColor(pair_fg, pair_bg)
	if err != nil {
		return err
	}
	pair_col.Fprintln(w, " The pair color: ", colorMsg(pair_fg, pair_bg))

	return nil
}

func colorMsg(fg, bg string) string {
	if fg != "" && bg != "" {
		return fmt.Sprintf("%q on %q is a lovely combination.", fg, bg)
	}
	if fg != "" {
		return fmt.Sprintf("%q is a lovely color by itself.", fg)
	}
	if bg != "" {
		return fmt.Sprintf("%q is our lovely background color.", bg)
	}
	return "the default colors will have to do."
}

// ColorDistance calculates a simple Euclidean distance between colors in
// RGB space.
func ColorDistance(c1, c2 color.RGBA) float64 {
	rDiff := float64(c1.R) - float64(c2.R)
	gDiff := float64(c1.G) - float64(c2.G)
	bDiff := float64(c1.B) - float64(c2.B)

	return math.Sqrt(rDiff*rDiff + gDiff*gDiff + bDiff*bDiff)
}

// FindComplementaryColor finds a color that is different but not too
// different.
//
// It is potentially quite slow, and should only be called at startup.
func FindComplementaryColor(name string) (string, error) {

	if name == "" {
		return "", nil
	}

	baseColor, have := colornames.Map[strings.ToLower(name)]
	if !have {
		return "", fmt.Errorf("%w: %s", ErrColorNotSupported, name)
	}

	// Define what's "too similar" and "too different"
	const minDistance = 100.0 // Colors should be at least this different
	const maxDistance = 200.0 // But not more than this different

	type colorPair struct {
		name     string
		distance float64
	}

	var candidates []colorPair

	// Search through all colornames
	for name, c := range colornames.Map {
		rgba := color.RGBA{c.R, c.G, c.B, c.A}
		dist := ColorDistance(baseColor, rgba)

		// Skip colors that are too similar or too different
		if dist < minDistance || dist > maxDistance {
			continue
		}

		candidates = append(candidates, colorPair{name, dist})
	}

	// Sort candidates by distance
	sort.Slice(candidates, func(i, j int) bool {
		return candidates[i].distance < candidates[j].distance
	})

	// Return the closest candidate that meets our criteria
	if len(candidates) > 0 {
		return candidates[0].name, nil
	}

	// Fallback is...?
	if baseColor != colornames.Lightblue {
		return "lightblue", nil
	} else {
		return "black", nil // or what?
	}
}
