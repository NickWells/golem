package filecheck_test

import (
	"github.com/NickWells/golem/filecheck"
	"os"
	"path"
	"regexp"
	"testing"
	"time"
)

func TestCountDirEntries(t *testing.T) {
	regexpStr := "[a-zA-Z]*1$"
	goodDirName := path.Join("testdata", "countEntries")
	badDirName := path.Join("testdata", "NoSuchDir")
	modTime := time.Date(2009, time.November, 10, 23, 0, 0, 0, time.UTC)
	modTimePlus1 := time.Date(2009, time.November, 10, 23, 1, 0, 0, time.UTC)

	testCases := []struct {
		testName      string
		checks        []filecheck.InfoChecker
		errExpected   bool
		countExpected int
		dirName       string
		modTime       time.Time
	}{
		{
			testName:      "bad directory: " + badDirName,
			errExpected:   true,
			countExpected: 0,
			dirName:       badDirName,
		},
		{
			testName:      "all entries",
			errExpected:   false,
			countExpected: 6,
		},
		{
			testName: "all files",
			checks: []filecheck.InfoChecker{
				filecheck.ICIsRegularFile,
			},
			errExpected:   false,
			countExpected: 4,
		},
		{
			testName: "matching regex: " + regexpStr,
			checks: []filecheck.InfoChecker{
				filecheck.ICNameMatches(regexp.MustCompile(regexpStr)),
			},
			errExpected:   false,
			countExpected: 2,
		},
		{
			testName: "name == file1",
			checks: []filecheck.InfoChecker{
				filecheck.ICNameEquals("file1"),
			},
			errExpected:   false,
			countExpected: 1,
		},
		{
			testName: "name != file1",
			checks: []filecheck.InfoChecker{
				filecheck.ICNot(filecheck.ICNameEquals("file1")),
			},
			errExpected:   false,
			countExpected: 5,
		},
		{
			testName: "size == 0",
			checks: []filecheck.InfoChecker{
				filecheck.ICSizeEQ(0),
			},
			errExpected:   false,
			countExpected: 3,
		},
		{
			testName: "size > 0",
			checks: []filecheck.InfoChecker{
				filecheck.ICSizeGT(0),
			},
			errExpected:   false,
			countExpected: 3,
		},
		{
			testName: "size < 100",
			checks: []filecheck.InfoChecker{
				filecheck.ICSizeLT(100),
			},
			errExpected:   false,
			countExpected: 4,
		},
		{
			testName: "files with suffix 1",
			checks: []filecheck.InfoChecker{
				filecheck.ICIsRegularFile,
				filecheck.ICNameHasSuffix("1"),
			},
			errExpected:   false,
			countExpected: 1,
		},
		{
			testName: "entries with prefix d",
			checks: []filecheck.InfoChecker{
				filecheck.ICNameHasPrefix("d"),
			},
			errExpected:   false,
			countExpected: 1,
		},
		{
			testName: "files with suffix 1 or directories",
			checks: []filecheck.InfoChecker{
				filecheck.ICOr(
					filecheck.ICAnd(
						filecheck.ICIsRegularFile,
						filecheck.ICNameHasSuffix("1"),
					),
					filecheck.ICIsDir,
				),
			},
			errExpected:   false,
			countExpected: 3,
		},
		{
			testName: "entries with modTime before " + modTimePlus1.String(),
			checks: []filecheck.InfoChecker{
				filecheck.ICModTimeBefore(modTimePlus1),
			},
			errExpected:   false,
			countExpected: 1,
			modTime:       modTime,
		},
		{
			testName: "entries with modTime equals " + modTime.String(),
			checks: []filecheck.InfoChecker{
				filecheck.ICModTimeEqual(modTime),
			},
			errExpected:   false,
			countExpected: 1,
			modTime:       modTime,
		},
		{
			testName: "entries with modTime after " + modTimePlus1.String(),
			checks: []filecheck.InfoChecker{
				filecheck.ICModTimeAfter(modTimePlus1),
			},
			errExpected:   false,
			countExpected: 5,
			modTime:       modTime,
		},
	}

	for i, tc := range testCases {
		dirName := goodDirName
		if tc.dirName != "" {
			dirName = tc.dirName
		}

		if !tc.modTime.IsZero() {
			filename := path.Join(goodDirName, "file1")
			err := os.Chtimes(filename, tc.modTime, tc.modTime)
			if err != nil {
				t.Logf("test %d: %s :\n", i, tc.testName)
				t.Logf("\t: while setting access and mod times on: %s\n",
					filename)
				t.Errorf("\t: unexpected error : %s\n", err)
			}
		}

		n, err := filecheck.CountDirEntries(dirName, tc.checks...)

		if err != nil {
			if !tc.errExpected {
				t.Logf("test %d: %s :\n", i, tc.testName)
				t.Logf("\t: in dir: %s\n", dirName)
				t.Errorf("\t: unexpected error: %s\n", err)
			}
		} else {
			if tc.errExpected {
				t.Logf("test %d: %s :\n", i, tc.testName)
				t.Logf("\t: in dir: %s\n", dirName)
				t.Errorf("\t: an error was expected but not seen\n")
			}
		}
		if n != tc.countExpected {
			t.Logf("test %d: %s :\n", i, tc.testName)
			t.Logf("\t: in dir: %s\n", dirName)
			t.Errorf("\t: expected count: %d got: %d\n", tc.countExpected, n)
		}
	}
}
