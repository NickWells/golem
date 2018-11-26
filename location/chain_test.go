package location_test

import (
	"github.com/nickwells/golem/location"
	"testing"
)

func TestChainDesc(t *testing.T) {
	fl1 := location.New("test1")
	fl2 := location.New("test2")
	fl3 := location.New("test3")

	includeChain := location.LocChain{*fl1, *fl2, *fl3}

	expectedChainDesc := "test1:0 --> test2:0 --> test3:0"
	if s := includeChain.String(); s != expectedChainDesc {
		t.Error("the string representing the chain should be: '" +
			expectedChainDesc + "', is: '" + s + "'")
	}
}
