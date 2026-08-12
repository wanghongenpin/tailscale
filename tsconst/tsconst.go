// Copyright (c) Tailscale Inc & contributors
// SPDX-License-Identifier: BSD-3-Clause

// Package tsconst exports some constants used elsewhere in the
// codebase.
package tsconst

// WintunInterfaceDesc is the description attached to this client's
// interfaces on Windows. This is set by the WinTun driver, which renders
// it as "<tunnelType> Tunnel" from the tunnelType passed to
// wintun.CreateAdapter.
//
// meshpin fork: renamed from "Tailscale Tunnel" to "MeshPin Tunnel" so
// the adapter carries no upstream branding — this string is *user
// visible*, in Get-NetAdapter, Control Panel, and the registry's
// DeviceDesc.
//
// Renaming is safe precisely because both ends of the match live in code
// we control: tsbridge's wintunTunnelType writes the description, and
// the three readers below compare against this constant.
//
//	net/netmon/interfaces_windows.go  notTailscaleInterface()
//	net/netns/netns_windows.go        isTailscaleInterface()
//
// The value MUST stay in sync with tsbridge's wintunTunnelType +
// " Tunnel" (native/tsbridge/tun_windows.go). If they drift, this
// client's own TUN is classified as a normal adapter, netns binds every
// outbound socket to the tunnel itself, and selecting an exit node kills
// all connectivity — see the comment on wintunTunnelType for the full
// mechanism.
const WintunInterfaceDesc = "MeshPin Tunnel"

// WintunInterfaceDesc0_14 is the legacy description used by Wintun
// 0.14 and earlier. Left as upstream: it is only ever *read*, to
// recognise adapters left behind by an old driver, so rebranding it
// would just stop matching those.
const WintunInterfaceDesc0_14 = "Wintun Userspace Tunnel"

// TailnetLockNotTrustedMsg is the error message used by tailnet lock
// and sniffed (via substring) out of an error sent over the network.
const TailnetLockNotTrustedMsg = "this node is not trusted by Tailnet Lock"
