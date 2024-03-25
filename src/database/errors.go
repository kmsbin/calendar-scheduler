package database

import "github.com/lib/pq"

const (
	UniqueViolationErr = pq.ErrorCode("23505")
)

func IsPqErrorCode(err error, errcode pq.ErrorCode) bool {
	if pgerr, ok := err.(*pq.Error); ok {
		return pgerr.Code == errcode
	}
	return false
}
