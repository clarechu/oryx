// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT

package main

import (
	"context"
	"platform/datasource"
)

var fastCache *FastCache

type FastCache struct {
	// Whether delivery HLS in high performance mode.
	HLSHighPerformance bool
	// Whether deliver HLS in low latency mode.
	HLSLowLatency bool
	ds            datasource.Datasource
}

func NewFastCache(ds datasource.Datasource) *FastCache {
	return &FastCache{ds: ds}
}

func (v *FastCache) Refresh(ctx context.Context) error {
	if vs, _ := v.ds.Get(ctx, SRS_LL_HLS, "hlsLowLatency"); vs == "true" {
		v.HLSLowLatency = true
	} else {
		v.HLSLowLatency = false
	}

	if vs, _ := v.ds.Get(ctx, SRS_HP_HLS, "noHlsCtx"); vs == "true" {
		v.HLSHighPerformance = true
	} else {
		v.HLSHighPerformance = false
	}

	return nil
}
