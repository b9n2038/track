package rating

import "errors"

var (
	ErrInvalidRating = errors.New("invalid rating value")
	ErrInvalidDate   = errors.New("invalid date")
	ErrNotFound      = errors.New("rating not found")
)
