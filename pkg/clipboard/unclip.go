package clipboard

import (
	"context"
	"crypto/sha256"
	"fmt"

	"github.com/itsonlycode/gosecret/internal/notify"

	"github.com/atotto/clipboard"
)

// Clear will attempt to erase the clipboard
func Clear(ctx context.Context, checksum string, force bool) error {
	if clipboard.Unsupported {
		return ErrNotSupported
	}

	cur, err := clipboard.ReadAll()
	if err != nil {
		return fmt.Errorf("failed to read clipboard: %w", err)
	}

	hash := fmt.Sprintf("%x", sha256.Sum256([]byte(cur)))
	if hash != checksum && !force {
		return nil
	}

	if err := clipboard.WriteAll(""); err != nil {
		_ = notify.Notify(ctx, "gosecret - clipboard", "Failed to clear clipboard")
		return fmt.Errorf("failed to write clipboard: %w", err)
	}

	if err := clearClipboardHistory(ctx); err != nil {
		_ = notify.Notify(ctx, "gosecret - clipboard", "Failed to clear clipboard history")
		return fmt.Errorf("failed to clear clipboard history: %w", err)
	}

	if err := notify.Notify(ctx, "gosecret - clipboard", "Clipboard has been cleared"); err != nil {
		return fmt.Errorf("failed to send unclip notification: %w", err)
	}

	return nil
}
