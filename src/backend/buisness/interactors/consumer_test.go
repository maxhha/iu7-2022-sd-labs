package interactors

import (
	"fmt"
	"iu7-2022-sd-labs/buisness/entities"
	"iu7-2022-sd-labs/buisness/ports/bus"
	"iu7-2022-sd-labs/buisness/ports/interactors"
	"iu7-2022-sd-labs/buisness/ports/repositories"
	"testing"

	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
)

type ConsumerSuite struct {
	InteractorSuite
	interactor ConsumerInteractor
}

func TestConsumerSuite(t *testing.T) {
	suite.Run(t, new(ConsumerSuite))
}

func (s *ConsumerSuite) SetupTest() {
	s.InteractorSuite.SetupTest()

	s.interactor = NewConsumerInteractor(s.repo, s.eventBus, s.validatorService)
}

func (s *ConsumerSuite) TestCreate() {
	id := "test-consumer"
	nickname := "conusmer-nick"
	form := map[string]interface{}{
		"hello": "world",
	}

	type Case struct {
		Name   string
		Mock   func(c *Case)
		Result entities.Consumer
		Error  error
	}

	cases := []Case{
		{
			"Case: fail validate form",
			func(c *Case) {
				s.validatorService.On("Validate", mock.Anything).
					Return(entities.ErrIsEmpty).
					Once()
			},
			*s.NewConsumerPtr(),
			entities.ErrIsEmpty,
		},
		{
			"Case: success",
			func(c *Case) {
				s.validatorService.On("Validate", mock.Anything).
					Return(nil).
					Once()

				s.repo.ConsumerMock.On("Create", mock.Anything).
					Return(func(consumer *entities.Consumer) error {
						consumer.SetID(id)
						return nil
					}).
					Once()
			},
			*s.NewConsumerPtr().
				SetID(id).
				SetNickname(nickname).
				SetForm(form),
			nil,
		},
	}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.interactor.Create(nickname, form)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *ConsumerSuite) TestFind() {
	n := 10
	tables := make([]entities.Consumer, n)
	for i := 0; i < n; i++ {
		tables[i].SetID(fmt.Sprintf("table%d", i))
	}

	params := repositories.ConsumerFindParams{}
	s.repo.ConsumerMock.On("Find", &params).Return(tables, nil)

	result, err := s.interactor.Find(&params)
	require.NoError(s.T(), err)
	require.Equal(s.T(), tables, result)
}

