package caddyja3

import (
	"github.com/antoniomika/syncmap"
	caddy "github.com/caddyserver/caddy/v2"
	"github.com/open-ch/ja3"
)

const (
	CacheAppId = "ja3.cache"
)

func init() {
	caddy.RegisterModule(Cache{})
}

type Cache struct {
	data *syncmap.Map[string, *ja3.JA3]
}

func (c *Cache) Provision(ctx caddy.Context) error {
	c.data = syncmap.New[string, *ja3.JA3]()
	return nil
}

func (c *Cache) SetClientHello(addr string, ch []byte) error {
	ja3Data, err := ja3.ComputeJA3FromSegment(ch)
	if err != nil {
		return err
	}

	c.data.Store(addr, ja3Data)
	return nil
}

func (c *Cache) ClearJA3(addr string) {
	c.data.Delete(addr)
}

func (c *Cache) Reset() {
	var toDelete []string
	c.data.Range(func(I string, J *ja3.JA3) bool {
		toDelete = append(toDelete, I)
		return true
	})
	for _, v := range toDelete {
		c.data.Delete(v)
	}
}

func (c *Cache) GetJA3(addr string) *ja3.JA3 {
	ja3Data, found := c.data.Load(addr)

	if found {
		return ja3Data
	}

	return nil
}

// CaddyModule implements caddy.Module.
func (Cache) CaddyModule() caddy.ModuleInfo {
	return caddy.ModuleInfo{
		ID: CacheAppId,
		New: func() caddy.Module {
			return &Cache{
				data: syncmap.New[string, *ja3.JA3](),
			}
		},
	}
}

// Start implements caddy.App.
func (c *Cache) Start() error {
	return nil
}

// Stop implements caddy.App.
func (c *Cache) Stop() error {
	return nil
}

// Interface guards.
var (
	_ caddy.App         = (*Cache)(nil)
	_ caddy.Provisioner = (*Cache)(nil)
)
