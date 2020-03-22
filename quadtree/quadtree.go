package quadtree

import (
	"container/heap"
	"container/list"
	"fmt"
	"image"
	"image/color"
	"image/draw"
	"math"
)

// Top Left, Top Right, Bottom Left, Bottom Right
const (
	TL = 0
	TR = 1
	BL = 2
	BR = 3
)

// QuadTree represents the quardtree information
type QuadTree struct {
	MaxWidth  int
	MaxHeight int
	root      *treeNode
	colorMap  [][]color.Color
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
	// tree.colorMap = buildColorAccumulation(imageData, tree.MaxWidth, tree.MaxHeight)
}

// TODO : incomplete
func buildColorAccumulation(imageData image.Image, width, height int) [][]color.Color {

	colorMap := make([][]color.Color, height)
	for i := range colorMap {
		colorMap[i] = make([]color.Color, width)
	}

	colorMap[0][0] = imageData.At(0, 0)

	// the first row
	for i := 1; i < width; i++ {
		// colorMap[0][i] = colorMap[0][i-1] + imageData.At(i, 0)
	}

	// the first column
	for j := 1; j < height; j++ {
		// colorMap[j][0] = colorMap[j-1][0] + imageData.At(0, j)
	}

	// the rest
	for j := 1; j < height; j++ {
		for i := 1; i < width; i++ {
			// colorMap[j][i] = imageData.At(i, j) - colorMap[j-1][i] - colorMap[j][i-1] + colorMap[j-1][i-1]
		}
	}

	return colorMap
}

func calculateDiff(median float64) func(uint8) float64 {
	return func(c uint8) float64 {
		return math.Abs(float64(c) - median)
	}
}

func convertColor(c uint8) float64 {
	return float64(c)
}

func averageColor(sumR, sumG, sumB, sumA float64, width, height int) color.Color {
	area := float64(width * height)
	sumR /= area
	sumG /= area
	sumB /= area
	sumA /= area

	return color.NRGBA{R: uint8(sumR), G: uint8(sumG), B: uint8(sumB), A: uint8(sumA)}
}

func accumulate(imageData image.Image, x, y, width, height int, opR, opG, opB, opA func(uint8) float64) (float64, float64, float64, float64) {

	var sumR, sumG, sumB, sumA float64
	for j := y; j < y+height; j++ {
		for i := x; i < x+width; i++ {

			pixel := imageData.At(i, j)
			currentColor := color.NRGBAModel.Convert(pixel).(color.NRGBA)

			sumR += opR(currentColor.R)
			sumG += opG(currentColor.G)
			sumB += opB(currentColor.B)
			sumA += opA(currentColor.A)
		}
	}

	return sumR, sumG, sumB, sumA
}

func buildTree(imageData image.Image, x, y, width, height int) *treeNode {

	node := treeNode{x: x, y: y, width: width, height: height}

	if width == 1 || height == 1 {

		sumR, sumG, sumB, sumA := accumulate(imageData, x, y, width, height, convertColor, convertColor, convertColor, convertColor)

		node.color = averageColor(sumR, sumG, sumB, sumA, width, height)

		// calculate the color diff
		diffR, diffG, diffB, diffA := accumulate(imageData, x, y, width, height,
			calculateDiff(sumR), calculateDiff(sumG), calculateDiff(sumB), calculateDiff(sumA))

		node.diff = (diffR + diffG + diffB + diffA) / float64(width*height)
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
		currentColor := color.NRGBAModel.Convert(child.color).(color.NRGBA)
		childArea := uint32(child.width * child.height)
		sumR += float64(uint32(currentColor.R) * childArea)
		sumG += float64(uint32(currentColor.G) * childArea)
		sumB += float64(uint32(currentColor.B) * childArea)
		sumA += float64(uint32(currentColor.A) * childArea)
	}

	node.color = averageColor(sumR, sumG, sumB, sumA, width, height)

	diffR, diffG, diffB, diffA := accumulate(imageData, x, y, width, height,
		calculateDiff(sumR), calculateDiff(sumG), calculateDiff(sumB), calculateDiff(sumA))

	node.diff = (diffR + diffG + diffB + diffA) / float64(width*height)
	return &node
}

// TODO: build a DP map for color

// For testing
func printInfo(node *treeNode) {
	fmt.Printf("x: %d, y: %d, w: %d, h: %d, c: %d, diff: %f\n", node.x, node.y, node.width, node.height, node.color, node.diff)
}

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

func setImageBuffer(node *treeNode, imageBuffer *image.NRGBA) {
	for j := node.y; j < node.y+node.height; j++ {
		for i := node.x; i < node.x+node.width; i++ {
			(*imageBuffer).Set(i, j, node.color)
		}
	}
}

// CreateImages : Create images based of the result in the QuadTree
// stepCount : number of steps before stop
// isAnimated : true means create a series of images, false will only create the last result
func (tree *QuadTree) CreateImages(stepCount int, isAnimated bool) []image.Image {

	if stepCount <= 0 {
		return nil
	}

	outputBound := image.Rect(0, 0, int(tree.MaxWidth), int(tree.MaxHeight))
	imageBuffer := image.NewNRGBA(outputBound)

	result := make([]image.Image, 0)

	pq := BuildPQ()
	heap.Push(&pq, tree.root)

	for pq.Len() > 0 && stepCount > 0 {
		currentNode := heap.Pop(&pq).(*treeNode)

		// printInfo(currentNode)
		setImageBuffer(currentNode, imageBuffer)

		// write to a new image and append the result to the list
		if isAnimated {
			currentImage := image.NewNRGBA(outputBound)
			draw.Draw(currentImage, currentImage.Bounds(), imageBuffer, image.ZP, draw.Src)
			result = append(result, currentImage)
		}

		stepCount--

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
