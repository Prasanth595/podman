package pods

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"

	"go.podman.io/podman/v6/pkg/bindings"
	"go.podman.io/podman/v6/pkg/domain/entities"
)

// Clone creates a copy of the given pod, applying the overrides in options.
func Clone(ctx context.Context, options entities.PodCloneOptions) (*entities.PodCloneReport, error) {
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	response, err := conn.DoRequest(ctx, bytes.NewReader(body), http.MethodPost, "/pods/%s/clone", nil, nil, options.ID)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var report entities.PodCloneReport
	if err := response.Process(&report); err != nil {
		return nil, err
	}
	return &report, nil
}
