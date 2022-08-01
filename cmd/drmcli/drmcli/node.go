package drmcli

import (
	"encoding/json"
	"os"

	"git.sr.ht/~emersion/go-drm"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

type node struct {
	drm.Node
	f *os.File
}

func openNode(device string) (*node, error) {
	f, err := os.Open(device)
	if err != nil {
		return nil, errors.Wrap(err, "cannot open DRM module file")
	}

	return &node{
		Node: *drm.NewNode(f.Fd()),
		f:    f,
	}, nil
}

func nodeFuncToJSON(ctx *cli.Context, f func(*node) (any, error)) error {
	n, err := openNode(ctx.String("device"))
	if err != nil {
		return err
	}

	v, err := f(n)
	if err != nil {
		return errors.Wrap(err, "cannot get")
	}

	if err := json.NewEncoder(os.Stdout).Encode(v); err != nil {
		return errors.Wrap(err, "cannot encode JSON")
	}

	return nil
}
