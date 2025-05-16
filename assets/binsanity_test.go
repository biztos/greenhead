/* binsanity_test.go - auto-generated; edit at your own peril!

To test the checksums for all content, set the environment variable
BINSANITY_TEST_CONTENT to one of: Y,YES,T,TRUE,1 (the Truthy Shortlist).

More info: https://github.com/biztos/binsanity

*/

package assets_test

import (
	"crypto/sha256"
	"fmt"
	"os"
	"strings"
	"testing"

	"github.com/biztos/greenhead/assets"
)

const BinsanityAssetMissing = "doc/tools/tictactoe.md--NOPE"
const BinsanityAssetPresent = "api-ui/ghd.js"
const BinsanityAssetPresentSum = "5c86e3a3ffba9e48990f71547b25befcfa27d56efa9bb48144d803580fc79749"

var BinsanityAssetNames = []string{

	"agents/chatty.toml",
	"agents/marvin.toml",
	"agents/pirate.toml",
	"agents/tictactoe.toml",
	"api-ui/README.md",
	"api-ui/app.html",
	"api-ui/err-badkey.html",
	"api-ui/err-noagents.html",
	"api-ui/favicon.svg",
	"api-ui/ghd.css",
	"api-ui/ghd.js",
	"api-ui/main.html",
	"api-ui/root.html",
	"doc/api.md",
	"doc/config.md",
	"doc/license.md",
	"doc/licenses.md",
	"doc/people.md",
	"doc/readme.md",
	"doc/tools/tictactoe.md",
}

var BinsanityAssetSums = []string{
	"c266499f9d5d56104f5060483ba3df466f49e7ee3f18992fdad4e018c1bf57a6",
	"e49aa1ea1710e0ad809e23a8429a683f017472cfb5ba8eec99ef56964acdf401",
	"ec08f63a25b7edda382315885e0ff961accb48e44dca179a28605a1b88e88eab",
	"9759c5116e99ea5db017462318129482e48d3697634d70a1d6ac8407e9a2ff68",
	"f8827a87cdbc60a1063758380210cb293ce57c5539050365c9200635b8262fdf",
	"450bd7339761efa07f7b4792be9bfaf9a81487e3fdc5e8965917f1bae26618f6",
	"a184d7b13e8dd11c9a7f1ac4073a9268a0dd053a7e45771fcbf84728d5b311fe",
	"a184d7b13e8dd11c9a7f1ac4073a9268a0dd053a7e45771fcbf84728d5b311fe",
	"42bc5e56e7412df8c21953de78c9e204dd7da8152e43854e8ed0e832bfac153e",
	"726d3cece1d97b0a26b297ccffe28e62b1c46095419a99692333404cb9b7f1d2",
	"5c86e3a3ffba9e48990f71547b25befcfa27d56efa9bb48144d803580fc79749",
	"5b0a97f5d0491f0d7aac962d5093fb11b883b110d4f5465fa9c00e2788e5caf1",
	"0e687d30134afbcaa994307be7cad78419847916db71b54f21f30014bf9dbc7d",
	"6376dd2474462ba10c9aa26a77de0654c2fc748aaaa6fac2680d0122e1419cbf",
	"b13e3e9c5a85a8ebaf29f01877d6cdf16b23cfe556b00bf85ecb726971b276a0",
	"7f84f67d02944c995bde337c4f70f3f3c1556684cf8298731a51bb3014130061",
	"7c4899a35d0d48162f54454cce04343fd402fa3da40a36c1ffc611a382969f68",
	"e66573e846581a3110de4f575a9be8937a52bb47979c5c631367051c4d72f398",
	"d3864eb192507738a60b34b04c6665819975a612eecfcb208a602c70115cc03d",
	"9d9c4535cf0577291fc57b0a6daebda9a3d9490bef897604631a220e0ceac60b",
}

