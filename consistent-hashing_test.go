package consistent_hashing

import (
	"fmt"
	"testing"
)

func TestConsistentHashing_Set(t *testing.T) {
	// server1: 2065571387, 1168844827, 850139277, 2879736119
	// server2: 3793178497, 1854307288, 427903822, 2156403444
	// server3: 2500886295, 2006796953, 10123791,  2576559029

	// Sorted servers: [10123791, 427903822, 850139277, 1168844827, 1854307288, 2006796953, 2065571387, 2156403444, 2500886295, 2576559029, 2879736119, 3793178497]
	//					server3    server2    server1     server1    server2      server3     server1     server2     server3     server3     server1     server2
	ch := NewConsistentHashing(3)

	ch.AddServer("server01")
	ch.AddServer("server02")
	ch.AddServer("server03")

	ch.printHashRing()

	// simply: 3581567309, text: 999008199, typesetting: 1167078435, long established fact: 1495964601
	keys := map[string]string{
		//"simply":                "server02",
		"text":                  "server01",
		"typesetting":           "server01",
		"long established fact": "server02",
		"English":               "server03",
	}

	for key, value := range keys {
		server, err := ch.Get(key)
		if err != nil {
			t.Errorf("want nil error but got: %+v", err)
			continue
		}

		if server != value {
			t.Errorf("want: %s but got: %s for key: %s", value, server, key)
		}
	}
}

func TestConsistentHashing_AddServer(t *testing.T) {
	ch := NewConsistentHashing(3)
	ch.AddServer("server01")
	ch.AddServer("server02")

	fmt.Printf("All available servers in the hashing ring: %+v\n", ch.ListServers())
	if len(ch.ListServers()) != 2 {
		t.Errorf("want 2 servers added but got %d servers", len(ch.ListServers()))
	}

	if len(ch.hashSortedKeys) != 8 {
		t.Errorf("want 8 hashes in the hashSortedKeys but got %d", len(ch.hashSortedKeys))
	}
}

func TestConsistentHashing_RemoveServer(t *testing.T) {
	ch := NewConsistentHashing(3)
	ch.AddServer("server01")
	ch.AddServer("server02")
	ch.RemoveServer("server01")

	fmt.Printf("All available servers in the hashing ring: %+v\n", ch.ListServers())
	if len(ch.ListServers()) != 1 {
		t.Errorf("want 2 servers added but got %d servers", len(ch.ListServers()))
	}

	if len(ch.hashSortedKeys) != 4 {
		t.Errorf("want 4 hashes in the hashSortedKeys but got %d", len(ch.hashSortedKeys))
	}
}
