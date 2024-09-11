package main

import (
	"container/heap"
	"math"
)

// Node represents a tile in the grid.
type Node struct {
	X, Y    int     // Position of the node
	G, H, F float64 // Cost values
	Parent  *Node   // Parent node in the path
	Index   int     // Index in the priority queue (needed for heap)
}

// PriorityQueue implements a priority queue for nodes.
type PriorityQueue []*Node

func (pq PriorityQueue) Len() int { return len(pq) }
func (pq PriorityQueue) Less(i, j int) bool {
	return pq[i].F < pq[j].F // Compare nodes by F value
}
func (pq PriorityQueue) Swap(i, j int) {
	pq[i], pq[j] = pq[j], pq[i]
	pq[i].Index = i
	pq[j].Index = j
}

func (pq *PriorityQueue) Push(x interface{}) {
	n := len(*pq)
	node := x.(*Node)
	node.Index = n
	*pq = append(*pq, node)
}

func (pq *PriorityQueue) Pop() interface{} {
	old := *pq
	n := len(old)
	node := old[n-1]
	node.Index = -1
	*pq = old[0 : n-1]
	return node
}

func (pq *PriorityQueue) update(node *Node, g, h, f float64) {
	node.G = g
	node.H = h
	node.F = f
	heap.Fix(pq, node.Index)
}

// Heuristic function (Manhattan distance).
func heuristic(x1, y1, x2, y2 int) float64 {
	return math.Abs(float64(x1-x2)) + math.Abs(float64(y1-y2))
}

// AStar performs the A* pathfinding algorithm.
func (b *Brain) AStar(startX, startY, goalX, goalY int) []*Node {
	openList := &PriorityQueue{}
	heap.Init(openList)

	startNode := &Node{X: startX, Y: startY, G: 0, H: heuristic(startX, startY, goalX, goalY), F: 0}
	heap.Push(openList, startNode)
	closedList := make(map[int]bool)

	for openList.Len() > 0 {
		current := heap.Pop(openList).(*Node)

		if current.X == goalX && current.Y == goalY {
			return reconstructPath(current)
		}

		closedList[current.X*1000+current.Y] = true

		for _, neighbor := range b.getNeighbors(current) {
			if closedList[neighbor.X*1000+neighbor.Y] {
				continue
			}

			tentativeG := current.G + 1 // Assuming uniform cost for each move

			if tentativeG < neighbor.G {
				neighbor.Parent = current
				neighbor.G = tentativeG
				neighbor.H = heuristic(neighbor.X, neighbor.Y, goalX, goalY)
				neighbor.F = neighbor.G + neighbor.H

				if !isInOpenList(neighbor, openList) {
					heap.Push(openList, neighbor)
				} else {
					openList.update(neighbor, neighbor.G, neighbor.H, neighbor.F)
				}
			}
		}
	}

	return nil // Path not found
}

func isInOpenList(node *Node, openList *PriorityQueue) bool {
	for _, n := range *openList {
		if n.X == node.X && n.Y == node.Y {
			return true
		}
	}
	return false
}

func reconstructPath(node *Node) []*Node {
	path := []*Node{}
	for node != nil {
		path = append(path, node)
		node = node.Parent
	}
	reversePath(path)
	return path
}

func reversePath(path []*Node) {
	for i, j := 0, len(path)-1; i < j; i, j = i+1, j-1 {
		path[i], path[j] = path[j], path[i]
	}
}

// getNeighbors returns the neighboring nodes for the current node.
func (b *Brain) getNeighbors(node *Node) []*Node {
	neighbors := []*Node{}
	directions := [][2]int{{0, 1}, {1, 0}, {0, -1}, {-1, 0}}

	for _, dir := range directions {
		nx, ny := node.X+dir[0], node.Y+dir[1]
		if nx >= 0 && ny >= 0 && nx < SIZE_OF_MAP && ny < SIZE_OF_MAP && b.Owner.WorldProvider.CanWalk(nx, ny) {
			neighbors = append(neighbors, &Node{X: nx, Y: ny, G: math.MaxFloat64})
		}
	}

	return neighbors
}