func TestAssetNames(t *testing.T) {

	names := assets.AssetNames()
	if len(names) != len(BinsanityAssetNames) {
		t.Fatalf("Wrong number of names:\n  expected: %d\n  actual: %d",
			len(BinsanityAssetNames), len(names))
	}

	// ...moments when you really miss Testify... but NO deps for the
	// generated files!
	for idx, n := range names {
		if n != BinsanityAssetNames[idx] {
			t.Fatalf("Mismatch at %d:\n  expected: %s\n  actual: %s",
				idx, BinsanityAssetNames[idx], n)
		}
	}

}

func TestAssetNotFound(t *testing.T) {

	_, err := assets.Asset(BinsanityAssetMissing)
	if err == nil {
		t.Fatal("No error for missing asset.")
	}
	if err.Error() != "Asset not found." {
		t.Fatal("Wrong error for missing asset.")
	}
}

func TestAssetFound(t *testing.T) {

	b, err := assets.Asset(BinsanityAssetPresent)
	if err != nil {
		t.Fatal("Error for asset that should not be missing.")
	}
	sum := fmt.Sprintf("%x", sha256.Sum256(b))
	if sum != BinsanityAssetPresentSum {
		t.Fatal("Wrong sha256 sum for asset data.")
	}
}

func TestMustAssetNotFound(t *testing.T) {

	exp := "Asset not found."
	panicky := func() { assets.MustAssetString(BinsanityAssetMissing) }
	AssertPanicsWith(t, panicky, exp, "MustAsset (not found)")

}

func TestMustAssetFound(t *testing.T) {

	b := assets.MustAsset(BinsanityAssetPresent)
	sum := fmt.Sprintf("%x", sha256.Sum256(b))
	if sum != BinsanityAssetPresentSum {
		t.Fatal("Wrong sha256 sum for asset data.")
	}

}

func TestMustAssetStringNotFound(t *testing.T) {

	exp := "Asset not found."
	panicky := func() { assets.MustAssetString(BinsanityAssetMissing) }
	AssertPanicsWith(t, panicky, exp, "MustAssetString (not found)")

}

func TestMustAssetStringFound(t *testing.T) {

	s := assets.MustAssetString(BinsanityAssetPresent)
	sum := fmt.Sprintf("%x", sha256.Sum256([]byte(s)))
	if sum != BinsanityAssetPresentSum {
		t.Fatal("Wrong sha256 sum for asset data.")
	}

}

func TestAssetSums(t *testing.T) {
	var want_tests bool
	// This is a little bit overkill but people have habits right?
	boolish := map[string]bool{
		"Y":    true,
		"YES":  true,
		"T":    true,
		"TRUE": true,
		"1":    true,
	}
	flag := strings.ToUpper(os.Getenv("BINSANITY_TEST_CONTENT"))
	want_tests = boolish[flag]
	if !want_tests {
		t.Skip()
		return
	}
	for idx, name := range BinsanityAssetNames {
		b, err := assets.Asset(name)
		if err != nil {
			t.Fatalf("%s: %v", name, err)
		}
		exp := BinsanityAssetSums[idx]
		sum := fmt.Sprintf("%x", sha256.Sum256(b))
		if sum != exp {
			t.Fatalf("Wrong sha256 sum for data of: %s\n  expected: %s\n    actual: %s",
				name, exp, sum)
		}
	}
}

// For a more useful version of this see: https://github.com/biztos/testig
func AssertPanicsWith(t *testing.T, f func(), exp string, msg string) {

	panicked := false
	got := ""
	func() {
		defer func() {
			if r := recover(); r != nil {
				panicked = true
				got = fmt.Sprintf("%s", r)
			}
		}()
		f()
	}()

	if !panicked {
		t.Fatalf("Function did not panic: %s", msg)
	} else if got != exp {

		t.Fatalf("Panic not as expected: %s\n  expected: %s\n    actual: %s",
			msg, exp, got)
	}

	// (In go testing, success is silent.)

}