func (s *ConsumerSuite) TestUpdate() {
	consumer := entities.NewConsumer()
	consumer.
		SetID("test-id").
		SetNickname("test-nickname")

	newNickname := "test-new-nickname"
	newForm := map[string]interface{}{
		"hello": "world",
	}

	type Case struct {
		Name   string
		Params interactors.ConsumerUpdateParams
		Mock   func(c *Case)
		Result entities.Consumer
		Error  error
	}

	cases := []Case{
		{
			"Case: fail get unknown",
			interactors.ConsumerUpdateParams{
				ID: "unknown-consumer",
			},
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", "unknown-consumer").
					Return(entities.Consumer{}, repositories.ErrNotFound).
					Once()
			},
			*s.NewConsumerPtr(),
			repositories.ErrNotFound,
		},
		{
			"Case: invalid form",
			interactors.ConsumerUpdateParams{
				ID:   consumer.ID(),
				Form: newForm,
			},
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumer.ID()).
					Return(consumer, nil).
					Once()

				s.validatorService.On("Validate", mock.Anything).
					Return(entities.ErrIsEmpty).
					Once()
			},
			*s.NewConsumerPtr(),
			entities.ErrIsEmpty,
		},
		{
			"Case: success",
			interactors.ConsumerUpdateParams{
				ID:       consumer.ID(),
				Nickname: newNickname,
				Form:     newForm,
			},
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumer.ID()).
					Return(consumer, nil).
					Once()

				s.validatorService.On("Validate", mock.Anything).
					Return(nil).
					Once()

				s.repo.ConsumerMock.On("Update", mock.Anything).
					Return(nil).
					Once()
			},
			*s.NewConsumerPtr().
				SetID(consumer.ID()).
				SetNickname(newNickname).
				SetForm(newForm),
			nil,
		},
	}

	for _, c := range cases {
		c.Mock(&c)
		result, err := s.interactor.Update(&c.Params)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
			require.Equal(s.T(), c.Result, result, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *ConsumerSuite) TestEnterRoom() {
	consumerID := "test-consumer"
	roomID := "test-room"

	consumer := entities.NewConsumer()
	consumer.SetID(consumerID)

	room := entities.NewRoom()
	room.SetID(roomID)

	type Case struct {
		Name  string
		Mock  func(c *Case)
		Error error
	}

	cases := []Case{
		{
			"Case: fail get consumer",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(entities.Consumer{}, repositories.ErrNotFound).
					Once()
			},
			repositories.ErrNotFound,
		},
		{
			"Case: fail update room",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, nil).
					Once()

				s.repo.RoomMock.On("Update", roomID, mock.Anything).
					Return(entities.Room{}, repositories.ErrNotFound).
					Once()
			},
			repositories.ErrNotFound,
		},
		{
			"Case: success",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, nil).
					Once()

				s.repo.RoomMock.On("Update", roomID, mock.Anything).
					Return(
						func(
							_ string,
							updateFn func(*entities.Room) error,
						) entities.Room {
							r := room
							err := updateFn(&r)
							require.NoError(s.T(), err, c.Name)
							return r
						},
						nil,
					).
					Once()

				s.eventBus.On("Notify", &bus.EvtConsumerEnteredRoom{
					Consumer: consumer,
					Room: *s.NewRoomPtr().
						SetID(roomID).
						SetConsumerIDs([]string{consumerID}),
				}).Once()
			},
			nil,
		},
	}

	for _, c := range cases {
		c.Mock(&c)
		err := s.interactor.EnterRoom(consumerID, roomID)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}

func (s *ConsumerSuite) TestExitRoom() {
	consumerID := "test-consumer"
	roomID := "test-room"

	consumer := entities.NewConsumer()
	consumer.SetID(consumerID)

	room := entities.NewRoom()
	room.SetID(roomID)
	room.SetConsumerIDs([]string{consumerID})

	type Case struct {
		Name  string
		Mock  func(c *Case)
		Error error
	}

	cases := []Case{
		{
			"Case: fail get consumer",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(entities.Consumer{}, repositories.ErrNotFound).
					Once()
			},
			repositories.ErrNotFound,
		},
		{
			"Case: fail update room",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, nil).
					Once()

				s.repo.RoomMock.On("Update", roomID, mock.Anything).
					Return(entities.Room{}, repositories.ErrNotFound).
					Once()
			},
			repositories.ErrNotFound,
		},
		{
			"Case: success",
			func(c *Case) {
				s.repo.ConsumerMock.On("Get", consumerID).
					Return(consumer, nil).
					Once()

				s.repo.RoomMock.On("Update", roomID, mock.Anything).
					Return(
						func(
							_ string,
							updateFn func(*entities.Room) error,
						) entities.Room {
							r := room
							err := updateFn(&r)
							require.NoError(s.T(), err, c.Name)
							return r
						},
						nil,
					).
					Once()

				s.eventBus.On("Notify", &bus.EvtConsumerExitedRoom{
					Consumer: consumer,
					Room:     *s.NewRoomPtr().SetID(roomID),
				}).Once()
			},
			nil,
		},
	}

	for _, c := range cases {
		c.Mock(&c)
		err := s.interactor.ExitRoom(consumerID, roomID)

		if c.Error == nil {
			require.NoError(s.T(), err, c.Name)
		} else {
			require.ErrorIs(s.T(), err, c.Error, c.Name)
		}
	}
}
