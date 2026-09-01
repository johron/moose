package highlight

import (
	"context"

	"github.com/tetratelabs/wazero"
	//"github.com/tetratelabs/wazero/imports/wasi_snapshot_preview1"
)

type HighlightEngine struct {
}

func NewHighlightEngine() *HighlightEngine {
	ctx := context.Background()
	r := wazero.NewRuntime(ctx)
	defer r.Close(ctx)

	return &HighlightEngine {
		
	}
}