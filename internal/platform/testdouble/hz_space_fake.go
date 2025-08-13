package testdouble

import (
    "context"
    "sync"

    hzint "golang-sba-hazelcast/internal/platform/hazelcast"
)

// FakeSpace is an in-memory implementation of hz Space for tests.
type FakeSpace struct {
    mu   sync.Mutex
    maps map[string]*FakeMap
}

// NewFakeSpace creates a new FakeSpace.
func NewFakeSpace() *FakeSpace {
    return &FakeSpace{maps: make(map[string]*FakeMap)}
}

// GetMap returns a fake map by name.
func (f *FakeSpace) GetMap(_ context.Context, name string) (hzint.Map, error) {
    f.mu.Lock()
    defer f.mu.Unlock()
    if m, ok := f.maps[name]; ok {
        return m, nil
    }
    m := &FakeMap{data: make(map[string]any)}
    f.maps[name] = m
    return m, nil
}

// FakeMap is a simple string-keyed map implementing hz Map interface.
type FakeMap struct {
    mu   sync.RWMutex
    data map[string]any
}

func (m *FakeMap) Put(_ context.Context, key string, value any) (any, error) {
    m.mu.Lock()
    defer m.mu.Unlock()
    old := m.data[key]
    m.data[key] = value
    return old, nil
}

func (m *FakeMap) GetEntrySet(_ context.Context) ([]hzint.KeyValue, error) {
    m.mu.RLock()
    defer m.mu.RUnlock()
    out := make([]hzint.KeyValue, 0, len(m.data))
    for k, v := range m.data {
        out = append(out, hzint.KeyValue{Key: k, Value: v})
    }
    return out, nil
}


