package atom

import (
	"fmt"
	"testing"
)

func assertEqual(t *testing.T, test string, ref string) {
	if test != ref {
		t.Error(fmt.Sprintf("Expected %s to equal %s", test, ref))
	}
}
