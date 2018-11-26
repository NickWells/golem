package location_test

import (
	"fmt"
	"github.com/nickwells/golem/location"
	"testing"
)

func TestHasLoop(t *testing.T) {
	l1 := location.New("f1")
	l2 := location.New("f2")
	l3 := location.New("f3")

	type hasLoopTest struct {
		name             string
		chain            location.LocChain
		expectedHasLoop  bool
		expectedLoopDesc string
	}

	testCases := []hasLoopTest{
		{
			name:             "f2",
			chain:            location.LocChain{*l1, *l2, *l3},
			expectedHasLoop:  true,
			expectedLoopDesc: "f1:0 --> f2:0 ==> f3:0 ==> Back to f2",
		},
		{
			name:             "notInIncludeChain",
			chain:            location.LocChain{*l1, *l2, *l3},
			expectedHasLoop:  false,
			expectedLoopDesc: "",
		},
	}

	for _, hlt := range testCases {
		testID := fmt.Sprintf("HasLoop('%s', %s)",
			hlt.name, hlt.chain)
		hasLoop, loopDesc := hlt.chain.HasLoop(hlt.name)
		if hasLoop != hlt.expectedHasLoop {
			t.Logf("%s: failed\n", testID)
			t.Logf("\t : expected: %v\n", hlt.expectedHasLoop)
			t.Errorf("\t :      got: %v\n", hasLoop)
		}
		if loopDesc != hlt.expectedLoopDesc {
			t.Logf("%s: failed\n", testID)
			t.Logf("\t : expected loop description: %v\n", hlt.expectedLoopDesc)
			t.Errorf("\t :                       got: %v\n", loopDesc)
		}
	}
}
