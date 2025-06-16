package kernel

import (
	"errors"
	"fmt"
	"math/rand"
	"time"
)

type Location struct {
	x int
	y int

	isSet bool
}

const (
	maxX = 10
	minX = 1
	maxY = 10
	minY = 1
)

var (
	randomSource         *rand.Rand
	ErrValueIsOutOfRange = errors.New("value is out of range")
	ErrInvalidLocation   = errors.New("location is invalid")
)

func init() {
	randomSource = rand.New(rand.NewSource(time.Now().UnixNano()))
}

func abs(x int) int {
	if x < 0 {
		return -x
	}

	return x
}

func NewLocation(x int, y int) (Location, error) {
	if x < minX || x > maxX {
		return Location{}, ErrValueIsOutOfRange
	}
	if y < minY || y > maxY {
		return Location{}, ErrValueIsOutOfRange
	}

	return Location{x, y, true}, nil
}

func NewRandomLocation() (Location, error) {
	x := randomSource.Intn(maxX-minX+1) + minX
	y := randomSource.Intn(maxY-minY+1) + minY

	location, err := NewLocation(x, y)
	if err != nil {
		panic(fmt.Sprintf("invalid random location: x=%d, y=%d, err=%v", x, y, err))

	}

	return location, nil

}

func (l Location) X() int {
	return l.x
}

func (l Location) Y() int {
	return l.y
}

func (l Location) Equals(other Location) bool {
	return l.x == other.x && l.y == other.y
}

func (l Location) IsEmpty() bool {
	return !l.isSet
}

func (l Location) DistanceTo(target Location) (int, error) {
	if target.IsEmpty() {
		return 0, ErrInvalidLocation
	}

	return abs(l.x-target.x) + abs(l.y-target.y), nil
}
