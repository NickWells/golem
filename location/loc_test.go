package location_test

import (
	"github.com/nickwells/golem/location"
	"testing"
)

func TestLocString(t *testing.T) {
	l := location.New("test1")
	expectedStr := "test1:0"
	if s := l.String(); s != expectedStr {
		t.Error("the string represening the location should be: '" +
			expectedStr + "', is: '" + s + "'")
	}

	l.Incr()
	expectedStr = "test1:1"
	if s := l.String(); s != expectedStr {
		t.Error(
			"after Incr the string representing the location should be: '" +
				expectedStr + "', is: '" + s + "'")
	}
}
