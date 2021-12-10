package optimizations

import "testing"

func TestOptimizePostgrest(t *testing.T) {
	settings := PostgrestServerSettings{
		DbPool: 50,
	}
	result, _ := generateSettings(settings)
	if `db-pool = 50
` != *result {
		t.Fatal(*result)
	}
}
