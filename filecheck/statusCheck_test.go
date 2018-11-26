package filecheck_test

import (
	"fmt"
	"os"
	"testing"

	"github.com/nickwells/golem/check"
	"github.com/nickwells/golem/filecheck"
)

func TestStatusCheck(t *testing.T) {
	os.Chmod("testdata/IsAFile.PBits0600", 0600) // force the file mode

	testCases := []struct {
		testName    string
		name        string
		es          filecheck.ExpectedStatus
		expectErr   bool
		expectedStr string
	}{
		{
			testName: "doesn't exist",
			name:     "testdata/nonesuch",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustNotExist,
				ObjectType: filecheck.FSObjTypeDontCare,
			},
			expectErr:   false,
			expectedStr: "The filesystem object must not exist",
		},
		{
			testName: "need not exist and doesn't",
			name:     "testdata/nonesuch",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.Optional,
				ObjectType: filecheck.FSObjTypeDontCare,
			},
			expectErr:   false,
			expectedStr: "The filesystem object need not exist",
		},
		{
			testName: "need not exist and doesn't - must be a named pipe",
			name:     "testdata/nonesuch",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.Optional,
				ObjectType: filecheck.FSObjTypeNamedPipe,
			},
			expectErr: false,
			expectedStr: "The filesystem object need not exist" +
				" but if it does it must be a named pipe",
		},
		{
			testName: "need not exist and doesn't - must be a socket",
			name:     "testdata/nonesuch",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.Optional,
				ObjectType: filecheck.FSObjTypeSocket,
			},
			expectErr: false,
			expectedStr: "The filesystem object need not exist" +
				" but if it does it must be a socket",
		},
		{
			testName: "need not exist and doesn't - must be a device",
			name:     "testdata/nonesuch",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.Optional,
				ObjectType: filecheck.FSObjTypeDevice,
			},
			expectErr: false,
			expectedStr: "The filesystem object need not exist" +
				" but if it does it must be a device",
		},
		{
			testName: "need not exist and does",
			name:     "testdata/IsAFile",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.Optional,
				ObjectType: filecheck.FSObjTypeDontCare,
			},
			expectErr:   false,
			expectedStr: "The filesystem object need not exist",
		},
		{
			testName: "need not exist and doesn't but should be a file",
			name:     "testdata/nonesuch",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.Optional,
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: false,
			expectedStr: "The filesystem object need not exist" +
				" but if it does it must be a regular file",
		},
		{
			testName: "need not exist and does but should be a file",
			name:     "testdata/IsAFile",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.Optional,
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: false,
			expectedStr: "The filesystem object need not exist" +
				" but if it does it must be a regular file",
		},
		{
			testName: "need not exist and does but should be a file",
			name:     "testdata/IsNotAFile",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.Optional,
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: true,
			expectedStr: "The filesystem object need not exist" +
				" but if it does it must be a regular file",
		},
		{
			testName: "dir - exists and should",
			name:     "testdata/IsADirectory",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustExist,
				ObjectType: filecheck.FSObjTypeDirectory,
			},
			expectErr: false,
			expectedStr: "The filesystem object must exist" +
				" and must be a directory",
		},
		{
			testName: "dir - exists and shouldn't",
			name:     "testdata/IsADirectory",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustNotExist,
				ObjectType: filecheck.FSObjTypeDirectory,
			},
			expectErr:   true,
			expectedStr: "The filesystem object must not exist",
		},
		{
			testName: "dir - exists but isn't a dir",
			name:     "testdata/IsNotADirectory",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustExist,
				ObjectType: filecheck.FSObjTypeDirectory,
			},
			expectErr: true,
			expectedStr: "The filesystem object must exist" +
				" and must be a directory",
		},
		{
			testName: "file - exists and should",
			name:     "testdata/IsAFile",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustExist,
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: false,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file - exists and shouldn't",
			name:     "testdata/IsAFile",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustNotExist,
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr:   true,
			expectedStr: "The filesystem object must not exist",
		},
		{
			testName: "file - exists but isn't a file",
			name:     "testdata/IsNotAFile",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustExist,
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: true,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file (symlink) - exists and should",
			name:     "testdata/IsASymlinkToAFile",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustExist,
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: false,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file (symlink) - exists but is a link to nothing",
			name:     "testdata/IsASymlinkToNothing",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustExist,
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: true,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "symlink - exists but is a link to nothing",
			name:     "testdata/IsASymlinkToNothing",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustExist,
				ObjectType: filecheck.FSObjTypeSymlink,
			},
			expectErr: false,
			expectedStr: "The filesystem object must exist" +
				" and must be a symbolic link",
		},
		{
			testName: "symlink - exists but is not a symlink",
			name:     "testdata/IsNotASymlink",
			es: filecheck.ExpectedStatus{
				Existence:  filecheck.MustExist,
				ObjectType: filecheck.FSObjTypeSymlink,
			},
			expectErr: true,
			expectedStr: "The filesystem object must exist" +
				" and must be a symbolic link",
		},
		{
			testName: "file - perms equal 0664",
			name:     "testdata/IsAFile",
			es: filecheck.ExpectedStatus{
				Existence: filecheck.MustExist,
				PermCheckers: []filecheck.PermChecker{
					filecheck.PermCheckEq(0644),
				},
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: false,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file - perms don't equal 0664",
			name:     "testdata/IsAFile.PBits0600",
			es: filecheck.ExpectedStatus{
				Existence: filecheck.MustExist,
				PermCheckers: []filecheck.PermChecker{
					filecheck.PermCheckEq(0644),
				},
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: true,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file - perms don't have 0001",
			name:     "testdata/IsAFile.PBits0600",
			es: filecheck.ExpectedStatus{
				Existence: filecheck.MustExist,
				PermCheckers: []filecheck.PermChecker{
					filecheck.PermCheckHasNone(0001),
				},
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: false,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file - perms have 0200",
			name:     "testdata/IsAFile.PBits0600",
			es: filecheck.ExpectedStatus{
				Existence: filecheck.MustExist,
				PermCheckers: []filecheck.PermChecker{
					filecheck.PermCheckHasAll(0200),
				},
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: false,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file - perms have 0200 and don't have 0001",
			name:     "testdata/IsAFile.PBits0600",
			es: filecheck.ExpectedStatus{
				Existence: filecheck.MustExist,
				PermCheckers: []filecheck.PermChecker{
					filecheck.PermCheckHasAll(0200),
					filecheck.PermCheckHasNone(0001),
				},
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: false,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file - perms have 0600 and don't have 0077",
			name:     "testdata/IsAFile.PBits0600",
			es: filecheck.ExpectedStatus{
				Existence: filecheck.MustExist,
				PermCheckers: []filecheck.PermChecker{
					filecheck.PermCheckHasAll(0600),
					filecheck.PermCheckHasNone(0077),
				},
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: false,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file - perms have 0200 and have 0001",
			name:     "testdata/IsAFile.PBits0600",
			es: filecheck.ExpectedStatus{
				Existence: filecheck.MustExist,
				PermCheckers: []filecheck.PermChecker{
					filecheck.PermCheckHasAll(0200),
					filecheck.PermCheckHasAll(0001),
				},
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: true,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file - perms don't have 0200 and have 0001",
			name:     "testdata/IsAFile.PBits0600",
			es: filecheck.ExpectedStatus{
				Existence: filecheck.MustExist,
				PermCheckers: []filecheck.PermChecker{
					filecheck.PermCheckHasNone(0200),
					filecheck.PermCheckHasAll(0001),
				},
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: true,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file - size must be 0 and is",
			name:     "testdata/IsAFile",
			es: filecheck.ExpectedStatus{
				Existence: filecheck.MustExist,
				SizeCheckers: []check.Int64{
					check.Int64EQ(0),
				},
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: false,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
		{
			testName: "file - size must be > 0 and is not",
			name:     "testdata/IsAFile",
			es: filecheck.ExpectedStatus{
				Existence: filecheck.MustExist,
				SizeCheckers: []check.Int64{
					check.Int64GT(0),
				},
				ObjectType: filecheck.FSObjTypeRegularFile,
			},
			expectErr: true,
			expectedStr: "The filesystem object must exist" +
				" and must be a regular file",
		},
	}

	for i, tc := range testCases {
		testID := fmt.Sprintf("test %d: %s", i, tc.testName)
		err := tc.es.StatusCheck(tc.name)
		if tc.expectErr {
			if err == nil {
				t.Logf("%s:\n", testID)
				t.Errorf("\t: an error was expected but none was seen\n")
			}
		} else {
			if err != nil {
				t.Logf("%s:\n", testID)
				t.Errorf("\t: unexpected error: %s\n", err)
			}
		}
		if tc.es.String() != tc.expectedStr {
			t.Logf("%s:\n", testID)
			t.Logf("\t: description was expected to be: %s\n", tc.expectedStr)
			t.Logf("\t:                        but was: %s\n", tc.es.String())
			t.Errorf("\t: unexpected string representation\n")
		}
	}

}
