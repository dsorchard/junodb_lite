package cluster

import (
	"junodb_lite/pkg/y_conn_mgr"
	"time"
)

func (p *OutboundSSProcessor) OnConnectSuccess(conn io.Conn, connector *io.OutboundConnector, timeTaken time.Duration) {

}

func (p *OutboundSSProcessor) OnConnectError(timeTaken time.Duration, connStr string, err error) {

}
