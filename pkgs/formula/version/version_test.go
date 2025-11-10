package version

import "testing"

func TestVersion(t *testing.T) {
	t.Log(Compare("1.0.0", "1.0.0-pre"))
}
