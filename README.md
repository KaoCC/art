
# Art

## QuadTree
Quadtree is a tree data structure in which each node has zero or four children. Here, it is used to create a special kind of encoding that transforms the original image into colorful grids.

The algorithm builds a Quadtree by recursively dividing the image into four subimages. For each subimage, the average color and deviation are computed.
Later a priority queue is utilized along with a traversal algorithm to decide which node should be processed and animated first based on the related deviation.


![](lenna_animated.gif)
![](ckao_animated.gif)