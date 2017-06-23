package models

// NewCollectionParams creates a new instance of CollectionParams.
func NewCollectionParams(cursor string, limit uint8) (params *CollectionParams) {

	if limit == 0 {
		limit = 25
	}

	return &CollectionParams{
		Cursor: cursor,
		Limit:  limit,
	}
}

// CollectionParams represents various parameters that can be applied to database methods
// that are working with collections (like pagination).
type CollectionParams struct {
	Cursor string // previous cursor
	Limit  uint8  // number of records to return
}

// NewCollectionInfo creates a new instance of CollectionInfo.
func NewCollectionInfo(cursor string, total uint64) (info *CollectionInfo) {
	return &CollectionInfo{
		Cursor: cursor,
		Total:  total,
	}
}

// CollectionInfo represents a statistical information returned from database methods that are working with collections.
type CollectionInfo struct {
	Cursor string // current cursor
	Total  uint64 // total number of records in the repo
}
