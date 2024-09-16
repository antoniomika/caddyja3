package caddyja3

import (
	"bytes"
	"encoding/binary"
	"errors"
	"io"
	"net"

	"go.uber.org/zap"
)

// JA3Conn a struct to hold all components for our JA3 data.
type JA3Conn struct {
	net.Conn
	cache     *Cache
	logger    *zap.Logger
	reader    *bytes.Reader
	helloRead bool
}

// NewJA3Conn returns a new JA3 conn for generating the JA3 data.
func NewJA3Conn(c net.Conn, cache *Cache, logger *zap.Logger) (net.Conn, error) {
	if c == nil {
		return nil, errors.New("cannot rewind nil connection")
	}

	if cache == nil {
		return nil, errors.New("cache cannot be nil")
	}

	if logger == nil {
		return nil, errors.New("logger cannot be nil")
	}

	return &JA3Conn{
		Conn:   c,
		logger: logger,
		cache:  cache,
	}, nil
}

// Read is the read handler for JA3 conns. It tries to read the client hello on first request.
func (c *JA3Conn) Read(b []byte) (int, error) {
	if !c.helloRead {
		raw := make([]byte, 5)
		if _, err := io.ReadFull(c.Conn, raw); err != nil {
			c.helloRead = true
			return 0, err
		}

		// Check if the first byte is 0x16 (TLS Handshake)
		if raw[0] != 0x16 {
			err := errors.New("not a TLS handshake record")
			c.helloRead = true
			return 0, err
		}

		// Read exactly length bytes from the reader
		raw = append(raw, make([]byte, binary.BigEndian.Uint16(raw[3:5]))...)
		_, err := io.ReadFull(c.Conn, raw[5:])

		if err == nil {
			if err := c.cache.SetClientHello(c.RemoteAddr().String(), raw); err != nil {
				c.logger.Error("Failed to cache JA3 for "+c.RemoteAddr().String(), zap.Error(err))
			}

			c.logger.Debug("Cached JA3 for " + c.RemoteAddr().String())
		} else {
			c.logger.Debug("Failed to read ClientHello for "+c.RemoteAddr().String(), zap.Error(err))
		}
		c.helloRead = true
		c.reader = bytes.NewReader(raw)
	}

	if c.reader == nil || c.reader.Size() == 0 {
		return c.Conn.Read(b)
	}

	n, err := c.reader.Read(b)
	if errors.Is(err, io.EOF) {
		c.reader.Reset([]byte{})
		return n, nil
	}
	return n, err
}

// Close cleans up the ja3 entry from cache.
func (c *JA3Conn) Close() error {
	c.cache.ClearJA3(c.RemoteAddr().String())
	c.logger.Debug("Cleaning JA3 for " + c.RemoteAddr().String())
	return c.Conn.Close()
}

// Interface guards.
var (
	_ net.Conn = (*JA3Conn)(nil)
)
