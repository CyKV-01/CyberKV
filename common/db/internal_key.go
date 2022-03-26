package db

import (
	. "github.com/yah01/CyberKV/common"
)

type InternalKey struct {
	key string
	ts  TimeStamp
}

var (
	DefaultUserKeyComparator = func(a, b string) int {
		if a < b {
			return -1
		} else if a > b {
			return 1
		}

		return 0
	}
)

func NewInternalKey(key string, ts TimeStamp) InternalKey {

	return InternalKey{key, ts}
}

func (key InternalKey) Compare(other InternalKey) int {
	a, b := key.UserKey(), other.UserKey()
	res := DefaultUserKeyComparator(a, b)
	if res != 0 {
		return res
	}

	tsa, tsb := key.GetTimeStamp(), other.GetTimeStamp()
	if tsa > tsb {
		return -1
	} else if tsa < tsb {
		return 1
	}

	return 0
}

func (key InternalKey) UserKey() string {
	return key.key
}

func (key InternalKey) GetTimeStamp() TimeStamp {
	return key.ts
}
