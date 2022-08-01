package drmcli

import (
	"encoding/hex"

	"git.sr.ht/~emersion/go-drm"
	"github.com/pkg/errors"
	"github.com/urfave/cli/v2"
)

func cmdResources(ctx *cli.Context) error {
	return nodeFuncToJSON(ctx, func(n *node) (any, error) {
		return n.ModeGetResources()
	})
}

type jsonModeConnector struct {
	PossibleEncoders []drm.EncoderID
	Modes            []drm.ModeModeInfo

	Encoder drm.EncoderID
	ID      drm.ConnectorID
	Type    string

	Status    string
	PhyWidth  uint32
	PhyHeight uint32 // mm
	Subpixel  string
}

func cmdConnectorInfo(ctx *cli.Context) error {
	id, err := firstArgID(ctx)
	if err != nil {
		return err
	}

	return nodeFuncToJSON(ctx, func(n *node) (any, error) {
		c, err := n.ModeGetConnector(drm.ConnectorID(id))
		if err != nil || ctx.Bool("raw") {
			return c, err
		}

		return jsonModeConnector{
			PossibleEncoders: c.PossibleEncoders,
			Modes:            c.Modes,
			Encoder:          c.Encoder,
			ID:               c.ID,
			Type:             c.Type.String(),
			Status:           c.Status.String(),
			PhyWidth:         c.PhyWidth,
			PhyHeight:        c.PhyHeight,
			Subpixel:         c.Subpixel.String(),
		}, nil
	})
}

func cmdConnectorProperties(ctx *cli.Context) error {
	return cmdProperties(ctx, drm.ObjectConnector)
}

func cmdProperties(ctx *cli.Context, typ drm.ObjectType) error {
	id, err := firstArgID(ctx)
	if err != nil {
		return err
	}

	return nodeFuncToJSON(ctx, func(n *node) (any, error) {
		return n.ModeObjectGetProperties(drm.NewAnyID(id, typ))
	})
}

type jsonProperty struct {
	ID        drm.PropertyID
	Name      string
	Type      string
	Atomic    bool
	Immutable bool

	Blobs       []string          `json:",omitempty"`
	Enums       map[uint64]string `json:",omitempty"`
	Range       []uint64          `json:",omitempty"`
	SignedRange []int64           `json:",omitempty"`
}

func cmdProperty(ctx *cli.Context) error {
	id, err := firstArgID(ctx)
	if err != nil {
		return err
	}

	return nodeFuncToJSON(ctx, func(n *node) (any, error) {
		p, err := n.ModeGetProperty(drm.PropertyID(id))
		if err != nil {
			return nil, err
		}

		if ctx.Bool("raw") {
			return p, nil
		}

		prop := jsonProperty{
			ID:        drm.PropertyID(id),
			Name:      p.Name,
			Type:      p.Type().String(),
			Atomic:    p.Atomic(),
			Immutable: p.Immutable(),
		}

		if blobs, ok := p.Blobs(); ok {
			for _, blob := range blobs {
				b, err := n.ModeGetBlob(blob.ID)
				if err != nil {
					return nil, errors.Wrapf(err, "cannot get blob %v for property %v", blob.ID, id)
				}
				prop.Blobs = append(prop.Blobs, hex.EncodeToString(b))
			}
		}

		if enums, ok := p.Enums(); ok {
			prop.Enums = make(map[uint64]string, len(enums))
			for _, enum := range enums {
				prop.Enums[enum.Value] = enum.Name
			}
		}

		if lo, hi, ok := p.Range(); ok {
			prop.Range = []uint64{lo, hi}
		}

		if lo, hi, ok := p.SignedRange(); ok {
			prop.SignedRange = []int64{lo, hi}
		}

		return prop, nil
	})
}
