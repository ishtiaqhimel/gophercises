package link

import (
	"strings"
	"testing"

	"gotest.tools/assert"
)

func TestParse(t *testing.T) {
	t.Run("test parse", func(t *testing.T) {
		html := `
		<html>
<body>
  <h1>Hello!</h1>
  <a href="/other-page">A link to another page</a>
  <a href="/page-two">A link to page two</a>
</body>
</html>
		`
		r := strings.NewReader(html)
		links, err := Parse(r)
		assert.NilError(t, err, "failed to parse")

		expected := []Link{
			{
				"/other-page",
				"A link to another page",
			},
			{
				"/page-two",
				"A link to page two",
			},
		}
		assert.DeepEqual(t, links, expected)
	})
}
