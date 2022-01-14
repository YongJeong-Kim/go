package sample

import (
	"github.com/stretchr/testify/require"
	"gogrpc/pb"
	"reflect"
	"testing"
)

func TestNewPerson(t *testing.T) {
	person := NewPerson()
	shirt := person.Shirt
	chino := person.Chino

	require.Equal(t, reflect.ValueOf(person).IsZero(), false)
	require.Equal(t, reflect.ValueOf(shirt).IsZero(), false)
	require.Equal(t, chino, (*pb.Chino)(nil))
}