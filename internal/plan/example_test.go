package plan_test

import (
	"os"

	"github.com/your-org/vaultpull/internal/plan"
)

func ExampleBuild() {
	current := map[string]string{
		"API_KEY": "old-key",
		"DB_HOST": "localhost",
	}
	incoming := map[string]string{
		"API_KEY": "new-key",
		"DB_PORT": "5432",
	}

	p := plan.Build(".env", current, incoming)
	p.Print(os.Stdout)

	// Output:
	// Plan for .env:
	//   ~ API_KEY
	//   + DB_PORT
	//   - DB_HOST
}
