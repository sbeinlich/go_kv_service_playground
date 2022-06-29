package cli

import (
	"fmt"
	"bufio"
	"os"
	"strings"
)

func Run() {
	// writer_c := make(chan string)
	kvs := make(map[string]string)
	reader := bufio.NewReader(os.Stdin)
	for {
		fmt.Print("read(r), put(p), delete(d), print store(s)? ")
		text, _ := reader.ReadString('\n');
		text = strings.TrimSuffix(text, "\n")
		switch  text {
			case "r":
				fmt.Print("Enter Key: ")
				key, _ := reader.ReadString('\n');
				key = strings.TrimSuffix(key, "\n")
				v, ok := read(&kvs, key)
				if ok {
					fmt.Println(v)
				} else {
					fmt.Println("ERROR: key not found")
				}
			case "p":
				fmt.Print("Enter Key: ")
				key, _ := reader.ReadString('\n');
				key = strings.TrimSuffix(key, "\n")
				fmt.Print("Enter Value: ")
				value, _ := reader.ReadString('\n');
				value = strings.TrimSuffix(value, "\n")
				write(&kvs, key,value)
			case "d":
				fmt.Print("Enter Key: ")
				key, _ := reader.ReadString('\n');
				key = strings.TrimSuffix(key, "\n")
				del(&kvs, key)
			case "s":
				printMap(&kvs)
			default:
				// continue
		}
	}

}

func client1(c chan string) {
	for i := 0; i < 10; i++ {
		c <- "ping!"
	}
}

func printMap(m *map[string]string) {
	for k, v := range *m {
        fmt.Println(k, ":", v)
    }
}

func del(kvs *map[string]string, k string) {
	delete(*kvs, k)
}

func write(kvs *map[string]string, k string, v string) {
	(*kvs)[k] = v
}

func read(kvs *map[string]string, k string) (string, bool) {
	v, ok := (*kvs)[k]
	return v, ok
}

// type KVStore interface {
// 	Get() string
// 	Put() string
// }

// type InMemoryKVStore struct {
// 	kv_map map[string]string
// }

// func (kvs *InMemoryKVStore) Get() string {
// 	return kvs[]
// }
