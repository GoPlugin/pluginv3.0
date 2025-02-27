package solana

import "github.com/goplugin/plugin-solana/pkg/solana/internal"

func LamportsToSol(lamports uint64) float64 { return internal.LamportsToSol(lamports) }
