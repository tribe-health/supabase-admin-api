package optimizations

import "testing"

func TestOptimizePostgrest(t *testing.T) {
	settings := PostgrestServerSettings{
		DbPool: 50,
	}
	result, _ := generateSettings(settings)
	if *result != `db-pool = 50
` {
		t.Fatal(*result)
	}
}
