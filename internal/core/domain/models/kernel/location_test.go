package kernel

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNewValidLocation(t *testing.T) {
	location, err := NewLocation(3, 7)

	assert.NoError(t, err)
	assert.NotEmpty(t, location)
}

func TestNewInvalidLocation(t *testing.T) {
	location, err := NewLocation(-1, 1)

	assert.ErrorIs(t, err, ErrValueIsOutOfRange)
	assert.Equal(t, location, Location{})
}

func TestCalculateCorrectDistance(t *testing.T) {
	firstLocation, _ := NewLocation(1, 1)
	secondLocation, _ := NewLocation(5, 5)
	distance, err := firstLocation.DistanceTo(secondLocation)

	assert.Equal(t, distance, 8)
	assert.NoError(t, err)
}

func TestCalculateDistanceInvalidTarget(t *testing.T) {
	firstLocation, _ := NewLocation(1, 1)

	distance, err := firstLocation.DistanceTo(Location{})

	assert.ErrorIs(t, err, ErrInvalidLocation)
	assert.Equal(t, distance, 0)
}

func TestRandomGeneration(t *testing.T) {
	for range 1000 {
		location, _ := NewRandomLocation()

		assert.False(t, location.IsEmpty())
		assert.GreaterOrEqual(t, location.X(), minX)
		assert.LessOrEqual(t, location.X(), maxX)
		assert.GreaterOrEqual(t, location.Y(), minY)
		assert.LessOrEqual(t, location.Y(), maxY)

	}
}
