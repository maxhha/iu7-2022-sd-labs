package interactors

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type RoomSuite struct {
	InteractorSuite
	interactor RoomInteractor
}

func TestRoomSuite(t *testing.T) {
	suite.Run(t, new(RoomSuite))
}

func (s *RoomSuite) SetupTest() {
	s.InteractorSuite.SetupTest()

	s.interactor = NewRoomInteractor(s.repo.OrganizerMock, s.repo.RoomMock)
}

func (s *RoomSuite) TestCreate() {
	roomID := "test-room"
	organizerID := "test-organizer"
	organizer := *s.NewOrganizerPtr().SetID(organizerID)
	name := "test-name"
	address := "test-address"

	type Case struct {
		Name   string
		Mock   func(c *Case)
		Result entities.Room
		Error  string
	}

	cases := []Case{{
		"Case: fail get organizer",
		func(c *Case) {
			s.repo.OrganizerMock.On("Get", organizerID).
				Return(entities.Organizer{}, repositories.ErrNotFound).
				Once()
		},
		entities.Room{},
		"organizer repo get: not found",
	}, {
		"Case: success",
		func(c *Case) {
			s.repo.OrganizerMock.On("Get", organizerID).
				Return(organizer, nil).
				Once()

			s.repo.RoomMock.On("Create", mock.Anything).
				Return(func(room *entities.Room) error {
					room.SetID(roomID)
					return nil
				}).
				Once()
		},
		*s.NewRoomPtr().
			SetID(roomID).
			SetOrganizerID(organizerID).
			SetName(name).
			SetAddress(address),
		"",
	}}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.interactor.Create(organizerID, name, address)

		if len(c.Error) > 0 {
			require.ErrorContains(s.T(), err, c.Error, c.Name)
		} else {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		}
	}
}

func (s *RoomSuite) TestFind() {
	n := 10
	rooms := make([]entities.Room, n)
	for i := 0; i < n; i++ {
		rooms[i].SetID(fmt.Sprintf("room%d", i))
	}

	params := repositories.RoomFindParams{}
	s.repo.RoomMock.On("Find", &params).Return(rooms, nil)

	result, err := s.interactor.Find(&params)
	require.NoError(s.T(), err)
	require.Equal(s.T(), rooms, result)
}

func (s *RoomSuite) TestDelete() {
	id := "test-room"
	s.repo.RoomMock.On("Delete", id).Return(entities.Room{}, nil)

	err := s.interactor.Delete(id)
	require.NoError(s.T(), err)
}
