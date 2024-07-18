package carrows

import "math"

type ArrowDescriptor struct {
	// Start Points
	Sx, Sy float64

	// Control point for start point
	C1x, C1y float64

	// Control point for end point
	C2x, C2y float64

	// End point
	Ex, Ey float64

	// Angle of end point
	Ae float64

	// Angle of start point
	As float64
}

type Opts struct {
	PadStart            float64
	PadEnd              float64
	ControlPointStretch float64
	AllowedStartSides   []RectSide
	AllowedEndSides     []RectSide
}

func GetArrow(x0, y0, x1, y1 float64, opts *Opts) ArrowDescriptor {
	return getBoxToBoxArrow(x0, y0, 0, 0, x1, y1, 0, 0, opts)
}

func getBoxToBoxArrow(x0, y0, w0, h0, x1, y1, w1, h1 float64, opts *Opts) ArrowDescriptor {
	if opts == nil {
		opts = &Opts{
			PadStart:            0,
			PadEnd:              0,
			ControlPointStretch: 50,
			AllowedStartSides:   []RectSide{Top, Right, Bottom, Left},
			AllowedEndSides:     []RectSide{Top, Right, Bottom, Left},
		}
	}

	startBox := Rect{x0, y0, w0, h0}
	startAtTop := Vec2{
		X: x0 + w0/2,
		Y: y0 - 2*opts.PadStart,
	}
	startAtBottom := Vec2{
		X: x0 + w0/2,
		Y: y0 + h0 + 2*opts.PadStart,
	}
	startAtLeft := Vec2{
		X: x0 - 2*opts.PadStart,
		Y: y0 + h0/2,
	}
	startAtRight := Vec2{
		X: x0 + w0 + 2*opts.PadStart,
		Y: y0 + h0/2,
	}

	endBox := Rect{x1, y1, w1, h1}
	endAtTop := Vec2{x1 + w1/2, y1 - 2*opts.PadEnd}
	endAtBottom := Vec2{
		X: x1 + w1/2,
		Y: y1 + h1 + 2*opts.PadEnd,
	}
	endAtLeft := Vec2{x1 - 2*opts.PadEnd, y1 + h1/2}
	endAtRight := Vec2{
		X: x1 + w1 + 2*opts.PadEnd,
		Y: y1 + h1/2,
	}

	allSides := []RectSide{Top, Right, Bottom, Left}
	var startSides, endSides []RectSide
	if len(opts.AllowedStartSides) > 0 {
		startSides = opts.AllowedStartSides
	} else {
		startSides = allSides
	}
	if len(opts.AllowedEndSides) > 0 {
		endSides = opts.AllowedEndSides
	} else {
		endSides = allSides
	}

	startPoints := map[RectSide]Vec2{
		Top:    startAtTop,
		Right:  startAtRight,
		Bottom: startAtBottom,
		Left:   startAtLeft,
	}
	endPoints := map[RectSide]Vec2{
		Top:    endAtTop,
		Right:  endAtRight,
		Bottom: endAtBottom,
		Left:   endAtLeft,
	}

	shortestDistance := math.Inf(-1)
	bestStartPoint := startAtTop
	bestEndPoint := endAtTop
	bestStartSide := Top
	bestEndSide := Top

	keepOutZone := 15.0
	for startSideID := 0; startSideID < len(startSides); startSideID++ {
		startSide := startSides[startSideID]
		startPoint := startPoints[startSide]
		if isPointInBox(startPoint, growRect(endBox, keepOutZone)) {
			continue
		}

		for endSideID := 0; endSideID < len(endSides); endSideID++ {
			endSide := endSides[endSideID]
			endPoint := endPoints[endSide]

			if isPointInBox(endPoint, growRect(startBox, keepOutZone)) {
				continue
			}

			d := distanceOf(startPoint, endPoint)
			if d < shortestDistance {
				shortestDistance = d
				bestStartPoint = startPoint
				bestEndPoint = endPoint
				bestStartSide = startSide
				bestEndSide = endSide
			}
		}
	}

	cpStart := controlPointOf(bestStartPoint, bestEndPoint, bestStartSide, opts.ControlPointStretch)
	cpEnd := controlPointOf(bestEndPoint, bestStartPoint, bestEndSide, opts.ControlPointStretch)

	return ArrowDescriptor{
		Sx:  bestStartPoint.X,
		Sy:  bestStartPoint.Y,
		C1x: cpStart.X,
		C1y: cpStart.Y,
		C2x: cpEnd.X,
		C2y: cpEnd.Y,
		Ex:  bestEndPoint.X,
		Ey:  bestEndPoint.Y,
		Ae:  float64(angleOf(bestEndSide)),
		As:  float64(angleOf(bestStartSide)),
	}
}
