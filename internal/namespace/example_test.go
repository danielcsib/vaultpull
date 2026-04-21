package namespace_test

import (
	"fmt"

	"github.com/your-org/vaultpull/internal/namespace"
)

func ExampleScoper_Apply() {
	s, _ := namespace.New("prod", "/")
	secrets := map[string]string{
		"DB_HOST": "db.example.com",
		"DB_PORT": "5432",
	}
	scoped := s.Apply(secrets)
	fmt.Println(scoped["prod/DB_HOST"])
	fmt.Println(scoped["prod/DB_PORT"])
	// Output:
	// db.example.com
	// 5432
}
