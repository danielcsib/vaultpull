package dedupe_test

import (
	"fmt"

	"github.com/your-org/vaultpull/internal/dedupe"
)

func ExampleMerger_Merge() {
	m := dedupe.New(dedupe.PolicyKeepLast)

	base := map[string]string{
		"DB_HOST": "localhost",
		"DB_PORT": "5432",
	}
	override := map[string]string{
		"DB_HOST": "prod.internal",
		"API_KEY": "topsecret",
	}

	if err := m.Merge(base, override); err != nil {
		fmt.Println("error:", err)
		return
	}

	fmt.Println("DB_HOST:", base["DB_HOST"])
	fmt.Println("conflicts:", len(m.Conflicts))
	// Output:
	// DB_HOST: prod.internal
	// conflicts: 1
}
