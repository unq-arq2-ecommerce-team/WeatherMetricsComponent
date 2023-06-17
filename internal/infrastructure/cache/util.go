package cache

import "fmt"

const tableAndKeySeparator = ":"

func BuildCacheKey(table, key string) string {
	return fmt.Sprintf("%s%s%s", table, tableAndKeySeparator, key)
}
