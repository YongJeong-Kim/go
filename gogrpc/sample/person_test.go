package sample

import (
	"github.com/stretchr/testify/require"
	"gogrpc/serializer"
	"reflect"
	"testing"
)

func TestNewPerson(t *testing.T) {
	person := NewPerson()
	require.Equal(t, reflect.ValueOf(person).IsZero(), false)

	shirt := person.Shirt
	chino := person.Chino

	shirtJSON, err := serializer.ProtobufToJSON(shirt)
	require.NoError(t, err)
	require.NotEmpty(t, shirtJSON)

	chinoJSON, err := serializer.ProtobufToJSON(chino)
	require.NoError(t, err)
	require.NotEmpty(t, chinoJSON)
}
