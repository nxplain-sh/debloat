package packages

import (
	"reflect"
	"strings"
	"testing"
)

func TestParse(t *testing.T) {
	in := `# Header comment
com.vivo.browser
  com.android.notes   # trailing comment
# full-line comment

	com.example.tabbed
com.vivo.browser

`
	want := []string{
		"com.vivo.browser",
		"com.android.notes",
		"com.example.tabbed",
		"com.vivo.browser", // duplicates are preserved for the caller
	}
	got, err := Parse(strings.NewReader(in))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if !reflect.DeepEqual(got, want) {
		t.Fatalf("Parse() = %#v, want %#v", got, want)
	}
}

func TestParseEmpty(t *testing.T) {
	got, err := Parse(strings.NewReader("# only comments\n\n   \n"))
	if err != nil {
		t.Fatalf("Parse returned error: %v", err)
	}
	if len(got) != 0 {
		t.Fatalf("Parse() = %#v, want empty", got)
	}
}
