package consistent_hashing

import (
	"fmt"
	"hash/crc32"
	"sort"
	"strconv"
	"sync"
)

type sortedKeys []uint32

func (sk sortedKeys) Len() int           { return len(sk) }
func (sk sortedKeys) Less(i, j int) bool { return sk[i] < sk[j] }
func (sk sortedKeys) Swap(i, j int)      { sk[i], sk[j] = sk[j], sk[i] }

type ConsistentHashing struct {
	// The number of replicas to increase the load distribution
	virtualNodeNum int

	hashSortedKeys sortedKeys

	hashRing  map[uint32]string
	serverMap map[string]struct{}

	mu sync.RWMutex
}

func NewConsistentHashing(virtualNodeNum int) *ConsistentHashing {
	return &ConsistentHashing{
		virtualNodeNum: virtualNodeNum,

		hashSortedKeys: []uint32{},

		hashRing:  map[uint32]string{},
		serverMap: map[string]struct{}{},
	}
}

// AddServer adds a server to the consistent hashing ring
func (ch *ConsistentHashing) AddServer(name string) {
	_, found := ch.serverMap[name]
	if found {
		return
	}

	hash := ch.hash(name)
	ch.serverMap[name] = struct{}{}
	ch.hashRing[hash] = name

	// Add virtual nodes: simply re-hash the server name concatenated with index
	for i := 0; i < ch.virtualNodeNum; i++ {
		virtualNode := ch.hash(name + strconv.Itoa(i))
		ch.hashRing[virtualNode] = name
	}

	ch.updateHashSortedKeys()
}

// RemoveServer removes the given server from the consistent hashing ring
func (ch *ConsistentHashing) RemoveServer(name string) {
	_, found := ch.serverMap[name]
	if !found {
		return
	}

	hash := ch.hash(name)
	delete(ch.hashRing, hash)
	delete(ch.serverMap, name)

	// Delete virtual nodes
	for i := 0; i < ch.virtualNodeNum; i++ {
		virtualNode := ch.hash(name + strconv.Itoa(i))
		delete(ch.hashRing, virtualNode)
	}

	ch.updateHashSortedKeys()
}

func (ch *ConsistentHashing) updateHashSortedKeys() {
	ch.mu.Lock()
	defer ch.mu.Unlock()

	var tempHashes sortedKeys
	for k := range ch.hashRing {
		tempHashes = append(tempHashes, k)
	}

	sort.Sort(tempHashes)
	ch.hashSortedKeys = tempHashes
}

// ListServers lists all available servers in the consistent hashing ring
func (ch *ConsistentHashing) ListServers() []string {
	var servers []string
	for k := range ch.serverMap {
		servers = append(servers, k)
	}
	return servers
}

// Get gets the servers to store an object with the given key
func (ch *ConsistentHashing) Get(key string) (string, error) {
	if len(ch.serverMap) == 0 {
		return "", fmt.Errorf("no server available")
	}

	hash := ch.hash(key)
	// Find the Nearest server using binary search
	nearest := findNearestIndex(ch.hashSortedKeys, hash)
	if nearest >= len(ch.hashSortedKeys) {
		nearest = 0
	}

	server, found := ch.hashRing[ch.hashSortedKeys[nearest]]
	if !found {
		return "", fmt.Errorf("server not found")
	}
	return server, nil
}

func (ch *ConsistentHashing) hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key))
}

func (ch *ConsistentHashing) printHashRing() {
	fmt.Printf("Hash ring: [ ")
	for _, val := range ch.hashSortedKeys {
		server, _ := ch.hashRing[val]
		fmt.Printf("%s:%+v ", server, val)
	}
	fmt.Printf("]\n")
}

func findNearestIndex(keys sortedKeys, key uint32) int {
	i, j := 0, len(keys)
	for i < j {
		h := (i + j) / 2
		if keys[h] < key {
			i = h + 1
		} else {
			j = h
		}
	}
	return i
}
