package httpfs_test

import (
	"strings"
	"testing"

	"github.com/reiver/go-httpfs"
)

func TestPathJoin(t *testing.T) {

	tests := []struct{
		Root       string
		SubPath    string
		Expected   string
		ExpectedOK bool
	}{
		{
			Root:       "/srv/data",
			SubPath:    "",
			Expected:   "/srv/data",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/",
			Expected: "/srv/data",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/file.txt",
			Expected: "/srv/data/file.txt",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/images/photo.jpg",
			Expected: "/srv/data/images/photo.jpg",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/one/two/three.txt",
			Expected: "/srv/data/one/two/three.txt",
			ExpectedOK: true,
		},



		{
			Root:     "/srv/data",
			SubPath:  "/./file.txt",
			Expected: "/srv/data/file.txt",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/one/./two/./three.txt",
			Expected: "/srv/data/one/two/three.txt",
			ExpectedOK: true,
		},



		{
			Root:     "/srv/data",
			SubPath:  "/one/two/../two/three.txt",
			Expected: "/srv/data/one/two/three.txt",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/one/two/three/../../two/three.txt",
			Expected: "/srv/data/one/two/three.txt",
			ExpectedOK: true,
		},



		// Path traversal attacks — these MUST NOT escape root.
		{
			Root:     "/srv/data",
			SubPath:  "/..",
			Expected: "/srv/data",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/../",
			Expected: "/srv/data",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/../etc/passwd",
			Expected: "/srv/data/etc/passwd",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/../../etc/passwd",
			Expected: "/srv/data/etc/passwd",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/../../../etc/passwd",
			Expected: "/srv/data/etc/passwd",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/../../../../../../../etc/passwd",
			Expected: "/srv/data/etc/passwd",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/..%2f..%2f..%2fetc/passwd",
			Expected: "/srv/data/..%2f..%2f..%2fetc/passwd",
			ExpectedOK: true,
		},
		{
			Root:     "/srv/data",
			SubPath:  "/one/two/../../../etc/passwd",
			Expected: "/srv/data/etc/passwd",
			ExpectedOK: true,
		},
	}

	for testNumber, test := range tests {

		actualPath, actualOK := httpfs.PathJoin(test.Root, test.SubPath)

		{
			actual   := actualOK
			expected := test.ExpectedOK

			if expected != actual {
				t.Errorf("For test #%d, the actual value for 'path-join' 'ok' is not what was expected.", testNumber)
				t.Logf("EXPECTED: %t", expected)
				t.Logf("ACTUAL:   %t", actual)
				t.Logf("ROOT:        %q", test.Root)
				t.Logf("SUB-PATH:    %q", test.SubPath)
				t.Logf("ACTUAL-PATH: %q", actualPath)
				continue
			}
		}

		{
			actual   := actualPath
			expected := test.Expected

			if expected != actual {
				t.Errorf("For test #%d, the actual value for 'path-join' is not what was expected.", testNumber)
				t.Logf("EXPECTED: %q", expected)
				t.Logf("ACTUAL:   %q", actual)
				t.Logf("ROOT:     %q", test.Root)
				t.Logf("SUB-PATH: %q", test.SubPath)
				continue
			}

			// Extra safety check: if a result was returned, it must be under root.
			if "" != actual && !strings.HasPrefix(actual, test.Root) {
				t.Errorf("For test #%d, the actual path escapes the root directory!", testNumber)
				t.Logf("ROOT:     %q", test.Root)
				t.Logf("ACTUAL:   %q", actual)
				t.Logf("SUB-PATH: %q", test.SubPath)
				continue
			}
		}

	}
}
