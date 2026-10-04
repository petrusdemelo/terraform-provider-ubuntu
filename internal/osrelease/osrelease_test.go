package osrelease

import (
	"maps"
	"testing"
)

func TestParse(t *testing.T) {
	content := `PRETTY_NAME="Ubuntu 24.04.1 LTS"
NAME="Ubuntu"
VERSION_ID="24.04"
VERSION_CODENAME=noble
# a comment

ID=ubuntu
SINGLE='quoted value'
ESCAPED="say \"hi\""
not an assignment
`

	want := map[string]string{
		"PRETTY_NAME":      "Ubuntu 24.04.1 LTS",
		"NAME":             "Ubuntu",
		"VERSION_ID":       "24.04",
		"VERSION_CODENAME": "noble",
		"ID":               "ubuntu",
		"SINGLE":           "quoted value",
		"ESCAPED":          `say "hi"`,
	}

	if got := Parse(content); !maps.Equal(got, want) {
		t.Errorf("Parse() = %v, want %v", got, want)
	}
}

func FuzzParse(f *testing.F) {
	f.Add("ID=ubuntu\nVERSION_ID=\"24.04\"\n")
	f.Add("A='x'\n# c\n=\n")

	f.Fuzz(func(_ *testing.T, content string) {
		Parse(content)
	})
}
