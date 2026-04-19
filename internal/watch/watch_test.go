package watch_test

import (
	"context"
	"sync"
	"testing"
	"time"

	"github.com/your-org/vaultpull/internal/watch"
)

type fakeClient struct {
	mu      sync.Mutex
	results map[string]string
	err     error
	calls   int
}

func (f *fakeClient) ReadSecrets(_ string) (map[string]string, error) {
	f.mu.Lock()
	defer f.mu.Unlock()
	f.calls++
	return f.results, f.err
}

func TestNew_DefaultInterval(t *testing.T) {
	w := watch.New(&fakeClient{}, []string{"secret/app"}, 0, nil)
	if w == nil {
		t.Fatal("expected non-nil watcher")
	}
}

func TestRun_CallsOnChange_WhenSecretsChange(t *testing.T) {
	client := &fakeClient{results: map[string]string{"KEY": "val1"}}
	var mu sync.Mutex
	var fired []string

	w := watch.New(client, []string{"secret/app"}, 20*time.Millisecond, func(path string, _ map[string]string) {
		mu.Lock()
		fired = append(fired, path)
		mu.Unlock()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 80*time.Millisecond)
	defer cancel()

	go w.Run(ctx)
	time.Sleep(30 * time.Millisecond)

	client.mu.Lock()
	client.results = map[string]string{"KEY": "val2"}
	client.mu.Unlock()

	<-ctx.Done()

	mu.Lock()
	defer mu.Unlock()
	if len(fired) < 2 {
		t.Fatalf("expected at least 2 onChange calls, got %d", len(fired))
	}
}

func TestRun_NoChange_NoCallback(t *testing.T) {
	client := &fakeClient{results: map[string]string{"KEY": "same"}}
	callCount := 0
	var mu sync.Mutex

	w := watch.New(client, []string{"secret/app"}, 20*time.Millisecond, func(_ string, _ map[string]string) {
		mu.Lock()
		callCount++
		mu.Unlock()
	})

	ctx, cancel := context.WithTimeout(context.Background(), 70*time.Millisecond)
	defer cancel()
	go w.Run(ctx)
	<-ctx.Done()

	mu.Lock()
	defer mu.Unlock()
	// Only the first poll triggers onChange; subsequent identical polls do not.
	if callCount != 1 {
		t.Fatalf("expected 1 onChange call, got %d", callCount)
	}
}
