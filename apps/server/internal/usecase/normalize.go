package usecase

func normalizeLimit(limit int) int {
	if limit <= 0 {
		return DefaultLimit
	}
	if limit > MaxLimit {
		return MaxLimit
	}
	return limit
}

func normalizeOffset(offset int) int {
	if offset < 0 {
		return 0
	}
	return offset
}

func normalizeOrder(order string) string {
	if order == "desc" || order == "DESC" {
		return "DESC"
	}
	return "ASC"
}
