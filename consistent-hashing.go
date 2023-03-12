package consistent_hashing

import (
	"fmt"
	"hash/crc32"
	"sort"
)

type sortedKeys []uint32

func (sk sortedKeys) Len() int           { return len(sk) }
func (sk sortedKeys) Less(i, j int) bool { return sk[i] < sk[j] }
func (sk sortedKeys) Swap(i, j int)      { sk[i], sk[j] = sk[j], sk[i] }

type ConsistentHashing struct {
	hashSortedKeys sortedKeys

	hashRing  map[uint32]string
	serverMap map[string]struct{}
}

func NewConsistentHashing() *ConsistentHashing {
	return &ConsistentHashing{
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
	ch.hashSortedKeys = append(ch.hashSortedKeys, hash)
	ch.hashRing[hash] = name

	// Sort the hash keys for searching the nearest server
	sort.Sort(ch.hashSortedKeys)
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

	pivot := 0
	for index, val := range ch.hashSortedKeys {
		if val == hash {
			pivot = index
		}
	}
	ch.hashSortedKeys = append(ch.hashSortedKeys[0:pivot], ch.hashSortedKeys[pivot:]...)

	// Sort the hash keys for searching the nearest server
	sort.Sort(ch.hashSortedKeys)
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
	fmt.Printf("Hash value of key: %s -> %+v\n", key, hash)
	// Find the Nearest server
	var serverHash uint32
	for _, val := range ch.hashSortedKeys {
		if hash <= val {
			serverHash = val
			break
		}
	}

	if serverHash == uint32(0) {
		serverHash = ch.hashSortedKeys[0]
	}

	server, found := ch.hashRing[serverHash]
	if !found {
		return "", fmt.Errorf("server not found")
	}
	return server, nil
}

func (ch *ConsistentHashing) hash(key string) uint32 {
	return crc32.ChecksumIEEE([]byte(key)) / (10_000_000)
}

func (ch *ConsistentHashing) printHashRing() {
	fmt.Printf("Hash ring: [ ")
	for _, val := range ch.hashSortedKeys {
		server, _ := ch.hashRing[val]
		fmt.Printf("%s:%+v ", server, val)
	}
	fmt.Printf("]\n")
}
