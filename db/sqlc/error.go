package db

import "github.com/lib/pq"


var ErrUniqueViolation = &pq.Error{
	Code: pq.ErrorCode("23505"),
}