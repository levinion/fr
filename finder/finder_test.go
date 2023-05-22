package finder

import "testing"

func TestExtNotInList(t *testing.T) {
	//asume to be false
	b1 := itemNotInList("mkv", []string{"html", "mp4", "mkv"})
	if b1 {
		t.Error("b1:", b1)
	}
	//asume to be true
	b2 := itemNotInList("mkv", []string{"html", "mp4"})
	if !b2 {
		t.Error("b2:", b2)
	}
}
