package natsqueue

import (
	"fmt"
    "encoding/json"

	"github.com/Hymiside/test-task-hezzl/pkg/models"
	"github.com/Hymiside/test-task-hezzl/pkg/repository/postgres"
	"github.com/nats-io/nats.go"
)

type natsQueue struct {
	nc *nats.Conn
	repoP  *postgres.PostgresRepository
}

func newNatsQueue(nc *nats.Conn, repoP *postgres.PostgresRepository) *natsQueue {
	return &natsQueue{nc: nc, repoP: repoP}
}

var logs [][]byte

func (n *natsQueue) Pub(data models.Good) error {
	byteData, err := json.Marshal(data)
	if err != nil {
		return fmt.Errorf("marshal data to byte array: %v", err)
	}

	return n.nc.Publish("hezzl", byteData)
}

func (n *natsQueue) Sub() error {
	_, err := n.nc.Subscribe("hezzl", func(msg *nats.Msg) {
		if len(logs) < 24 {
			logs = append(logs, msg.Data)
		} else {
			if err := n.repoP.Shop.WriteLogs(logs); err != nil {
				fmt.Println(err.Error())
			}
			logs = nil
		}
	})
	if err != nil {
		return err
	}
	return nil
}