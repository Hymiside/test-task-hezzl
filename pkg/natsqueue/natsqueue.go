package natsqueue

import (
	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/nats-io/nats.go"
)

type queue interface {
	Pub(data models.Good) error
	Sub() error
}

type Queue struct {
	Q queue
}

func NewNatsQueue(nc *nats.Conn, repoP *postgres.PostgresRepository) *Queue {
	return &Queue{Q: newNatsQueue(nc, repoP)}
}
