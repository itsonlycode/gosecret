package secrets

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"sort"
	"strings"

	"github.com/itsonlycode/gosecret/pkg/debug"
	"github.com/itsonlycode/gosecret/pkg/gosecret"
)

var _ gosecret.Secret = &KV{}

// NewKV creates a new KV secret
func NewKV() *KV {
	return &KV{
		data: make(map[string][]string, 10),
	}
}

// NewKVWithData returns a new KV secret populated with data
func NewKVWithData(kvps map[string][]string, body string, converted bool) *KV {
	kv := &KV{
		data:     make(map[string][]string, len(kvps)),
		body:     body,
		fromMime: converted,
	}
	for k, v := range kvps {
		kv.data[k] = v
	}
	return kv
}

// KV is a secret that contains any number of
// lines of key-value pairs (defined as: contains a colon) and any number of
// free text lines. This is the default secret format gosecret uses and encourages.
//
// Format
// ------
// Line | Description
// ---- | -----------
//  0-n | Key-Value pairs, e.g. "key: value". Can be omitted but the secret
//      | might get parsed as a "Plain" secret if zero key-value pairs are found.
//  n+1 | Body. Can contain any number of characters that will be parsed as
//      | UTF-8 and appended to an internal string. Note: Technically this can
//      | be any kind of binary data but we neither support nor test this with
//      | non-text data. Also we do not intent do support any kind of streaming
//      | access, i.e. this is not intended for huge files.
//
// Example
// -------
// Line | Content
// ---- | -------
//    0 | hello: world
//    1 | gosecret: secret
//    2 | Yo
//    3 | Hi
//
// This would be parsed as a KV secret that contains:
//   - key-value pairs:
//     - "hello": "world"
//     - "gosecret": "secret"
//   - body: "Yo\nHi"
type KV struct {
	data     map[string][]string
	body     string
	fromMime bool
}

// Bytes serializes
func (k *KV) Bytes() []byte {
	buf := &bytes.Buffer{}
	for ik, key := range k.Keys() {
		sv, ok := k.data[key]
		if !ok {
			continue
		}
		for iv, v := range sv {
			_, _ = buf.WriteString(key)
			_, _ = buf.WriteString(": ")
			_, _ = buf.WriteString(v)
			// the last one shouldn't add a newline, it's handled below
			if iv < len(sv)-1 {
				_, _ = buf.WriteString("\n")
			}
		}
		// we must only add a final newline if the body is non-empty
		if k.body != "" || ik < len(k.Keys())-1 {
			_, _ = buf.WriteString("\n")
		}
	}
	buf.WriteString(k.body)
	return buf.Bytes()
}

// Keys returns all keys
func (k *KV) Keys() []string {
	keys := make([]string, 0, len(k.data)+1)
	for key := range k.data {
		keys = append(keys, key)
	}
	sort.Strings(keys)
	return keys
}

// Get returns the first value of that key
func (k *KV) Get(key string) (string, bool) {
	key = strings.ToLower(key)

	if v, found := k.data[key]; found {
		return v[0], true
	}

	return "", false
}

// Values returns all values for that key
func (k *KV) Values(key string) ([]string, bool) {
	key = strings.ToLower(key)
	v, found := k.data[key]
	return v, found
}

// Set writes a single key
func (k *KV) Set(key string, value interface{}) error {
	key = strings.ToLower(key)
	if v, ok := k.data[key]; ok && len(v) > 1 {
		return fmt.Errorf("cannot set key %s: this entry contains multiple same keys. Please use 'gosecret edit' instead", key)
	}
	k.data[key] = []string{fmt.Sprintf("%s", value)}
	return nil
}

// Add appends data to a given key
func (k *KV) Add(key string, value interface{}) error {
	key = strings.ToLower(key)
	k.data[key] = append(k.data[key], fmt.Sprintf("%s", value))
	return nil
}

// Del removes a given key and all of its values
func (k *KV) Del(key string) bool {
	key = strings.ToLower(key)
	_, found := k.data[key]
	delete(k.data, key)
	return found
}

// Body returns the body
func (k *KV) Body() string {
	return k.body
}

// ParseKV tries to parse a KV secret
func ParseKV(in []byte) (*KV, error) {
	k := &KV{
		data: make(map[string][]string, 10),
	}
	r := bufio.NewReader(bytes.NewReader(in))

	var sb strings.Builder
	for {
		line, err := r.ReadString('\n')
		if err != nil && line == "" {
			if err == io.EOF {
				break
			}
			return nil, err
		}
		// append non KV pairs to the body
		if !strings.Contains(line, ":") {
			sb.WriteString(line)
			continue
		}
		line = strings.TrimRight(line, "\n")

		parts := strings.SplitN(line, ":", 2)
		// should not happen
		if len(parts) < 1 {
			continue
		}
		for i, part := range parts {
			parts[i] = strings.TrimSpace(part)
		}
		// we only store lower case keys for KV
		parts[0] = strings.ToLower(parts[0])
		// preserve key only entries
		if len(parts) < 2 {
			k.data[parts[0]] = append(k.data[parts[0]], "")
			continue
		}

		k.data[parts[0]] = append(k.data[parts[0]], parts[1])
	}
	if len(k.data) < 1 {
		debug.Log("no KV entries")
	}
	k.body = sb.String()
	return k, nil
}

// Write appends the buffer to the secret's body
func (k *KV) Write(buf []byte) (int, error) {
	k.body += string(buf)
	return len(buf), nil
}

// FromMime returns whether this secret was converted from a Mime secret of not
func (k *KV) FromMime() bool {
	return k.fromMime
}

// SafeStr always returnes "(elided)"
func (k *KV) SafeStr() string {
	return "(elided)"
}
