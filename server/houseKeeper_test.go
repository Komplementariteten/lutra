package server

import (
	"testing"
	"time"
)

func TestHouseKeeper_Ctrl(t *testing.T) {
	h := &HTTPServer{}
	ctrl := StartNewHouseKeeper(h, 1*time.Second)
	time.Sleep(5 * time.Second)
	ctrl.Quit <- true
}
