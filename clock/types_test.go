package clock

import (
	"testing"
)

func Test_NextRun(t *testing.T) {
	s, err := nextRun("17:32", "18:00", 0)

	if err != nil {
		t.Error(err)
	}
	t.Log(s)

}
