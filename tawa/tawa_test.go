package tawa

import (
	"fmt"
	"testing"
)

func TestUrl(t *testing.T) {
	ta, err := New("redis://localhost:6379/3")
	if err != nil {
		t.Error(err)
	}
	fmt.Println(ta)

}
