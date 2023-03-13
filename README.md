# Consistent Hashing
## What is Consistent Hashing
Quoted from Wikipedia: "Consistent hashing is a special kind of hashing such that when a
hash table is re-sized and consistent hashing is used, only k/n keys need to be remapped on
average, where k is the number of keys, and n is the number of slots. In contrast, in most
traditional hash tables, a change in the number of array slots causes nearly all keys to be
remapped [1]"
## How it works
1. Store both objects and servers on the circle(hash ring)
2. To find out which server to ask for the given key we need to locate the key on the circle and move in the ascending angle direction(counter-clockwise) until we find a server

![img.png](img.png)
## Install
```
go get github.com/nguyen1tech/consistent-hashing
```
## Usage
```go
// Initialize an instance of ConsistentHashing with 3 replicas
ch := NewConsistentHashing(3)

// Add servers to the hash ring
ch.AddServer("server01")
ch.AddServer("server02")
ch.AddServer("server03")

// Remove server from the hash ring
ch.RemoveServer("server02")
ch.RemoveServer("server03")

// Retrieve the server in which will store the key-value pair data
// 
// server: The server name
// error: The error may occur
server, err := ch.Get(key)
if err != nil {
    // Handle error
}
```
References:
- https://en.wikipedia.org/wiki/Consistent_hashing
- https://www.toptal.com/big-data/consistent-hashing#:~:text=according%20to%20Wikipedia).-,Consistent%20Hashing%20is%20a%20distributed%20hashing%20scheme%20that%20operates%20independently,without%20affecting%20the%20overall%20system.
- https://medium.com/system-design-blog/consistent-hashing-b9134c8a9062
