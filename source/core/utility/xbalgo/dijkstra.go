package xbalgo

import (
	"cmp"
	"fmt"
	"math"
	"slices"
	"strconv"

	"github.com/starryck/strk-tc-x-lib-go/source/core/toolkit/xbvalue"
	"github.com/starryck/strk-tc-x-lib-go/source/core/utility/xbctnr"
)

func NewDijkstra[T DijkstraVectex]() Dijkstra[T] {
	dijkstra := &Dijkstra[T]{
		vertexSet:    xbctnr.NewSet[T](),
		neighborsMap: make(map[T][]*DijkstraAdjacency[T]),
	}
	return *dijkstra
}

type Dijkstra[T DijkstraVectex] struct {
	vertexSet    xbctnr.Set[T]
	neighborsMap map[T][]*DijkstraAdjacency[T]
}

type DijkstraVectex = comparable

type DijkstraOperate = func()

func (dijkstra *Dijkstra[T]) Graph(edges ...*DijkstraEdge[T]) {
	graphDijkstra(dijkstra, edges...)
}

func (dijkstra *Dijkstra[T]) Calculate(source T, options *DijkstraCalculatorOptions[T]) map[T]*DijkstraRoute[T] {
	return calculateDijkstra(dijkstra, source, options)
}

func graphDijkstra[T DijkstraVectex](dijkstra *Dijkstra[T], edges ...*DijkstraEdge[T]) {
	grapher := &DijkstraGrapher[T]{dijkstra: dijkstra, edges: edges}
	grapher.iterateEdges(func() {
		grapher.checkDistance()
		grapher.storeVerticies()
		grapher.storeNeighbors()
	})
}

type DijkstraGrapher[T DijkstraVectex] struct {
	dijkstra *Dijkstra[T]
	edge     *DijkstraEdge[T]
	edges    []*DijkstraEdge[T]
}

type DijkstraEdge[T DijkstraVectex] struct {
	Source   T
	Target   T
	Distance float64
}

type DijkstraAdjacency[T DijkstraVectex] struct {
	vectex   T
	distance float64
}

func (grapher *DijkstraGrapher[T]) iterateEdges(operate DijkstraOperate) {
	for _, edge := range grapher.edges {
		grapher.edge = edge
		operate()
		grapher.edge = nil
	}
}

func (grapher *DijkstraGrapher[T]) checkDistance() {
	if distance := grapher.edge.Distance; distance < 0 {
		panic(fmt.Sprintf("Edge distance `%s` must be a nonnegative number.", strconv.FormatFloat(distance, 'f', -1, 64)))
	}
}

func (grapher *DijkstraGrapher[T]) storeVerticies() {
	edge := grapher.edge
	vertexSet := grapher.dijkstra.vertexSet
	vertexSet.Add(edge.Source)
	vertexSet.Add(edge.Target)
}

func (grapher *DijkstraGrapher[T]) storeNeighbors() {
	edge := grapher.edge
	neighborsMap := grapher.dijkstra.neighborsMap
	if _, ok := neighborsMap[edge.Source]; !ok {
		neighborsMap[edge.Source] = []*DijkstraAdjacency[T]{}
	}
	neighborsMap[edge.Source] = append(neighborsMap[edge.Source], &DijkstraAdjacency[T]{
		vectex:   edge.Target,
		distance: edge.Distance,
	})
}

func calculateDijkstra[T DijkstraVectex](dijkstra *Dijkstra[T], source T, options *DijkstraCalculatorOptions[T]) map[T]*DijkstraRoute[T] {
	calculator := &DijkstraCalculator[T]{dijkstra: dijkstra, source: source, options: options}
	calculator.initialize()
	calculator.iterateVertexSet(func() {
		calculator.setupPoint()
		calculator.storePoint()
	})
	calculator.iterateHeapValues(func() {
		calculator.visitPoint()
		calculator.iterateNeighbors(func() {
			calculator.setupDistance()
			if calculator.shouldUpdateNeighbor() {
				calculator.updateNeighbor()
			}
		})
	})
	calculator.updateRoutes()
	return calculator.routes
}

type DijkstraCalculator[T DijkstraVectex] struct {
	dijkstra *Dijkstra[T]
	source   T
	heap     *xbctnr.Heap[*DijkstraHeapValue[T]]
	points   map[T]*DijkstraPoint[T]
	routes   map[T]*DijkstraRoute[T]
	options  *DijkstraCalculatorOptions[T]

	vertex   T
	point    *DijkstraPoint[T]
	place    *DijkstraPoint[T]
	distance float64
}

type DijkstraCalculatorOptions[T DijkstraVectex] struct {
	Target      *T
	Rangefinder DijkstraRangefinder[T]
}

type DijkstraPoint[T DijkstraVectex] struct {
	vertex    T
	previous  T
	order     int
	distance  float64
	isVisited bool
}

