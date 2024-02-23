package cluster

import (
	"junodb_lite/pkg/y_conn_mgr"
	"time"
)

func (p *OutboundSSProcessor) OnConnectSuccess(conn y_conn_mgr.Conn, connector *y_conn_mgr.OutboundConnector, timeTaken time.Duration) {

}

func (p *OutboundSSProcessor) OnConnectError(timeTaken time.Duration, connStr string, err error) {

}
