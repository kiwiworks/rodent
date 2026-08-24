package command

import (
	"testing"

	"github.com/spf13/cobra"

	"github.com/kiwiworks/rodent/system/manifest"
)

func testManifest() *manifest.Manifest {
	return manifest.New("testapp", "0.1.0")
}

func findChild(cmd *cobra.Command, name string) *cobra.Command {
	for _, c := range cmd.Commands() {
		if c.Name() == name {
			return c
		}
	}
	return nil
}

func TestNewRoot_SameLeafAcrossGroupsDoesNotCollide(t *testing.T) {
	commands := []*Command{
		New("alpha", "alpha group", ""),
		New("beta", "beta group", ""),
		New("alpha.run", "alpha run", ""),
		New("beta.run", "beta run", ""),
	}
	root, err := NewRoot(RootParams{Manifest: testManifest(), Commands: commands})
	if err != nil {
		t.Fatalf("NewRoot returned error: %v", err)
	}

	alpha := findChild(root.root, "alpha")
	beta := findChild(root.root, "beta")
	if alpha == nil || beta == nil {
		t.Fatalf("expected both group commands to be registered as root children")
	}

	alphaRun := findChild(alpha, "run")
	betaRun := findChild(beta, "run")
	if alphaRun == nil || betaRun == nil {
		t.Fatalf("expected both run leaves to be registered under their own group")
	}
	if alphaRun == betaRun {
		t.Fatalf("alpha run and beta run resolved to the same *cobra.Command")
	}
	if alphaRun.Short != "alpha run" {
		t.Errorf("alpha run has wrong Short: got %q", alphaRun.Short)
	}
	if betaRun.Short != "beta run" {
		t.Errorf("beta run has wrong Short: got %q", betaRun.Short)
	}
}

func TestNewRoot_SameLeafAtIntermediateLevelDoesNotCollide(t *testing.T) {
	commands := []*Command{
		New("a", "a group", ""),
		New("x", "x group", ""),
		New("a.b", "a.b group", ""),
		New("x.b", "x.b group", ""),
		New("a.b.c", "a.b.c leaf", ""),
		New("x.b.c", "x.b.c leaf", ""),
	}
	root, err := NewRoot(RootParams{Manifest: testManifest(), Commands: commands})
	if err != nil {
		t.Fatalf("NewRoot returned error: %v", err)
	}

	a := findChild(root.root, "a")
	x := findChild(root.root, "x")
	if a == nil || x == nil {
		t.Fatalf("expected both top-level groups to be registered")
	}

	ab := findChild(a, "b")
	xb := findChild(x, "b")
	if ab == nil || xb == nil {
		t.Fatalf("expected intermediate 'b' groups to be registered under their own parent")
	}
	if ab == xb {
		t.Fatalf("a.b and x.b resolved to the same *cobra.Command")
	}

	abc := findChild(ab, "c")
	xbc := findChild(xb, "c")
	if abc == nil || xbc == nil {
		t.Fatalf("expected leaf 'c' to be registered under its own intermediate group")
	}
	if abc == xbc {
		t.Fatalf("a.b.c and x.b.c resolved to the same *cobra.Command")
	}
	if abc.Short != "a.b.c leaf" || xbc.Short != "x.b.c leaf" {
		t.Errorf("leaves resolved to wrong implementation: abc=%q xbc=%q", abc.Short, xbc.Short)
	}
}

func TestNewRoot_DuplicateRegistrationErrors(t *testing.T) {
	commands := []*Command{
		New("alpha.run", "first", ""),
		New("alpha.run", "second", ""),
	}
	_, err := NewRoot(RootParams{Manifest: testManifest(), Commands: commands})
	if err == nil {
		t.Fatalf("expected an error for duplicate command registration, got nil")
	}
}
