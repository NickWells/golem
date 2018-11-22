package filecheck

import (
	"os"
	"regexp"
	"strings"
	"time"
)

// InfoChecker is a function that takes a FileInfo and checks it returning
// true if the check passes and false otherwise
type InfoChecker func(os.FileInfo) bool

// ICNot returns an InfoChecker that will invert the result of the passed
// InfoChecker
func ICNot(ic InfoChecker) InfoChecker {
	return func(fi os.FileInfo) bool {
		return !ic(fi)
	}
}

// ICOr returns an InfoChecker that will return true if any one of the
// InfoChecker's returns true. It will return true as soon as any of them
// returns true and so it is not guaranteed that they will all be called
func ICOr(ics ...InfoChecker) InfoChecker {
	return func(fi os.FileInfo) bool {
		for _, ic := range ics {
			if ic(fi) {
				return true
			}
		}
		return false
	}
}

// ICAnd returns an InfoChecker that will return true if all of the
// InfoChecker's returns true. It will return false as soon as any of them
// returns false and so it is not guaranteed that they will all be called
func ICAnd(ics ...InfoChecker) InfoChecker {
	return func(fi os.FileInfo) bool {
		for _, ic := range ics {
			if !ic(fi) {
				return false
			}
		}
		return true
	}
}

// ICNameMatches returns an InfoChecker that will check that the file name
// matches the passed regular expression
func ICNameMatches(re *regexp.Regexp) InfoChecker {
	return func(fi os.FileInfo) bool {
		return re.MatchString(fi.Name())
	}
}

// ICNameEquals returns an InfoChecker that will check that the file name
// equals the passed value
func ICNameEquals(n string) InfoChecker {
	return func(fi os.FileInfo) bool {
		return fi.Name() == n
	}
}

// ICNameHasPrefix returns an InfoChecker that will check that the file name
// has the supplied prefix
func ICNameHasPrefix(prefix string) InfoChecker {
	return func(fi os.FileInfo) bool {
		return strings.HasPrefix(fi.Name(), prefix)
	}
}

// ICNameHasSuffix returns an InfoChecker that will check that the file name
// has the supplied suffix
func ICNameHasSuffix(suffix string) InfoChecker {
	return func(fi os.FileInfo) bool {
		return strings.HasSuffix(fi.Name(), suffix)
	}
}

// ICIsDir returns an InfoChecker that will check that the file is a
// directory
func ICIsDir(fi os.FileInfo) bool {
	return fi.IsDir()
}

// ICIsRegularFile returns an InfoChecker that will check that the file is
// a regular file
func ICIsRegularFile(fi os.FileInfo) bool {
	return fi.Mode()&os.ModeType == 0
}

// ICSizeLT returns an InfoChecker that will check that the file size is
// less than the parameter
func ICSizeLT(size int64) InfoChecker {
	return func(fi os.FileInfo) bool {
		return fi.Size() < size
	}
}

// ICSizeEQ returns true if the file size is equal to the parameter
func ICSizeEQ(size int64) InfoChecker {
	return func(fi os.FileInfo) bool {
		return fi.Size() == size
	}
}

// ICSizeGT returns an InfoChecker that will check that the file size is
// greater than the parameter
func ICSizeGT(size int64) InfoChecker {
	return func(fi os.FileInfo) bool {
		return fi.Size() > size
	}
}

// ICModTimeBefore returns true if the file modification time is before the
// parameter
func ICModTimeBefore(t time.Time) InfoChecker {
	return func(fi os.FileInfo) bool {
		return fi.ModTime().Before(t)
	}
}

// ICModTimeEqual returns an InfoChecker that will check that the file
// modification time is equal to the parameter
func ICModTimeEqual(t time.Time) InfoChecker {
	return func(fi os.FileInfo) bool {
		return fi.ModTime().Equal(t)
	}
}

// ICModTimeAfter returns an InfoChecker that will check that the file
// modification time is after the parameter
func ICModTimeAfter(t time.Time) InfoChecker {
	return func(fi os.FileInfo) bool {
		return fi.ModTime().After(t)
	}
}

// DirEntries reads the directory and returns the FileInfo details for each
// entry for which all the checks return true. Any errors are also returned
func DirEntries(dirName string, checks ...InfoChecker) ([]os.FileInfo, error) {
	d, err := os.Open(dirName)
	if err != nil {
		return []os.FileInfo{}, err
	}

	defer d.Close()

	allInfo, err := d.Readdir(0)
	info := make([]os.FileInfo, 0, len(allInfo))
	for _, fi := range allInfo {
		isGood := true

		for _, ic := range checks {
			if !ic(fi) {
				isGood = false
				break
			}
		}

		if isGood {
			info = append(info, fi)
		}
	}
	return info, err
}

// CountDirEntries returns the number of entries in the directory that match
// the supplied checks (if any) and any errors detected
func CountDirEntries(dirName string, checks ...InfoChecker) (int, error) {
	info, err := DirEntries(dirName, checks...)
	return len(info), err
}
