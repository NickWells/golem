package param_test

import (
	"fmt"
	"github.com/nickwells/golem/filecheck"
	"github.com/nickwells/golem/param"
	"github.com/nickwells/golem/param/paramset"
	"github.com/nickwells/golem/param/psetter"
	"testing"
)

var CFValExample1 bool
var CFValExample2 int64

func TestConfigFile(t *testing.T) {
	CFValExample1 = false
	CFValExample2 = 0

	ps, err := paramset.NewNoHelpNoExitNoErrRpt(CFAddParams1, CFAddParams2)
	if err != nil {
		t.Fatal("TestConfigFile : couldn't construct the ParamSet: ", err)
	}
	const mustExistDoes = "./testdata/config.test"
	ps.SetConfigFile(mustExistDoes, filecheck.MustExist)

	const mustExistDoesNot = "./testdata/config.nosuch"
	ps.AddConfigFile(mustExistDoesNot, filecheck.MustExist)

	const mayExistDoes = "./testdata/config.opt"
	ps.AddConfigFile(mayExistDoes, filecheck.Optional)

	const mayExistDoesNot = "./testdata/config.opt.nosuch"
	ps.AddConfigFile(mayExistDoesNot, filecheck.Optional)

	ps.Parse([]string{})

	if errs, ok := ps.Errors()["config file: "+mustExistDoes]; ok {
		t.Logf("Unexpected problem with config file:\n")
		t.Errorf("\t: %s\n", mustExistDoes)
		t.Errorf("\t: got: %v\n", errs)
	}

	if _, ok := ps.Errors()["config file: "+mustExistDoesNot]; !ok {
		t.Logf("A problem was expected with missing, must-exist config file:\n")
		t.Errorf("\t: %s\n", mustExistDoesNot)
		t.Errorf("\t: none found\n")
	}

	if errs, ok := ps.Errors()["config file: "+mayExistDoes]; ok {
		t.Logf("Unexpected problem with config file:\n")
		t.Errorf("\t: %s\n", mayExistDoes)
		t.Errorf("\t: got: %v\n", errs)
	}

	if errs, ok := ps.Errors()["config file: "+mayExistDoesNot]; ok {
		t.Logf("Unexpected problem with missing, optional config file:\n")
		t.Errorf("\t: %s\n", mayExistDoesNot)
		t.Errorf("\t: got: %v\n", errs)
	}

	if CFValExample2 != 5 {
		t.Errorf("CFValExample2 should be 5 but is: %d\n", CFValExample2)
	}
}

// CFAddParams1 will set the "example1" parameter in the ParamSet
func CFAddParams1(ps *param.ParamSet) error {
	ps.Add("example1",
		psetter.BoolSetter{Value: &CFValExample1},
		"here is where you would describe the parameter",
		param.AltName("e1"))

	return nil
}

// CFAddParams2 will set the "example2" parameter in the ParamSet
func CFAddParams2(ps *param.ParamSet) error {
	ps.Add("example2",
		psetter.Int64Setter{Value: &CFValExample2},
		"the description of the parameter",
		param.AltName("e2"))

	return nil
}

var groupCFName1 = "grp1"
var groupCFName2 = "grp2"
var paramInt1 int64
var paramInt2 int64
var paramBool1 bool
var paramBool2 bool

type expVals struct {
	pi1Val int64
	pi2Val int64
	pb1Val bool
	pb2Val bool
}

