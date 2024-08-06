// Copyright (c) 2022-2024 Winlin
//
// SPDX-License-Identifier: MIT

//go:build linux

package main

import (
	"context"
)

// queryLatestVersion is to query the latest and stable version from Oryx API.
func queryLatestVersion(ctx context.Context) (*Versions, error) {
	return &Versions{
		Version: version,
		Stable:  "v1.0.193",
		Latest:  "v1.0.307",
	}, nil
}