type DijkstraRoute[T DijkstraVectex] struct {
	Verticies []T
	Order     int
	Distance  float64
}

type DijkstraRangefinder[T DijkstraVectex] interface {
	Measure(source, target T, order int, distance float64) float64
}

func (calculator *DijkstraCalculator[T]) initialize() {
	vertexSet := calculator.dijkstra.vertexSet
	calculator.heap = &xbctnr.Heap[*DijkstraHeapValue[T]]{}
	calculator.points = make(map[T]*DijkstraPoint[T], len(vertexSet))
	calculator.routes = make(map[T]*DijkstraRoute[T], len(vertexSet))
	if calculator.options == nil {
		calculator.options = &DijkstraCalculatorOptions[T]{}
	}
}

func (calculator *DijkstraCalculator[T]) iterateVertexSet(operate DijkstraOperate) {
	vertexSet := calculator.dijkstra.vertexSet
	for vertex := range vertexSet {
		calculator.vertex = vertex
		operate()
		calculator.vertex = xbvalue.Zero[T]()
		calculator.point = nil
	}
}

func (calculator *DijkstraCalculator[T]) setupPoint() {
	vertex := calculator.vertex
	point := &DijkstraPoint[T]{
		vertex:   vertex,
		previous: vertex,
		distance: math.Inf(1),
	}
	if vertex == calculator.source {
		point.distance = 0
	}
	calculator.point = point
}

func (calculator *DijkstraCalculator[T]) storePoint() {
	point := calculator.point
	calculator.heap.Push(&DijkstraHeapValue[T]{
		vertex:   point.vertex,
		distance: point.distance,
	})
	calculator.points[point.vertex] = point
}

func (calculator *DijkstraCalculator[T]) iterateHeapValues(operate DijkstraOperate) {
	for {
		var value *DijkstraHeapValue[T]
		if mValue, ok := calculator.heap.Pull(); !ok {
			break
		} else {
			value = mValue
		}
		if target := calculator.options.Target; target != nil && *target == value.vertex {
			break
		}
		point := calculator.points[value.vertex]
		if point.isVisited {
			continue
		}
		calculator.point = point
		operate()
		calculator.point = nil
	}
}

func (calculator *DijkstraCalculator[T]) visitPoint() {
	point := calculator.point
	point.isVisited = true
}

func (calculator *DijkstraCalculator[T]) iterateNeighbors(operate DijkstraOperate) {
	point := calculator.point
	var neighbors []*DijkstraAdjacency[T]
	if mNeighbors, ok := calculator.dijkstra.neighborsMap[point.vertex]; !ok {
		return
	} else {
		neighbors = mNeighbors
	}
	for _, adjacency := range neighbors {
		place := calculator.points[adjacency.vectex]
		if place.isVisited {
			continue
		}
		calculator.place = place
		calculator.distance = adjacency.distance
		operate()
		calculator.place = nil
		calculator.distance = 0
	}
}

func (calculator *DijkstraCalculator[T]) setupDistance() {
	point := calculator.point
	place := calculator.place
	distance := calculator.distance
	if rangefinder := calculator.options.Rangefinder; rangefinder != nil {
		distance = rangefinder.Measure(point.vertex, place.vertex, point.order, distance)
		if distance < 0 {
			panic(fmt.Sprintf("Measured distance `%s` must be a nonnegative number.", strconv.FormatFloat(distance, 'f', -1, 64)))
		}
	}
	calculator.distance = point.distance + distance
}

func (calculator *DijkstraCalculator[T]) shouldUpdateNeighbor() bool {
	return calculator.distance < calculator.place.distance
}

func (calculator *DijkstraCalculator[T]) updateNeighbor() {
	point := calculator.point
	place := calculator.place
	place.previous = point.vertex
	place.order = point.order + 1
	place.distance = calculator.distance
	calculator.heap.Push(&DijkstraHeapValue[T]{
		vertex:   place.vertex,
		distance: place.distance,
	})
}

func (calculator *DijkstraCalculator[T]) updateRoutes() {
	for _, point := range calculator.points {
		calculator.routes[point.vertex] = &DijkstraRoute[T]{
			Verticies: calculator.makeVerticies(point),
			Order:     point.order,
			Distance:  point.distance,
		}
	}
}

func (calculator *DijkstraCalculator[T]) makeVerticies(point *DijkstraPoint[T]) []T {
	verticies := []T{point.vertex}
	for {
		place := calculator.points[point.previous]
		if place == point {
			break
		}
		verticies = append(verticies, place.vertex)
		point = place
	}
	slices.Reverse(verticies)
	return verticies
}

type DijkstraHeapValue[T DijkstraVectex] struct {
	vertex   T
	distance float64
}

func (value *DijkstraHeapValue[T]) Compare(other *DijkstraHeapValue[T]) int {
	return cmp.Compare(value.distance, other.distance)
}
