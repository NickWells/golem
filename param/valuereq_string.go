// Code generated by "stringer -type=ValueReq"; DO NOT EDIT.

package param

import "strconv"

const _ValueReq_name = "MandatoryOptionalNone"

var _ValueReq_index = [...]uint8{0, 9, 17, 21}

func (i ValueReq) String() string {
	if i < 0 || i >= ValueReq(len(_ValueReq_index)-1) {
		return "ValueReq(" + strconv.FormatInt(int64(i), 10) + ")"
	}
	return _ValueReq_name[_ValueReq_index[i]:_ValueReq_index[i+1]]
}
