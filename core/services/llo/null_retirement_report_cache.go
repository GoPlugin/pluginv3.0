package llo

import (
	"context"

	"github.com/goplugin/plugin-libocr/offchainreporting2plus/types"
	ocr2types "github.com/goplugin/plugin-libocr/offchainreporting2plus/types"

	datastreamsllo "github.com/goplugin/plugin-data-streams/llo"
)

type NullRetirementReportCache struct{}

func (n *NullRetirementReportCache) StoreAttestedRetirementReport(ctx context.Context, cd ocr2types.ConfigDigest, retirementReport []byte, sigs []types.AttributedOnchainSignature) error {
	return nil
}
func (n *NullRetirementReportCache) StoreConfig(ctx context.Context, cd ocr2types.ConfigDigest, signers [][]byte, f uint8) error {
	return nil
}
func (n *NullRetirementReportCache) AttestedRetirementReport(predecessorConfigDigest ocr2types.ConfigDigest) ([]byte, error) {
	return nil, nil
}
func (n *NullRetirementReportCache) CheckAttestedRetirementReport(predecessorConfigDigest ocr2types.ConfigDigest, attestedRetirementReport []byte) (datastreamsllo.RetirementReport, error) {
	return datastreamsllo.RetirementReport{}, nil
}
