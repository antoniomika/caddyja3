package caddyja3

import (
	"net"

	caddy "github.com/caddyserver/caddy/v2"
	"github.com/caddyserver/caddy/v2/caddyconfig/caddyfile"
	"go.uber.org/zap"
)

func init() {
	caddy.RegisterModule(JA3Listener{})
}

type JA3Listener struct {
	cache *Cache
	log   *zap.Logger
}

type tlsClientHelloListener struct {
	net.Listener
	cache *Cache
	log   *zap.Logger
}

// CaddyModule implements caddy.Module.
func (JA3Listener) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID:  "caddy.listeners.ja3",
		New: func() caddy.Module { return new(JA3Listener) },
	}
}

func (l *JA3Listener) Provision(ctx caddy.Context) error {
	a, err := ctx.App(CacheAppId)
	if err != nil {
		return err
	}

	l.cache = a.(*Cache)
	l.log = ctx.Logger(l)
	return nil
}

// WrapListener implements caddy.ListenerWrapper.
func (l *JA3Listener) WrapListener(ln net.Listener) net.Listener {
	return &tlsClientHelloListener{
		ln,
		l.cache,
		l.log,
	}
}

// UnmarshalCaddyfile implements caddyfile.Unmarshaler.
func (l *JA3Listener) UnmarshalCaddyfile(_ *caddyfile.Dispenser) error {
	// no-op impl
	return nil
}

// Accept implements net.Listener.
func (l *tlsClientHelloListener) Accept() (net.Conn, error) {
	conn, err := l.Listener.Accept()
	if err != nil {
		return conn, err
	}

	l.log.Debug("Accepting JA3 conn for " + conn.RemoteAddr().String())

	return NewJA3Conn(conn, l.cache, l.log)
}

// Close implements net.Listener.
func (l *tlsClientHelloListener) Close() error {
	l.cache.Reset()
	l.log.Debug("Clearing JA3 cache")

	return l.Listener.Close()
}

// Interface guards.
var (
	_ caddy.Provisioner     = (*JA3Listener)(nil)
	_ caddy.ListenerWrapper = (*JA3Listener)(nil)
	_ caddyfile.Unmarshaler = (*JA3Listener)(nil)
)
