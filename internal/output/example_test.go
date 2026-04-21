package output_test

import (
	"os"

	"github.com/your-org/vaultpull/internal/output"
)

// ExampleWriter_Write demonstrates writing secrets in the export format.
func ExampleWriter_Write() {
	secrets := map[string]string{
		"APP_ENV": "production",
		"LOG_LEVEL": "info",
	}

	w := output.New(output.FormatExport, os.Stdout)
	_ = w.Write(secrets)

	// Output:
	// export APP_ENV=production
	// export LOG_LEVEL=info
}
