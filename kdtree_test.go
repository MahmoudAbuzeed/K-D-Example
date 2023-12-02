package kdtree

import (
	"fmt"
	"testing"
)

type point2D struct {
	values [2]float64
}

func (p point2D) GetDimensionValue(n int) float64 {
	return p.values[n]
}

func (p point2D) Dimensions() int {
	return len(p.values)
}

func (p point2D) String() string {
	return fmt.Sprintf("<%v,%v>", p.values[0], p.values[1])
}

type testCase struct {
	tree     *KDTree[float64]
	target   point2D
	expected point2D
}

func TestNew_KDTree(t *testing.T) {
	points := []KDPoint[float64]{
		point2D{values: [2]float64{5, 4}},
		point2D{values: [2]float64{2, 6}},
		point2D{values: [2]float64{13, 3}},
		point2D{values: [2]float64{3, 1}},
		point2D{values: [2]float64{10, 2}},
		point2D{values: [2]float64{8, 7}},
	}
	tree := NewKDTree[float64](points, func(a, b KDPoint[float64], dim int) float64 {
		return a.GetDimensionValue(dim) - b.GetDimensionValue(dim)
	})

	traverse(tree.Root, 0, func(node *Node[float64], depth int) {
		dim := depth % node.Point.Dimensions()
		if node.Left != nil && node.Left.Point.GetDimensionValue(dim) > node.Point.GetDimensionValue(dim) {
			t.Errorf("expected left node to be less than %v, got %v", node.Point, node.Left.Point)
		}
		if node.Right != nil && node.Right.Point.GetDimensionValue(dim) < node.Point.GetDimensionValue(dim) {
			t.Errorf("expected right node to be greater than %v, got %v", node.Point, node.Right.Point)
		}
		fmt.Printf("depth: %v, node: %v\n", depth, node.Point)
	})

}

func TestSearch_Nearest(t *testing.T) {
	points := []KDPoint[float64]{
		point2D{values: [2]float64{5, 4}},
		point2D{values: [2]float64{2, 6}},
		point2D{values: [2]float64{13, 3}},
		point2D{values: [2]float64{3, 1}},
		point2D{values: [2]float64{10, 2}},
		point2D{values: [2]float64{8, 7}},
	}
	tree := NewKDTree[float64](points, func(a, b KDPoint[float64], dim int) float64 {
		return a.GetDimensionValue(dim) - b.GetDimensionValue(dim)
	})

	testCases := []testCase{
		{
			tree:     tree,
			target:   point2D{values: [2]float64{9, 4}},
			expected: point2D{values: [2]float64{10, 2}},
		},
		{
			tree:     tree,
			target:   point2D{values: [2]float64{9, 6}},
			expected: point2D{values: [2]float64{8, 7}},
		},
		{
			tree:     tree,
			target:   point2D{values: [2]float64{9, 4.5}},
			expected: point2D{values: [2]float64{8, 7}},
		},
		{
			tree:     tree,
			target:   point2D{values: [2]float64{9, 4.5}},
			expected: point2D{values: [2]float64{10, 2}},
		},
		{
			tree:     tree,
			target:   point2D{values: [2]float64{3, 1}},
			expected: point2D{values: [2]float64{3, 1}},
		},
		{
			tree:     tree,
			target:   point2D{values: [2]float64{0, 0}},
			expected: point2D{values: [2]float64{3, 1}},
		},
	}

	for _, tc := range testCases {
		actual := tc.tree.SearchNearest(&tc.target)
		if actual != tc.expected && distance(tc.target, tc.expected, tree.dstFn) != distance(tc.target, actual, tree.dstFn) {
			t.Errorf("expected %v, got %v", tc.expected, actual)
		}
	}

}

func TestInsert(t *testing.T) {
	points := []KDPoint[float64]{
		point2D{values: [2]float64{5, 4}},
		point2D{values: [2]float64{2, 6}},
		point2D{values: [2]float64{13, 3}},
		point2D{values: [2]float64{3, 1}},
		point2D{values: [2]float64{10, 2}},
		point2D{values: [2]float64{8, 7}},
	}
	tree := NewKDTree[float64](points, func(a, b KDPoint[float64], dim int) float64 {
		return a.GetDimensionValue(dim) - b.GetDimensionValue(dim)
	})

	newPoint := point2D{values: [2]float64{1, 1}}
	tree.Insert(newPoint)
	if tree.Root.Left.Left.Left.Point != newPoint {
		t.Errorf("expected %v, got %v", newPoint, tree.Root.Left.Left.Left)
	}

	newPoint = point2D{values: [2]float64{3, 5}}
	tree.Insert(newPoint)
	if tree.Root.Left.Right.Right.Point != newPoint {
		t.Errorf("expected %v, got %v", newPoint, tree.Root.Left.Left.Left)
	}

	newPoint = point2D{values: [2]float64{1, 5}}
	tree.Insert(newPoint)
	if tree.Root.Left.Right.Left.Point != newPoint {
		t.Errorf("expected %v, got %v", newPoint, tree.Root.Left.Left.Left)
	}

	newPoint = point2D{values: [2]float64{11, 4}}
	tree.Insert(newPoint)
	if tree.Root.Right.Right.Point != newPoint {
		t.Errorf("expected %v, got %v", newPoint, tree.Root.Left.Left.Left)
	}

	newPoint = point2D{values: [2]float64{10, 4}}
	tree.Insert(newPoint)
	if tree.Root.Right.Right.Left.Point != newPoint {
		t.Errorf("expected %v, got %v", newPoint, tree.Root.Left.Left.Left)
	}

	newPoint = point2D{values: [2]float64{11, 2}}
	tree.Insert(newPoint)
	if tree.Root.Right.Left.Right.Point != newPoint {
		t.Errorf("expected %v, got %v", newPoint, tree.Root.Left.Left.Left)
	}
}
