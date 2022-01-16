package service_test

import (
	"context"
	"github.com/google/uuid"
	"github.com/stretchr/testify/require"
	"gogrpc/pb"
	"gogrpc/sample"
	"gogrpc/service"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"testing"
)

func TestCreatePerson(t *testing.T) {
	t.Parallel()

	invalidIDPerson := sample.NewPerson()
	invalidIDPerson.Id = "invalidID"

	duplicatedPerson := sample.NewPerson()
	duplicatedPerson.Id = uuid.NewString()
	duplicatedStore := service.NewInMemoryPersonStore()
	err := duplicatedStore.Save(duplicatedPerson)
	require.NoError(t, err)

	testCases := []struct {
		name   string
		person *pb.Person
		store  service.PersonStore
		code   codes.Code
	}{
		{
			name:   "success",
			person: sample.NewPerson(),
			store:  service.NewInMemoryPersonStore(),
			code:   codes.OK,
		},
		{
			name:   "invalid id person",
			person: invalidIDPerson,
			store:  service.NewInMemoryPersonStore(),
			code:   codes.InvalidArgument,
		},
		{
			name:   "duplicated person",
			person: duplicatedPerson,
			store:  duplicatedStore,
			code:   codes.AlreadyExists,
		},
	}

	for i := range testCases {
		tc := testCases[i]

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			req := &pb.CreatePersonRequest{
				Person: tc.person,
			}

			server := service.NewPersonServer(tc.store)
			res, err := server.CreatePerson(context.Background(), req)
			if tc.code == codes.OK {
				require.NoError(t, err)
				require.NotNil(t, res)
				require.NotEmpty(t, res.PersonId)
				if len(tc.person.Id) > 0 {
					require.Equal(t, tc.person.Id, res.PersonId)
				}
			} else {
				require.Error(t, err)
				require.Nil(t, res)
				st, ok := status.FromError(err)
				require.True(t, ok)
				require.Equal(t, st.Code(), tc.code)
			}
		})
	}
}
