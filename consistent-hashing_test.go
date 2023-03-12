package consistent_hashing

import (
	"fmt"
	"testing"
)

func TestConsistentHashing_Set(t *testing.T) {
	ch := NewConsistentHashing()

	ch.AddServer("server01")
	ch.AddServer("server02")
	ch.AddServer("server03")

	fmt.Printf("All available servers in the hashing ring: %+v\n", ch.ListServers())
	ch.printHashRing()

	keys := []string{"simply", "text", "typesetting", "long established fact"}

	for _, key := range keys {
		server, err := ch.Get(key)
		if err != nil {
			fmt.Printf("Error while geting key: %s, error: %+v\n", key, err)
			continue
		}

		fmt.Printf("key: %s in server: %s\n", key, server)
	}

	ch.AddServer("server04")
	ch.AddServer("server05")
	fmt.Printf("All available servers in the hashing ring: %+v\n", ch.ListServers())
	ch.printHashRing()
	for _, key := range keys {
		server, err := ch.Get(key)
		if err != nil {
			fmt.Printf("Error while geting key: %s, error: %+v\n", key, err)
			continue
		}

		fmt.Printf("key: %s in server: %s\n", key, server)
	}
}
