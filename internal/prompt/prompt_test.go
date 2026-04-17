package prompt

import (
	"bytes"
	"strings"
	"testing"
)

func newTestConfirmer(input string) (*Confirmer, *bytes.Buffer) {
	out := &bytes.Buffer{}
	return &Confirmer{In: strings.NewReader(input), Out: out}, out
}

func TestAsk_YesAnswer(t *testing.T) {
	c, _ := newTestConfirmer("y\n")
	ok, err := c.Ask("Continue?", false)
	if err != nil || !ok {
		t.Fatalf("expected true, nil; got %v, %v", ok, err)
	}
}

func TestAsk_NoAnswer(t *testing.T) {
	c, _ := newTestConfirmer("no\n")
	ok, err := c.Ask("Continue?", true)
	if err != nil || ok {
		t.Fatalf("expected false, nil; got %v, %v", ok, err)
	}
}

func TestAsk_DefaultYes_EmptyInput(t *testing.T) {
	c, _ := newTestConfirmer("\n")
	ok, err := c.Ask("Continue?", true)
	if err != nil || !ok {
		t.Fatalf("expected true, nil; got %v, %v", ok, err)
	}
}

func TestAsk_DefaultNo_EmptyInput(t *testing.T) {
	c, _ := newTestConfirmer("\n")
	ok, err := c.Ask("Continue?", false)
	if err != nil || ok {
		t.Fatalf("expected false, nil; got %v, %v", ok, err)
	}
}

func TestAsk_InvalidAnswer(t *testing.T) {
	c, _ := newTestConfirmer("maybe\n")
	_, err := c.Ask("Continue?", false)
	if err == nil {
		t.Fatal("expected error for unrecognised answer")
	}
}

func TestAsk_EOFTreatedAsDefault(t *testing.T) {
	c, _ := newTestConfirmer("")
	ok, err := c.Ask("Continue?", true)
	if err != nil || !ok {
		t.Fatalf("expected true, nil on EOF; got %v, %v", ok, err)
	}
}

func TestAsk_PrintsHint(t *testing.T) {
	c, out := newTestConfirmer("y\n")
	_, _ = c.Ask("Overwrite?", false)
	if !bytes.Contains(out.Bytes(), []byte("y/N")) {
		t.Errorf("expected hint y/N in output, got: %s", out.String())
	}
}
