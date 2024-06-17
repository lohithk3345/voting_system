package reporesult

type InsertResult interface{}

type MultiResult interface{}

type StoreResponse struct {
	Code    int64
	Message string
}

func IsModified(v int64) bool {
	return v == 1
}

func IsUpserted(v int64) bool {
	return v == 1
}

func IsMatched(v int64) bool {
	return v == 1
}

func IsDeletedCount(v int64) bool {
	return v > 0
}
