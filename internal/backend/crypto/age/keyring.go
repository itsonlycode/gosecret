package age

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"path/filepath"

	"filippo.io/age"
	"github.com/itsonlycode/gosecret/pkg/ctxutil"
	"github.com/itsonlycode/gosecret/pkg/debug"
	"github.com/itsonlycode/gosecret/pkg/termio"
)

// Keyring is an age keyring
type Keyring []Keypair

// Keypair is a public / private keypair
type Keypair struct {
	Name     string `json:"name"`
	Email    string `json:"email"`
	Identity string `json:"identity"`
}

func (a *Age) pkself(ctx context.Context) (age.Recipient, error) {
	kr, err := a.loadKeyring(ctx)

	var id *age.X25519Identity
	if err != nil || len(kr) < 1 {
		id, err = a.genKey(ctx)
	} else {
		id, err = age.ParseX25519Identity(kr[len(kr)-1].Identity)
	}
	if err != nil {
		return nil, err
	}
	return id.Recipient(), nil
}

func (a *Age) genKey(ctx context.Context) (*age.X25519Identity, error) {
	debug.Log("No native age key found. Generating ...")
	id, err := a.generateIdentity(ctx, termio.DetectName(ctx, nil), termio.DetectEmail(ctx, nil))
	if err != nil {
		return nil, err
	}
	return id, nil
}

// GenerateIdentity will create a new native private key
func (a *Age) GenerateIdentity(ctx context.Context, name, email, _ string) error {
	_, err := a.generateIdentity(ctx, name, email)
	return err
}

func (a *Age) generateIdentity(ctx context.Context, name, email string) (*age.X25519Identity, error) {
	id, err := age.GenerateX25519Identity()
	if err != nil {
		return id, err
	}

	var newKeyring bool
	kr, err := a.loadKeyring(ctx)
	if err != nil {
		debug.Log("failed to load existing keyring from %s: %s", a.keyring, err)
		if !errors.Is(err, os.ErrNotExist) {
			return nil, err
		}
		debug.Log("no existing keyring, creating new one")
		newKeyring = true
	}

	kr = append(kr, Keypair{
		Name:     name,
		Email:    email,
		Identity: id.String(),
	})

	return id, a.saveKeyring(ctx, kr, newKeyring)
}

func (a *Age) loadKeyring(ctx context.Context) (Keyring, error) {
	if !ctxutil.HasPasswordCallback(ctx) {
		debug.Log("no password callback found, redirecting to askPass")
		ctx = ctxutil.WithPasswordCallback(ctx, func(prompt string, _ bool) ([]byte, error) {
			pw, err := a.askPass.Passphrase(prompt, fmt.Sprintf("to load the age keyring at %s", a.keyring), false)
			return []byte(pw), err
		})
	}

	buf, err := a.decryptFile(ctx, a.keyring)
	if err != nil {
		debug.Log("can't decrypt keyring at %s: %s", a.keyring, err)
		return Keyring{}, err
	}

	var kr Keyring
	if err := json.Unmarshal(buf, &kr); err != nil {
		debug.Log("can't parse keyring at %s: %s", a.keyring, err)
		return Keyring{}, err
	}

	// remove invalid IDs
	valid := make(Keyring, 0, len(kr))
	for _, k := range kr {
		if k.Identity == "" {
			continue
		}
		valid = append(valid, k)
	}
	debug.Log("loaded keyring with %d valid entries from %s", len(kr), a.keyring)
	return valid, nil
}

func (a *Age) saveKeyring(ctx context.Context, k Keyring, newKeyring bool) error {
	if !ctxutil.HasPasswordCallback(ctx) {
		debug.Log("no password callback found, redirecting to askPass")
		ctx = ctxutil.WithPasswordCallback(ctx, func(prompt string, confirm bool) ([]byte, error) {
			pw, err := a.askPass.Passphrase(prompt, fmt.Sprintf("to save the age keyring to %s", a.keyring), confirm)
			return []byte(pw), err
		})
	}

	if err := os.MkdirAll(filepath.Dir(a.keyring), 0700); err != nil {
		debug.Log("failed to create directory for the keyring at %s", a.keyring)
		return err
	}

	// encrypt final keyring
	buf, err := json.Marshal(k)
	if err != nil {
		return err
	}

	if err := a.encryptFile(ctx, a.keyring, buf, newKeyring); err != nil {
		return err
	}

	debug.Log("saved encrypted keyring with %d entries to %s", len(k), a.keyring)
	return nil
}
