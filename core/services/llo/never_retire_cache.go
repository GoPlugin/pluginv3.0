package llo

import (
	ocrtypes "github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	llotypes "github.com/goplugin/plugin-common/pkg/types/llo"
)

type neverShouldRetireCache struct{}

func NewNeverShouldRetireCache() llotypes.ShouldRetireCache {
	return &neverShouldRetireCache{}
}

func (n *neverShouldRetireCache) ShouldRetire(digest ocrtypes.ConfigDigest) (bool, error) {
	return false, nil
}
