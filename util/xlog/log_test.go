package xlog

import (
	"errors"
	"fmt"
	"testing"
)

func TestWrap(t *testing.T) {
	fmt.Println(Wrap(errors.New("test error")))
}
