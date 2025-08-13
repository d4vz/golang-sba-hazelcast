package hazelcast

import (
    "context"
    "time"

    hz "github.com/hazelcast/hazelcast-go-client"
)

// KeyValue represents a key-value pair from the distributed map.
type KeyValue struct {
    Key   any
    Value any
}

// Map defines the minimal contract used by services.
type Map interface {
    Put(ctx context.Context, key string, value any) (any, error)
    GetEntrySet(ctx context.Context) ([]KeyValue, error)
}

// Space provides access to distributed maps.
type Space interface {
    GetMap(ctx context.Context, name string) (Map, error)
}

// Client wraps hazelcast-go-client and implements Space.
type Client struct {
    hz *hz.Client
}

// New creates and starts a Hazelcast client with the provided configuration.
func New(ctx context.Context, clusterName string, members []string) (*Client, error) {
    cfg := hz.Config{}
    cfg.Cluster.Name = clusterName
    cfg.Cluster.Network.SetAddresses(members...)

    c, err := hz.StartNewClientWithConfig(ctx, cfg)
    if err != nil {
        return nil, err
    }
    return &Client{hz: c}, nil
}

// Shutdown terminates the client connection.
func (c *Client) Shutdown(ctx context.Context) error {
    return c.hz.Shutdown(ctx)
}

// GetMap returns a distributed map reference by name as the Map interface.
func (c *Client) GetMap(ctx context.Context, name string) (Map, error) {
    m, err := c.hz.GetMap(ctx, name)
    if err != nil {
        return nil, err
    }
    return &hzMap{m: m}, nil
}

// WithTimeout returns a derived context with default timeout for I/O operations.
func WithTimeout(parent context.Context) (context.Context, context.CancelFunc) {
    return context.WithTimeout(parent, 10*time.Second)
}

// hzMap adapts hz.Map to the local Map interface.
type hzMap struct {
    m *hz.Map
}

func (h *hzMap) Put(ctx context.Context, key string, value any) (any, error) {
    return h.m.Put(ctx, key, value)
}

func (h *hzMap) GetEntrySet(ctx context.Context) ([]KeyValue, error) {
    es, err := h.m.GetEntrySet(ctx)
    if err != nil {
        return nil, err
    }
    out := make([]KeyValue, 0, len(es))
    for _, kv := range es {
        out = append(out, KeyValue{Key: kv.Key, Value: kv.Value})
    }
    return out, nil
}


