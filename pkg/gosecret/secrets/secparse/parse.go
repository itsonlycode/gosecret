package secparse

import (
	"github.com/itsonlycode/gosecret/internal/out"
	"github.com/itsonlycode/gosecret/pkg/debug"
	"github.com/itsonlycode/gosecret/pkg/gosecret"
	"github.com/itsonlycode/gosecret/pkg/gosecret/secrets"
)

// Parse tries to parse a secret. It will start with the most specific
// secrets type.
func Parse(in []byte) (gosecret.Secret, error) {
	var s gosecret.Secret
	var err error
	s, err = parseLegacyMIME(in)
	if err == nil {
		debug.Log("parsed as MIME: %+v", s)
		return s, nil
	}
	debug.Log("failed to parse as MIME: %s", out.Secret(err.Error()))

	if _, ok := err.(*secrets.PermanentError); ok {
		return secrets.ParsePlain(in), err
	}
	s, err = secrets.ParseYAML(in)
	if err == nil {
		debug.Log("parsed as YAML: %+v", s)
		return s, nil
	}
	debug.Log("failed to parse as YAML: %s\n%s", err, out.Secret(string(in)))

	s, err = secrets.ParseKV(in)
	if err == nil {
		debug.Log("parsed as KV: %+v", s)
		return s, nil
	}
	debug.Log("failed to parse as KV: %s", err)

	s = secrets.ParsePlain(in)
	debug.Log("parsed as plain: %s", out.Secret(s.Bytes()))
	return s, nil
}