func TestGroupConfigFile(t *testing.T) {
	configFileNameA := "testdata/groupConfigFile.A"
	configFileNameB := "testdata/groupConfigFile.B"
	configFileNameC := "testdata/groupConfigFile.C"
	configFileNameNonesuch := "testdata/groupConfigFile.nonesuch"
	testCases := []struct {
		name         string
		gName        string
		fileName     string
		check        filecheck.Exists
		errsExpected map[string][]string
		valsExpected expVals
	}{
		{
			name:         "all good - file must exist and does",
			gName:        groupCFName1,
			fileName:     configFileNameA,
			check:        filecheck.MustExist,
			valsExpected: expVals{pi1Val: 42, pb1Val: true},
		},
		{
			name:     "config file exists but has an unknown parameter",
			gName:    groupCFName1,
			fileName: configFileNameB,
			check:    filecheck.MustExist,
			errsExpected: map[string][]string{
				"unknown-param": []string{
					"this is not a parameter of this program",
					"group-specific parameter config file",
					configFileNameB,
				},
			},
			valsExpected: expVals{pi1Val: 42, pb1Val: true},
		},
		{
			name:     "config file exists but has a parameter from another group",
			gName:    groupCFName1,
			fileName: configFileNameC,
			check:    filecheck.MustExist,
			errsExpected: map[string][]string{
				"pi2": []string{
					"this parameter is not a member of group: " + groupCFName1,
					"group-specific parameter config file",
					configFileNameC,
				},
			},
			valsExpected: expVals{pi1Val: 42, pb1Val: true},
		},
		{
			name:     "missing file",
			gName:    groupCFName1,
			fileName: configFileNameNonesuch,
			check:    filecheck.MustExist,
			errsExpected: map[string][]string{
				"config file: " + configFileNameNonesuch: []string{
					"no such file or directory",
					configFileNameNonesuch,
				},
			},
		},
	}

	for i, tc := range testCases {
		testID := fmt.Sprintf("%d: %s", i, tc.name)

		ps, err := paramset.NewNoHelpNoExitNoErrRpt()
		if err != nil {
			t.Fatal(testID, " : couldn't construct the ParamSet: ", err)
		}
		addParamsForGroupCF(ps)
		ps.AddGroupConfigFile(tc.gName, tc.fileName, tc.check)

		resetParamVals()
		errMap := ps.Parse([]string{})

		errMapCheck(t, testID, errMap, tc.errsExpected)
		valsCheck(t, testID, tc.valsExpected)
	}

}

// resetParamVals resets the param values to their initial state
func resetParamVals() {
	paramInt1 = 0
	paramInt2 = 0
	paramBool1 = false
	paramBool2 = false
}

// valsCheck checks that the values match the expected values
func valsCheck(t *testing.T, testID string, vals expVals) {
	t.Helper()

	showId := true
	if paramInt1 != vals.pi1Val {
		t.Logf("test: %s:\n", testID)
		showId = false
		t.Errorf("\t: unexpected values: paramInt1 = %d, should be %d\n",
			paramInt1, vals.pi1Val)
	}

	if paramInt2 != vals.pi2Val {
		if showId {
			t.Logf("test: %s:\n", testID)
			showId = false
		}
		t.Errorf("\t: unexpected values: paramInt2 = %d, should be %d\n",
			paramInt2, vals.pi2Val)
	}

	if paramBool1 != vals.pb1Val {
		if showId {
			t.Logf("test: %s:\n", testID)
			showId = false
		}
		t.Errorf("\t: unexpected values: paramBool1 = %v, should be %v\n",
			paramBool1, vals.pb1Val)
	}

	if paramBool2 != vals.pb2Val {
		if showId {
			t.Logf("test: %s:\n", testID)
			showId = false
		}
		t.Errorf("\t: unexpected values: paramBool2 = %v, should be %v\n",
			paramBool2, vals.pb2Val)
	}
}

func addParamsForGroupCF(ps *param.ParamSet) {
	ps.SetGroupDescription(groupCFName1, "blah blah blah - 1")
	ps.SetGroupDescription(groupCFName2, "blah blah blah - 2")
	ps.Add("pi1", psetter.Int64Setter{Value: &paramInt1},
		"param int val 1",
		param.GroupName(groupCFName1))
	ps.Add("pi2", psetter.Int64Setter{Value: &paramInt2},
		"param int val 2",
		param.GroupName(groupCFName2))
	ps.Add("pb1", psetter.BoolSetter{Value: &paramBool1},
		"param bool val 1",
		param.GroupName(groupCFName1))
	ps.Add("pb2", psetter.BoolSetter{Value: &paramBool2},
		"param bool val 2",
		param.GroupName(groupCFName2))
}
