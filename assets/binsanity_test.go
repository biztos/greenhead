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

const BinsanityAssetMissing = "webui/root.html--NOPE"
const BinsanityAssetPresent = "doc/readme.md"
const BinsanityAssetPresentSum = "9e5d4dd53626efe1e858ca146d950b0284eaed425271a444685c7cc0811ed1ca"

var BinsanityAssetNames = []string{

	"agents/chatty.toml",
	"agents/marvin.toml",
	"agents/pirate.toml",
	"agents/tictactoe.toml",
	"doc/.DS_Store",
	"doc/api.md",
	"doc/config.md",
	"doc/license.md",
	"doc/licenses.md",
	"doc/people.md",
	"doc/readme.md",
	"doc/tools/tictactoe.md",
	"webui/app.html",
	"webui/err-badkey.html",
	"webui/err-noagents.html",
	"webui/favicon.png",
	"webui/favicon.svg",
	"webui/greenhead-150x225.png",
	"webui/greenhead-300x450.png",
	"webui/greenhead.png",
	"webui/root.html",
}

var BinsanityAssetSums = []string{
	"c266499f9d5d56104f5060483ba3df466f49e7ee3f18992fdad4e018c1bf57a6",
	"e49aa1ea1710e0ad809e23a8429a683f017472cfb5ba8eec99ef56964acdf401",
	"ec08f63a25b7edda382315885e0ff961accb48e44dca179a28605a1b88e88eab",
	"07b6369cd8d991b14432ad1468bf29c0a5f78c3d4b2d5d57c80638481466f3db",
	"d65165279105ca6773180500688df4bdc69a2c7b771752f0a46ef120b7fd8ec3",
	"e8619367edc76ba0ca73a0857368440ae2e1642dfbb9d58088679c3724eacb41",
	"cbb1008c9e35246a05ad0176e954901ea87da171858f4b0188958e9f82c6ee1f",
	"7f84f67d02944c995bde337c4f70f3f3c1556684cf8298731a51bb3014130061",
	"7c4899a35d0d48162f54454cce04343fd402fa3da40a36c1ffc611a382969f68",
	"e66573e846581a3110de4f575a9be8937a52bb47979c5c631367051c4d72f398",
	"9e5d4dd53626efe1e858ca146d950b0284eaed425271a444685c7cc0811ed1ca",
	"9d9c4535cf0577291fc57b0a6daebda9a3d9490bef897604631a220e0ceac60b",
	"6a70b0cb86f7ed71e4f32f2559514944bc2465811d45e889e7a3e2e7712a52a8",
	"a184d7b13e8dd11c9a7f1ac4073a9268a0dd053a7e45771fcbf84728d5b311fe",
	"a184d7b13e8dd11c9a7f1ac4073a9268a0dd053a7e45771fcbf84728d5b311fe",
	"0e11ef0572db7297c342d6a9c1cc78b54df3871dd92edb548e8ae06579bc1a66",
	"42bc5e56e7412df8c21953de78c9e204dd7da8152e43854e8ed0e832bfac153e",
	"86b742dcea3b22359c9e413b4011afe9521d51c00ff5b03262248152d3dca7ab",
	"cae0ac498aeef33f365c09df814a104f96d4df7a228d383e2fec292c89d4c3c8",
	"d7ca9a8bdf6043e8426e9e1aad254da7d5a35990ef0814ac24985f189e5031de",
	"a071ad011aa90eba8affb681aa3a17192105861edba5f2b97d863171d55531f2",
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
