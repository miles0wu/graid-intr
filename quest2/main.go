package main

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"
)

func initNodes(size int) []*Node {
	nodes := make([]*Node, size)
	for i := 0; i < size; i++ {
		nodes[i] = NewNode(i)
	}

	for i, node := range nodes {
		for j := 0; j < size; j++ {
			if i == j {
				continue
			}
			node.AddPeers(nodes[j])
		}
	}

	return nodes
}

func main() {
	size := 5

	if len(os.Args) > 1 {
		n, err := strconv.Atoi(os.Args[1])
		if err == nil && n >= 3 {
			size = n
		}
	}

	nodes := initNodes(size)

	fmt.Printf("> Starting quorum with %d members\n", len(nodes))

	for _, node := range nodes {
		node.Start()
	}

	time.Sleep(1 * time.Second)

	fmt.Println("> Enter 'kill <ID>' to remove a member or 'exit' to stop:")
	scanner := bufio.NewScanner(os.Stdin)
	for {
		if !scanner.Scan() {
			break
		}

		input := scanner.Text()
		input = strings.ToLower(strings.TrimSpace(input))
		if input == "exit" {
			fmt.Println("> Exiting...")
			return
		}

		parts := strings.Split(input, " ")
		if len(parts) == 2 && parts[0] == "kill" {
			id, err := strconv.Atoi(parts[1])
			if err != nil || id < 0 || id >= len(nodes) || nodes[id] == nil {
				fmt.Println("> Invalid ID")
				continue
			}

			// stop a node by id
			nodes[id].Stop()
			for _, node := range nodes {
				if node != nil {
					node.RemovePeer(nodes[id])
				}
			}
			nodes[id] = nil

			// check alive nodes
			cnt := 0
			for _, node := range nodes {
				if node != nil {
					cnt++
				}
			}
			if cnt < 3 {
				fmt.Printf("> Member cnt is %d (less than 3), exiting...\n", cnt)
				os.Exit(0)
			}
		} else {
			fmt.Println("> Unknown command")
		}
	}
	if err := scanner.Err(); err != nil {
		fmt.Println("> Error reading input:", err)
	}
}
