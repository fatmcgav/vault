package vault

import (
	"testing"

	"github.com/hashicorp/vault/physical"
)

// This file contains a number of methods that are useful for unit
// tests within other packages.

// TestCore returns a pure in-memory, uninitialized core for testing.
func TestCore(t *testing.T) *Core {
	physicalBackend := physical.NewInmem()
	c, err := NewCore(&CoreConfig{
		Physical: physicalBackend,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}

	return c
}

// TestCoreInit initializes the core with a single key, and returns
// the key that must be used to unseal the core and a root token.
func TestCoreInit(t *testing.T, core *Core) ([]byte, string) {
	result, err := core.Initialize(&SealConfig{
		SecretShares:    1,
		SecretThreshold: 1,
	})
	if err != nil {
		t.Fatalf("err: %s", err)
	}
	return result.SecretShares[0], result.RootToken
}

// TestCoreUnsealed returns a pure in-memory core that is already
// initialized and unsealed.
func TestCoreUnsealed(t *testing.T) (*Core, []byte) {
	core, key, _ := TestCoreUnsealedToken(t)
	return core, key
}

// TestCoreUnsealedToken returns a pure in-memory core that is already
// initialized and unsealed along with the root token.
func TestCoreUnsealedToken(t *testing.T) (*Core, []byte, string) {
	core := TestCore(t)
	key, token := TestCoreInit(t, core)
	if _, err := core.Unseal(TestKeyCopy(key)); err != nil {
		t.Fatalf("unseal err: %s", err)
	}

	sealed, err := core.Sealed()
	if err != nil {
		t.Fatalf("err checking seal status: %s", err)
	}
	if sealed {
		t.Fatal("should not be sealed")
	}

	return core, key, token
}

// TestKeyCopy is a silly little function to just copy the key so that
// it can be used with Unseal easily.
func TestKeyCopy(key []byte) []byte {
	result := make([]byte, len(key))
	copy(result, key)
	return result
}
