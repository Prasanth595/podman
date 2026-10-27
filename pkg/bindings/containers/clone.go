package containers

import (
	"bytes"
	"context"
	"net/http"

	"go.podman.io/podman/v6/pkg/bindings"
	"go.podman.io/podman/v6/pkg/domain/entities"
)

// Clone creates a copy of the given container, applying the overrides in options.
func Clone(ctx context.Context, options entities.ContainerCloneOptions) (*entities.ContainerCreateReport, error) {
	conn, err := bindings.GetClient(ctx)
	if err != nil {
		return nil, err
	}
	body, err := json.Marshal(options)
	if err != nil {
		return nil, err
	}
	response, err := conn.DoRequest(ctx, bytes.NewReader(body), http.MethodPost, "/containers/%s/clone", nil, nil, options.ID)
	if err != nil {
		return nil, err
	}
	defer response.Body.Close()

	var report entities.ContainerCreateReport
	if err := response.Process(&report); err != nil {
		return nil, err
	}
	return &report, nil
}
