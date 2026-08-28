// Copyright (c) Tailscale Inc & AUTHORS
// SPDX-License-Identifier: BSD-3-Clause

package netstack

import "sync/atomic"

// This file is a meshpin-local addition to the vendored fork. Upstream has no
// way to size an individual netstack endpoint's send buffer: gonet.TCPConn
// keeps its tcpip.Endpoint unexported and offers no SO_SNDBUF equivalent, so an
// application holding only the net.Conn cannot reach it. The accept path here
// still has the endpoint in hand, which makes this the only place the option
// can be applied.
//
// Why an application would want that: netstack creates every TCP endpoint with
// the stack-wide default send buffer (tcpTXBufDefSize, 1 MiB). For bulk
// transfers that is exactly right — it is roughly the bandwidth-delay product
// of a fast path and keeps the pipe full. For a real-time stream it is a
// latency reservoir: the sender's Write returns as soon as the bytes are
// copied into that buffer, so a sender that is producing faster than the link
// drains sees no backpressure at all while up to a megabyte of already-stale
// frames sits queued ahead of the ones being produced. At 2 Mbps a full 1 MiB
// buffer is about four seconds of video the viewer must display before it sees
// anything current, and nothing in the sending application can observe it.
//
// Sizing such a connection to a couple of times its bandwidth-delay product
// instead restores the property a real socket has: Write blocks when the link
// is behind, which is the signal the application needs in order to slow its
// producer down.

// tcpSendBufferForPort holds a func(port uint16) int64 or nil. See
// SetTCPSendBufferForPort.
var tcpSendBufferForPort atomic.Value

// SetTCPSendBufferForPort registers a hook consulted for each accepted inbound
// TCP connection, keyed by the local (listening) port. A return value > 0 is
// applied as that endpoint's send buffer size in bytes; 0 leaves the stack
// default in place. A nil hook clears any previous registration.
//
// The size is clamped by the stack's TCPSendBufferSizeRangeOption min/max, and
// setting it explicitly disables netstack's send-buffer auto-tuning for that
// endpoint (gvisor tcp.Endpoint.OnSetSendBufferSize) — which is the point: an
// auto-tuned buffer grows back to the bandwidth-delay product the stack infers,
// undoing the bound the caller asked for.
//
// Keyed by port rather than applied globally so that one latency-sensitive
// listener can be bounded without capping throughput for everything else: a
// small send buffer limits a single connection's throughput to buffer/RTT, which
// is the right trade for a video stream and the wrong one for a file transfer
// over the same stack.
//
// Intended to be called once during setup, before the listener it targets
// accepts anything. It is safe to call at any time — the load and store are
// atomic — but a connection accepted before the call keeps the stack default.
func SetTCPSendBufferForPort(f func(port uint16) int64) {
	if f == nil {
		tcpSendBufferForPort.Store(func(uint16) int64 { return 0 })
		return
	}
	tcpSendBufferForPort.Store(f)
}

// tcpSendBufferSizeForPort returns the registered size for port, or 0 if there
// is no hook or the hook declines the port.
func tcpSendBufferSizeForPort(port uint16) int64 {
	f, _ := tcpSendBufferForPort.Load().(func(uint16) int64)
	if f == nil {
		return 0
	}
	return f(port)
}
