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

func (q *Queue) DefaultTicketConfig() interface{} {
	return &TicketConfig{
		Reason: "Server is starting...\n\nPlease try to reconnect in a minute.",
	}
}

func (q *Queue) Enter(ticket *queue.Ticket) bool {
	tcfg := ticket.Config.(*TicketConfig)

	msg := tcfg.Reason.TextComponent()
	ticket.Player.Disconnect(msg)

	return true
}
