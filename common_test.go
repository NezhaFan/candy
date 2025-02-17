package candy

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestStrLen(t *testing.T) {
	s := "1BC四五六👍"
	assert.Equal(t, StrLen(s), 7)
}

func TestDiv(t *testing.T) {
	assert.Equal(t, Div(1, 0, 2), 0.)

	// 33.3333
	assert.Equal(t, Div(100, 3, 2), 33.33)
	assert.Equal(t, Div(100, 3, 1), 33.3)
	assert.Equal(t, Div(100, 3, 0), 33.)

	// 1.5714...
	assert.Equal(t, Div(11, 7, 3), 1.571)
	assert.Equal(t, Div(11, 7, 2), 1.57)
	assert.Equal(t, Div(11, 7, 1), 1.6)

	assert.Equal(t, Div(-11, 7, 3), -1.571)
	assert.Equal(t, Div(-11, 7, 2), -1.57)
	assert.Equal(t, Div(-11, 7, 1), -1.6)
}

func TestRound(t *testing.T) {
	pi := 3.14159
	assert.Equal(t, Round(pi, 4), 3.1416)
	assert.Equal(t, Round(pi, 3), 3.142)
	assert.Equal(t, Round(pi, 2), 3.14)
	assert.Equal(t, Round(pi, 1), 3.1)
	assert.Equal(t, Round(pi, 0), 3.)

	pi = -pi
	assert.Equal(t, Round(pi, 4), -3.1416)
	assert.Equal(t, Round(pi, 3), -3.142)
	assert.Equal(t, Round(pi, 2), -3.14)
	assert.Equal(t, Round(pi, 1), -3.1)
	assert.Equal(t, Round(pi, 0), -3.)
}

func TestCallers(t *testing.T) {
	t1 := func(f func()) {
		f()
	}

	t2 := func(f func()) {
		f()
	}

	t1(func() {
		t2(func() {
			defer func() {
				_ = recover()
				fmt.Println(Callers())
			}()
			panic("xxx")
		})
	})
}
