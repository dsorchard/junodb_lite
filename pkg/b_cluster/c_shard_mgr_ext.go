package cluster

import (
	"junodb_lite/pkg/z_conn_mgr"
	"time"
)

func (p *OutboundSSProcessor) OnConnectSuccess(conn z_conn_mgr.Conn, connector *z_conn_mgr.OutboundConnector, timeTaken time.Duration) {

}

func (p *OutboundSSProcessor) OnConnectError(timeTaken time.Duration, connStr string, err error) {

}
