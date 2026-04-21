package quota_test

import (
	"fmt"
	"os"

	"github.com/your-org/vaultpull/internal/quota"
)

func ExampleStore_Check() {
	s := quota.New(2)

	for i := 0; i < 3; i++ {
		err := s.Check("secret/data/myapp")
		if err != nil {
			fmt.Println("blocked:", err)
		} else {
			fmt.Printf("read %d ok, remaining=%d\n", i+1, s.Remaining("secret/data/myapp"))
		}
	}

	quota.PrintSummary(s, os.Stdout)

	// Output:
	// read 1 ok, remaining=1
	// read 2 ok, remaining=0
	// blocked: quota exceeded: path "secret/data/myapp" has been read 3 times (max 2)
	// quota usage:
	//   secret/data/myapp                        reads=3    remaining=0    [EXCEEDED]
}
