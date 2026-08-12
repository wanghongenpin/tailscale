// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

package osrouter

import (
	"net/netip"
	"path/filepath"
	"testing"
)

func TestGetNetshPath(t *testing.T) {
	ft := &firewallTweaker{
		logf: t.Logf,
	}
	path := ft.getNetshPath()
	if !filepath.IsAbs(path) {
		t.Errorf("expected absolute path for netsh.exe: %q", path)
	}
}

func TestHasDefaultRoute(t *testing.T) {
	for name, routes := range map[string][]netip.Prefix{
		"no routes":      nil,
		"specific route": {netip.MustParsePrefix("192.168.1.0/24")},
	} {
		if hasDefaultRoute(routes) {
			t.Errorf("%s: hasDefaultRoute() = true, want false", name)
		}
	}
	for _, route := range []string{"0.0.0.0/0", "::/0"} {
		if !hasDefaultRoute([]netip.Prefix{netip.MustParsePrefix(route)}) {
			t.Errorf("%s: hasDefaultRoute() = false, want true", route)
		}
	}
}
