package filecheck

import (
	"errors"
	"fmt"
	"github.com/NickWells/golem/check"
	"os"
)

// Exists records whether the file-system object should exist or
// not. In each case the check is only valid at the time the check is
// made and so any code using this should be aware of this
type Exists uint

// Optional indicates that no existence check should be made
//
// MustExist indicates that the object must exist
//
// MustNotExist indicates that the object must not exist
const (
	Optional Exists = iota
	MustExist
	MustNotExist
)

// FSObjTypeDontCare means we don't care what type the file system object has
//
// FSObjTypeRegularFile indicates that the object should be a regular file
//
// FSObjTypeDirectory indicates that the object should be a directory
//
// FSObjTypeSymlink indicates that the object should be a symlink
//
// FSObjTypeNamedPipe indicates that the object should be a named pipe
//
// FSObjTypeSocket indicates that the object should be a socket
//
// FSObjTypeDevice indicates that the object should be a device
const (
	FSObjTypeDontCare    os.FileMode = os.ModeType
	FSObjTypeRegularFile             = 0
	FSObjTypeDirectory               = os.ModeDir
	FSObjTypeSymlink                 = os.ModeSymlink
	FSObjTypeNamedPipe               = os.ModeNamedPipe
	FSObjTypeSocket                  = os.ModeSocket
	FSObjTypeDevice                  = os.ModeDevice
)

// PermChecker will be called with just the permission bits of the File Mode
// bits set.
type PermChecker func(os.FileMode) error

// PermCheckEq returns a PermChecker which will return a non-nil error if the
// permission bits are not equal to the passed value
func PermCheckEq(p os.FileMode) PermChecker {
	return func(filePerms os.FileMode) error {
		if filePerms == p {
			return nil
		}
		return fmt.Errorf(
			"bad permissions: should be %04o but we have %04o",
			p, filePerms)
	}
}

// PermCheckHasAll returns a PermChecker which will return a non-nil error if
// the permission bits do not have all of the bits in the passed value set
func PermCheckHasAll(p os.FileMode) PermChecker {
	return func(filePerms os.FileMode) error {
		if (filePerms & p) == p {
			return nil
		}
		return fmt.Errorf(
			"bad permissions: all of %04o should be set but we have %04o",
			p, filePerms)
	}
}

// PermCheckHasNone returns a PermChecker which will return a non-nil error if
// the permission bits have any of the bits in the passed value set
func PermCheckHasNone(p os.FileMode) PermChecker {
	return func(filePerms os.FileMode) error {
		if (filePerms & p) == 0 {
			return nil
		}
		return fmt.Errorf(
			"bad permissions: none of %04o should be set but we have %04o",
			p, filePerms)
	}
}

// ExpectedStatus records the expectations of the file-system object
type ExpectedStatus struct {
	ObjectType   os.FileMode
	PermCheckers []PermChecker
	SizeCheckers []check.Int64
	Existence    Exists
}

// StatusCheck checks that the file system object called 'name' satisfies
// the constraints. it returns a non-nil error if the constraint is not
// met. Note that if the file does not exist and it is not expected to
// exist then no further checks are performed (this may be obvious to you)
func (es ExpectedStatus) StatusCheck(name string) error {
	var info os.FileInfo
	var err error
	if es.ObjectType == FSObjTypeSymlink {
		info, err = os.Lstat(name)
	} else {
		info, err = os.Stat(name)
	}

	if os.IsNotExist(err) {
		if es.Existence == MustExist {
			return errors.New("path: '" + name +
				"' does not exist but should")
		}
		return nil
	}

	if es.Existence == MustNotExist {
		return errors.New("path: '" + name +
			"' exists but shouldn't")
	}

	if err != nil {
		return errors.New("path: '" + name +
			"' error: " + err.Error())
	}

	if es.ObjectType != FSObjTypeDontCare {
		typeBits := info.Mode() & os.ModeType
		if typeBits != es.ObjectType {
			return fmt.Errorf("path: '%s' should be a %s but is a %s",
				name,
				ObjectTypeAsString(es.ObjectType),
				ObjectTypeAsString(typeBits))
		}
	}

	// TODO: add handling for not following symlinks

	for _, pc := range es.PermCheckers {
		if err := pc(info.Mode() & os.ModePerm); err != nil {
			return err
		}
	}

	for _, sc := range es.SizeCheckers {
		if err := sc(info.Size()); err != nil {
			return err
		}
	}
	return nil
}

// ObjectTypeAsString returns a string representation of the object type
// information
func ObjectTypeAsString(objType os.FileMode) string {
	switch objType {
	case FSObjTypeDontCare:
		return "any type"
	case FSObjTypeRegularFile:
		return "regular file"
	case FSObjTypeDirectory:
		return "directory"
	case FSObjTypeSymlink:
		return "symbolic link"
	case FSObjTypeNamedPipe:
		return "named pipe"
	case FSObjTypeSocket:
		return "socket"
	case FSObjTypeDevice:
		return "device"
	}
	return fmt.Sprintf("unknown: %v", objType)
}

func (es ExpectedStatus) String() string {
	if es.Existence == MustNotExist {
		return "The filesystem object must not exist"
	}

	var rval, prefix string

	if es.Existence == MustExist {
		rval = "The filesystem object must exist"
		prefix = " and"
	} else if es.Existence == Optional {
		rval = "The filesystem object need not exist"
		prefix = " but if it does it"
	}
	if es.ObjectType != FSObjTypeDontCare {
		rval += prefix + " must be a " + ObjectTypeAsString(es.ObjectType)
	}

	return rval
}
