package interactors

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type OrganizerSuite struct {
	InteractorSuite
	interactor OrganizerInteractor
}

func TestOrganizerSuite(t *testing.T) {
	suite.Run(t, new(OrganizerSuite))
}

func (s *OrganizerSuite) SetupTest() {
	s.InteractorSuite.SetupTest()

	s.interactor = NewOrganizerInteractor(s.organizerRepo)
}

func (s *OrganizerSuite) TestCreate() {
	name := "user-name"
	id := "user-id"

	s.organizerRepo.On("Create", mock.Anything).Run(func(args mock.Arguments) {
		org, ok := args.Get(0).(*entities.Organizer)
		require.True(s.T(), ok)
		org.SetID(id)
	}).Return(nil)

	org, err := s.interactor.Create(name)
	require.NoError(s.T(), err)
	require.Equal(s.T(), name, org.Name())
	require.Equal(s.T(), id, org.ID())
}

func (s *OrganizerSuite) TestFind() {
	n := 10
	orgs := make([]entities.Organizer, n)
	for i := 0; i < n; i++ {
		orgs[i].SetID(fmt.Sprintf("org%d", i))
	}

	params := repositories.OrganizerFindParams{}
	s.organizerRepo.On("Find", &params).Return(orgs, nil)

	result, err := s.interactor.Find(&params)
	require.NoError(s.T(), err)
	require.Equal(s.T(), orgs, result)
}

func (s *OrganizerSuite) TestUpdate() {
	id := "user-id"
	newName := "user-new-name"

	cases := []struct {
		Name   string
		Mock   func()
		Result entities.Organizer
		Error  error
	}{
		{
			"Get error",
			func() {
				s.organizerRepo.On("Get", id).
					Return(entities.Organizer{}, repositories.ErrNotFound).
					Once()
			},
			entities.Organizer{},
			repositories.ErrNotFound,
		},
		{
			"Success",
			func() {
				org := entities.NewOrganizer()
				org.SetID(id)
				s.organizerRepo.On("Get", id).
					Return(org, nil).
					Once()

				org.SetName(newName)

				s.organizerRepo.On("Update", &org).
					Return(nil)
			},
			*s.NewOrganizerPtr().SetID(id).SetName(newName),
			nil,
		},
	}

	for _, c := range cases {
		c.Mock()
		result, err := s.interactor.Update(&interactors.OrganizerUpdateParams{
			ID:   id,
			Name: newName,
		})

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}
