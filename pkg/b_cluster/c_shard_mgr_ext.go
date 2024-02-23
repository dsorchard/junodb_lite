package cluster

import (
	"junodb_lite/pkg/z_io"
	"time"
)

func (p *OutboundSSProcessor) OnConnectSuccess(conn z_io.Conn, connector *z_io.OutboundConnector, timeTaken time.Duration) {

}

func (p *OutboundSSProcessor) OnConnectError(timeTaken time.Duration, connStr string, err error) {

}
