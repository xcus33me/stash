package fileapp

import (
	"context"
	"fmt"
)

type CleanupExpiredOutput struct {
	Processed int
	Failed    int
}

func (uc *Usecase) CleanupExpired(ctx context.Context) (CleanupExpiredOutput, error) {
	logger := uc.log.With("op", "cleanup_expired")

	items, err := uc.files.ListExpired(ctx)
	if err != nil {
		return CleanupExpiredOutput{}, fmt.Errorf("list expired: %w", err)
	}

	var out CleanupExpiredOutput
	for _, f := range items {
		if err := uc.storage.Delete(ctx, f.StorageKey); err != nil {
			logger.ErrorContext(ctx, "storage delete failed for expired file",
				"err", err,
				"file_id", f.ID,
				"storage_key", f.StorageKey,
			)
		}
		if err := uc.files.SoftDelete(ctx, f.ID); err != nil {
			logger.ErrorContext(ctx, "soft delete failed for expired file",
				"err", err,
				"file_id", f.ID,
			)
			out.Failed++
			continue
		}
		out.Processed++
	}

	logger.InfoContext(ctx, "cleanup finished",
		"processed", out.Processed,
		"failed", out.Failed,
	)

	return out, nil
}
