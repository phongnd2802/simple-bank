package db

import (
	"github.com/jackc/pgx/v5"
	"github.com/lib/pq"
)


var ErrUniqueViolation = &pq.Error{
	Code: pq.ErrorCode("23505"),
}

var ErrRecordNotFound = pgx.ErrNoRows