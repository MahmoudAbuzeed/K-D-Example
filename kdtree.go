package kdtree

import (
	"fmt"
	"math"
	"sort"
	"strings"
)

type KDPoint[T any] interface {
	GetDimensionValue(n int) T
	Dimensions() int
}

type KDistanceCalculator[T any] func(a, b KDPoint[T], dim int) float64

type Node[T any] struct {
	Point KDPoint[T]
	Left  *Node[T]
	Right *Node[T]
}

type KDTree[T any] struct {
	Root  *Node[T]
	Size  int
	dstFn KDistanceCalculator[T]
}

func NewKDTree[T any](points []KDPoint[T], dstFn KDistanceCalculator[T]) *KDTree[T] {
	if dstFn == nil {
		panic("dstFn cannot be nil")
	}

	return &KDTree[T]{dstFn: dstFn, Root: buildTree(points, dstFn, 0), Size: len(points)}
}

// To Implement

func buildTree[T any](points []KDPoint[T], dstFn KDistanceCalculator[T], depth int) *Node[T] {
	if len(points) == 0 {
		return nil
	}

	// Select axis based on depth
	axis := depth % points[0].Dimensions()

	// Sort points by the selected axis
	sort.Slice(points, func(i, j int) bool {
		return dstFn(points[i], points[j], axis) < 0
	})

	// Find median
	median := len(points) / 2

	// Create a new node
	return &Node[T]{
		Point: points[median],
		Left:  buildTree(points[:median], dstFn, depth+1),
		Right: buildTree(points[median+1:], dstFn, depth+1),
	}
}

func (t *KDTree[T]) SearchNearest(target KDPoint[T]) KDPoint[T] {
	return searchNearest(t.Root, target, 0, t.dstFn, nil, math.MaxFloat64).Point
}

func searchNearest[T any](node *Node[T], target KDPoint[T], depth int, dstFn KDistanceCalculator[T], bestNode *Node[T], bestDist float64) *Node[T] {
	if node == nil {
		return bestNode
	}

	axis := depth % target.Dimensions()
	dist := distance(target, node.Point, dstFn)
	var nextNode, otherNode *Node[T]

	if dstFn(target, node.Point, axis) < 0 {
		nextNode, otherNode = node.Left, node.Right
	} else {
		nextNode, otherNode = node.Right, node.Left
	}

	bestNode = searchNearest(nextNode, target, depth+1, dstFn, bestNode, bestDist)
	if dist < bestDist {
		bestDist = dist
		bestNode = node
	}

	// Check if other subtree might contain a closer point
	if math.Pow(dstFn(target, node.Point, axis), 2) < bestDist {
		bestNode = searchNearest(otherNode, target, depth+1, dstFn, bestNode, bestDist)
	}

	return bestNode
}

func (t *KDTree[T]) Insert(p KDPoint[T]) {
	t.Root = insert(t.Root, p, 0, t.dstFn)
	t.Size++
}

func insert[T any](node *Node[T], point KDPoint[T], depth int, dstFn KDistanceCalculator[T]) *Node[T] {
	if node == nil {
		return &Node[T]{Point: point}
	}

	axis := depth % point.Dimensions()

	if dstFn(point, node.Point, axis) < 0 {
		node.Left = insert(node.Left, point, depth+1, dstFn)
	} else {
		node.Right = insert(node.Right, point, depth+1, dstFn)
	}

	return node
}

// Utils

func (t *KDTree[T]) print() {
	grid := buildTreeGrid(t.Root)
	for _, row := range grid {
		for _, c := range row {
			fmt.Print(c)
		}
		fmt.Println()
		fmt.Println()
	}
}

func maxDepth[T any](n *Node[T]) int {
	if n == nil {
		return 0
	}

	lDepth := maxDepth(n.Left)
	rDepth := maxDepth(n.Right)

	if lDepth > rDepth {
		return lDepth + 1
	}
	return rDepth + 1
}

func buildTreeGrid[T any](root *Node[T]) [][]string {
	if root == nil {
		return [][]string{}
	}

	h := maxDepth(root)
	col := int(math.Pow(2, float64(h+1)) - 1)
	res := make([][]string, h+1)

	for i := 0; i < h+1; i++ {
		row := make([]string, col)
		// init res 2d arr
		for j := 0; j < col; j++ {
			row[j] = ""
		}
		res[i] = row
	}

	maxLen := fillNode(root, 0, col, 0, res)

	for i := 0; i < h+1; i++ {
		for j := 0; j < col; j++ {
			if res[i][j] == "" {
				res[i][j] = strings.Repeat(" ", maxLen)
			}
		}
	}

	return res
}

func fillNode[T any](n *Node[T], l, r, h int, res [][]string) int {
	if n == nil {
		return 1
	}

	maxLen := 0
	var mid int = (l + r) / 2
	if s, ok := n.Point.(fmt.Stringer); ok == true {
		res[h][mid] = s.String()
	} else {
		res[h][mid] = fmt.Sprintf("%v", n.Point)
	}

	if len(res[h][mid]) > maxLen {
		maxLen = len(res[h][mid])
	}

	if n.Left != nil {
		fillNode(n.Left, l, mid, h+1, res)
	}

	if n.Right != nil {
		fillNode(n.Right, mid+1, r, h+1, res)
	}

	return maxLen
}

func traverse(node *Node[float64], depth int, fn func(*Node[float64], int)) {
	if node == nil {
		return
	}

	fn(node, depth)
	traverse(node.Left, depth+1, fn)
	traverse(node.Right, depth+1, fn)
}

func distance[T any](a, b KDPoint[T], dstFn KDistanceCalculator[T]) float64 {
	d := 0.0
	for i := 0; i < a.Dimensions(); i++ {
		d += math.Pow(dstFn(a, b, i), 2)
	}
	return d
}
