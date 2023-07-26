package cache

import (
	"fmt"
	"log"
	"sync"

	pb "gk-cache/cachepb"
	"gk-cache/peers"
	"gk-cache/singleflight"
)

// A Group is a cache namespace and associated data loaded spread over
type Group struct {
	name      string
	getter    Getter
	mainCache cache
	peers     peers.PeerPicker
	loader    *singleflight.Group
}

// A Getter loads data for a key.
type Getter interface {
	Get(key string) ([]byte, error)
}

// A GetterFunc implements Getter with a function.
type GetterFunc func(key string) ([]byte, error)

// Get implements Getter interface function
func (f GetterFunc) Get(key string) ([]byte, error) {
	return f(key)
}

var (
	mu     sync.RWMutex
	groups = make(map[string]*Group)
)

// NewGroup create a new instance of Group
func NewGroup(name string, cacheBytes int64, getter Getter) *Group {
	if getter == nil {
		panic("nil Getter")
	}

	mu.Lock()
	defer mu.Unlock()

	g := &Group{
		name:      name,
		getter:    getter,
		mainCache: cache{cacheBytes: cacheBytes},
		loader:    &singleflight.Group{},
	}

	groups[name] = g
	return g
}

// GetGroup returns the named group previously created with NewGroup, or
// nil if there's no such group.
func GetGroup(name string) *Group {
	mu.RLock()
	defer mu.RUnlock()
	return groups[name]
}

// Get value for a key from cache
func (g *Group) Get(key string) (ByteView, error) {
	if key == "" {
		return ByteView{}, fmt.Errorf("key is required")
	}

	if v, ok := g.mainCache.Get(key); ok {
		log.Println("[GkCache] hit")
		return v, nil
	}

	return g.load(key)
}

func (g *Group) RegisterPeers(peers peers.PeerPicker) {
	if peers == nil {
		panic("RegisterPeerPicker called more than once")
	}

	g.peers = peers
}

func (g *Group) load(key string) (value ByteView, err error) {
	viewi, err := g.loader.Do(key, func() (interface{}, error) {
		if g.peers != nil {
			if peer, ok := g.peers.PickPeer(key); ok {
				if value, err = g.getFromPeer(peer, key); err == nil {
					return value, nil
				}

				log.Println("[GkCache] Failed to get from peer", err)
			}
		}
		return g.getLocally(key)
	})

	if err != nil {
		return viewi.(ByteView), err
	}

	return
}

func (g *Group) getLocally(key string) (ByteView, error) {
	bytes, err := g.getter.Get(key)
	if err != nil {
		return ByteView{}, err

	}

	value := ByteView{b: cloneBytes(bytes)}
	g.populateCache(key, value)

	return value, nil
}

func (g *Group) populateCache(key string, value ByteView) {
	g.mainCache.Add(key, value)
}

func (g *Group) getFromPeer(peer peers.PeerGetter, key string) (ByteView, error) {

	req := &pb.Request{
		Group: g.name,
		Key:   key,
	}
	resp := &pb.Response{}
	err := peer.Get(req, resp)
	if err != nil {
		return ByteView{}, err
	}

	return ByteView{resp.Value}, nil
}
