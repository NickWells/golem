package location

import "testing"

func TestError(t *testing.T) {
	loc := New("test")
	testCases := []struct {
		testName string
		e        Err
		expStr   string
	}{
		{
			testName: "nil",
			expStr:   " (at :0)",
		},
		{
			testName: "no Loc",
			e: Err{
				Msg: "no Loc",
			},
			expStr: "no Loc (at :0)",
		},
		{
			testName: "with Loc",
			e: Err{
				Msg: "Msg",
				Loc: *loc,
			},
			expStr: "Msg (at test:0)",
		},
	}

	for i, tc := range testCases {
		e := tc.e
		if s := e.Error(); s != tc.expStr {
			t.Errorf(
				"test %d: %s : e.Error() should have returned '%s' not '%s'",
				i, tc.testName, tc.expStr, s)
		}
	}

}
