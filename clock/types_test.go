package clock

import (
	"testing"
)

func Test_NextRun(t *testing.T) {
	opt := NewOptions()
	opt.Start = "09:00"
	opt.End = "18:00"
	timing, _ := NewTiming(opt)
	timing.NextRun()

}
