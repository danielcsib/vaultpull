package lineage_test

import (
	"os"

	"github.com/your-org/vaultpull/internal/lineage"
)

func ExampleStore_Record() {
	s, err := lineage.NewStore("")
	if err != nil {
		panic(err)
	}

	s.Record("DB_PASSWORD", "secret/myapp/db", "password")
	s.Record("API_KEY", "secret/myapp/api", "key")

	lineage.PrintSummary(os.Stdout, s)
	// Output is a formatted table; exact timestamps vary, so we omit Output tag.
}
