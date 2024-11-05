package kick

import "github.com/kasefuchs/lazygate/pkg/queue"

var _ queue.Queue = (*Queue)(nil)

const name = "kick"

type Queue struct{}

func (q *Queue) Name() string {
	return name
}

func (q *Queue) Init(_ *queue.InitOptions) error {
	return nil
}

func (q *Queue) Enter(ticket *queue.Ticket) bool {
	msg := ticket.Config.Queue.Kick.Starting.TextComponent()
	ticket.Player.Disconnect(msg)

	return true
}
