package quadtree

import (
	"container/heap"
	"container/list"
	"fmt"
	"image"
	"image/color"
	"image/draw"
)

// Top Left, Top Right, Bottom Left, Bottom Right
const (
	TL = 0
	TR = 1
	BL = 2
	BR = 3
)

// QuadTree keeps the information of the tree
type QuadTree struct {
	MaxWidth  int
	MaxHeight int
	root      *treeNode
}

type treeNode struct {
	x        int
	y        int
	width    int
	height   int
	color    color.Color
	diff     float64
	quadrant [4]*treeNode
}

// BuildTree : construct the QuadTree
func (tree *QuadTree) BuildTree(imageData image.Image) {
	maxBound, minBound := imageData.Bounds().Max, imageData.Bounds().Min
	tree.MaxWidth = (maxBound.X - 1) - minBound.X + 1 // Range: [min, max)
	tree.MaxHeight = (maxBound.Y - 1) - minBound.Y + 1

	xStart, yStart := minBound.X, minBound.Y
	tree.root = buildTree(imageData, xStart, yStart, tree.MaxWidth, tree.MaxHeight)
}

func calculateDiff(median float64) func(uint32) float64 {
	return func(c uint32) float64 {
		diff := float64(c) - median
		return diff * diff
	}
}

func convertColor(c uint32) float64 {
	return float64(c)
}

func averageColor(sumR, sumG, sumB, sumA float64, width, height int) color.Color {
	area := float64(width * height)
	sumR /= area
	sumG /= area
	sumB /= area
	sumA /= area

	return color.NRGBA64{R: uint16(sumR), G: uint16(sumG), B: uint16(sumB), A: uint16(sumA)}
}

func accumulate(imageData image.Image, x, y, width, height int, opR, opG, opB, opA func(uint32) float64) (float64, float64, float64, float64) {

	var sumR, sumG, sumB, sumA float64
	for j := y; j < y+height; j++ {
		for i := x; i < x+width; i++ {

			pixel := imageData.At(i, j)
			r, g, b, a := pixel.RGBA()

			sumR += opR(r)
			sumG += opG(g)
			sumB += opB(b)
			sumA += opA(a)
		}
	}

	return sumR, sumG, sumB, sumA
}

func buildTree(imageData image.Image, x, y, width, height int) *treeNode {

	node := treeNode{x: x, y: y, width: width, height: height}
	area := float64(width * height)

	if width == 1 || height == 1 {

		sumR, sumG, sumB, sumA := accumulate(imageData, x, y, width, height, convertColor, convertColor, convertColor, convertColor)

		node.color = averageColor(sumR, sumG, sumB, sumA, width, height)

		// calculate the color diff
		diffR, diffG, diffB, diffA := accumulate(imageData, x, y, width, height,
			calculateDiff(sumR/area), calculateDiff(sumG/area), calculateDiff(sumB/area), calculateDiff(sumA/area))

		node.diff = (diffR + diffG + diffB + diffA) / area
		return &node
	}

	halfW := width / 2
	halfH := height / 2

	node.quadrant[TL] = buildTree(imageData, x, y, halfW, halfH)
	node.quadrant[TR] = buildTree(imageData, x+halfW, y, width-halfW, halfH)
	node.quadrant[BL] = buildTree(imageData, x, y+halfH, halfW, height-halfH)
	node.quadrant[BR] = buildTree(imageData, x+halfW, y+halfH, width-halfW, height-halfH)

	var sumR, sumG, sumB, sumA float64
	for _, child := range node.quadrant {
		r, g, b, a := child.color.RGBA()
		childArea := uint32(child.width * child.height)
		sumR += float64(uint32(r) * childArea)
		sumG += float64(uint32(g) * childArea)
		sumB += float64(uint32(b) * childArea)
		sumA += float64(uint32(a) * childArea)
	}

	node.color = averageColor(sumR, sumG, sumB, sumA, width, height)

	diffR, diffG, diffB, diffA := accumulate(imageData, x, y, width, height,
		calculateDiff(sumR/area), calculateDiff(sumG/area), calculateDiff(sumB/area), calculateDiff(sumA/area))

	node.diff = (diffR + diffG + diffB + diffA) / area
	return &node
}

// For testing
func printInfo(node *treeNode) {
	fmt.Printf("x: %d, y: %d, w: %d, h: %d, c: %d, diff: %f\n", node.x, node.y, node.width, node.height, node.color, node.diff)
}

// Traversal goes through each node with a selected algorithm
func (tree *QuadTree) Traversal() {

	traversalAlgo := treeTraversalLevelOrder
	visitAlgo := printInfo

	traversalAlgo(tree.root, visitAlgo)
}

func treeTraversalPreOrder(node *treeNode, visit func(node *treeNode)) {

	if node == nil {
		return
	}

	visit(node)

	for _, child := range node.quadrant {
		treeTraversalPreOrder(child, visit)
	}
}

func treeTraversalLevelOrder(node *treeNode, visit func(node *treeNode)) {

	if node == nil {
		return
	}

	queue := list.New()

	queue.PushBack(node)

	for queue.Len() > 0 {
		current := queue.Front()
		currentNode := current.Value.(*treeNode)

		visit(currentNode)

		for _, child := range currentNode.quadrant {
			if child != nil {
				queue.PushBack(child)
			}
		}

		queue.Remove(current)
	}

}

func setImageBuffer(node *treeNode, imageBuffer *image.NRGBA64) {
	for j := node.y; j < node.y+node.height; j++ {
		for i := node.x; i < node.x+node.width; i++ {
			(*imageBuffer).Set(i, j, node.color)
		}
	}
}

// CreateImages : Create images based of the result in the QuadTree
// count : Number of steps before stop
// isAnimated : True means create a series of images, false will only create the final result
// samplePeriod: The sample period. When animate flag is set, draw the result every n frame.
func (tree *QuadTree) CreateImages(count uint, isAnimated bool, samplePeriod uint) []image.Image {

	if count == 0 || samplePeriod == 0 {
		return nil
	}

	outputBound := image.Rect(0, 0, tree.MaxWidth, tree.MaxHeight)
	imageBuffer := image.NewNRGBA64(outputBound)

	result := make([]image.Image, 0)

	pq := BuildPQ()
	heap.Push(&pq, tree.root)

	sample := uint(0)

	for pq.Len() > 0 && sample < count {
		currentNode := heap.Pop(&pq).(*treeNode)
		setImageBuffer(currentNode, imageBuffer)

		// write to a new image and append the result to the list
		if isAnimated && sample%samplePeriod == 0 {
			currentImage := image.NewNRGBA64(outputBound)
			draw.Draw(currentImage, currentImage.Bounds(), imageBuffer, image.ZP, draw.Src)
			result = append(result, currentImage)
		}

		sample++

		for _, child := range currentNode.quadrant {
			if child != nil {
				heap.Push(&pq, child)
			}
		}
	}

	if !isAnimated {
		result = append(result, imageBuffer)
	}

	return result
}
