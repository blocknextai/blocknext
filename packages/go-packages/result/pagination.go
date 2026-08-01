package result

import (
	"math"
)

const (
	defaultLimit = 10
	maxLimit     = 100
	maxOffset    = math.MaxInt32
)

// PaginationRequest holds the requested pagination offset and limit.
type PaginationRequest struct {
	Offset int `json:"offset" query:"offset"`
	Limit  int `json:"limit" query:"limit"`
}

// Normalize clamps the request's offset and limit to their valid bounds in
// place and returns the normalized value, so it is correct both as a statement
// and as an expression.
func (p *PaginationRequest) Normalize() PaginationRequest {
	*p = NewPaginationRequest(p.Offset, p.Limit)
	return *p
}

// Pagination describes the computed pagination state of a result set.
type Pagination struct {
	Total      int64 `json:"total"`
	Page       int   `json:"page"`
	TotalPages int   `json:"totalPages"`
	Limit      int   `json:"limit"`
	Offset     int   `json:"offset"`
	HasNext    bool  `json:"hasNext"`
	HasPrev    bool  `json:"hasPrev"`
}

// NewPaginationRequest returns a PaginationRequest with offset and limit clamped
// to their valid ranges, applying default and maximum limits as needed.
func NewPaginationRequest(offset int, limit int) PaginationRequest {
	return PaginationRequest{
		Offset: normalizeOffset(offset),
		Limit:  normalizeLimit(limit),
	}
}

// NewPagination computes a Pagination from the total item count, offset, and
// limit. Offset, limit and total are clamped to the same bounds
// NewPaginationRequest applies, so callers cannot bypass validation by reaching
// this function directly, and every field is computed without integer overflow.
//
// HasNext and HasPrev are deliberately conservative so that a client paging loop
// always terminates: HasNext is false once offset reaches maxOffset, even when
// more rows exist, because no further offset can be requested past that ceiling;
// and HasPrev is false whenever total is zero.
func NewPagination(total int64, offset int, limit int) Pagination {
	offset = normalizeOffset(offset)
	limit = normalizeLimit(limit)
	total = max(total, 0)

	off := int64(offset)
	lim := int64(limit)

	totalPages := total / lim
	if total%lim != 0 {
		totalPages++
	}
	totalPages = min(totalPages, math.MaxInt)

	page := min(off/lim+1, math.MaxInt)

	return Pagination{
		Total:      total,
		Page:       int(page),
		TotalPages: int(totalPages),
		Offset:     offset,
		Limit:      limit,
		HasNext:    total-off > lim && offset < maxOffset,
		HasPrev:    offset > 0 && total > 0,
	}
}

func normalizeOffset(offset int) int {
	return min(max(offset, 0), maxOffset)
}

func normalizeLimit(limit int) int {
	if limit < 1 {
		return defaultLimit
	}

	return min(limit, maxLimit)
}
