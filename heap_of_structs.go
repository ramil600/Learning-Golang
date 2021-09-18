// We will implement heap of the structs(points) and find the point with least distance to some vertex

package main

import (
	"container/heap"
	"fmt"
	"math"
)

//Struct we are going to use in the heap
type Point struct {
	x, y int
}

//the desired vertex. We want to find the closest point to this vertex
var vertex = Point{2, 2}


//Heap we are going to implement from container/heap package
type PointHeap []Point


//Functions that implement Heap Interface: Len, Swap, Less, Push and Pop
func (h PointHeap) Len() int {
	return len(h)
}

func (h PointHeap) Less(x, y int) bool {
	return distToVrtx(h[x]) < distToVrtx(h[y])
}

func (h PointHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *PointHeap) Pop() interface{} {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func (h *PointHeap) Push(x interface{}) {
	p := x.(Point)
	*h = append(*h, p)
}


//Find the distance from a given p to our vertex
func distToVrtx(p Point) float64 {
	sqrx := (p.x - vertex.x) * (p.x - vertex.x);
	sqry := (p.y - vertex.y) * (p.y - vertex.y);
	dist := math.Sqrt(float64(sqrx + sqry));
	return dist
}

func main() {

	points := PointHeap{{6, 7}, {1, 3},
		{4, 5},
	}
	heap.Init(&points)

  //Points should now be sorted. We are listing points from shortest distance to vertex to longest.
	for points.Len() > 0 {
		item := heap.Pop(&points).(Point)
		fmt.Printf("The element of array is point: X: %d, Y: %d. The distance is %.2f\n", item.x, item.y, distToVrtx(item))
	}

}
