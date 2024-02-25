package redist

import (
	"context"
	"io"
	proto "junodb_lite/pkg/ac_proto"
	. "junodb_lite/pkg/y_conn_mgr"
	redistst "junodb_lite/pkg/y_stats/redist"
	util "junodb_lite/pkg/y_util"
	"time"
)

type RedistRequestContext struct {
	util.QueItemBase
	message      proto.RawMessage
	retry_cnt    uint16
	timeReceived time.Time
	reqCh        chan IRequestContext // channel for retry
	stats        *redistst.Stats
}

func NewRedistRequestContext(msg *proto.RawMessage,
	reqCh chan IRequestContext, stats *redistst.Stats) *RedistRequestContext {
	r := &RedistRequestContext{
		retry_cnt:    0,
		timeReceived: time.Now(),
		reqCh:        reqCh,
		stats:        stats,
	}
	//r.SetQueTimeout(RedistConfig.RedistRespTimeout.Duration)
	//r.message.DeepCopy(msg)
	return r
}

func (r *RedistRequestContext) OnCleanup() {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) OnExpiration() {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) Deadline() (deadline time.Time) {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) ResetDeadline() {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) SetId(id uint32) {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) GetId() uint32 {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) SetInUse(flag bool) {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) SetQueTimeout(t time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) GetQueTimeout() (t time.Duration) {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) IsInUse() bool {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) GetMessage() *proto.RawMessage {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) GetCtx() context.Context {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) Cancel() {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) Read(re io.Reader) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) WriteWithOpaque(opaque uint32, w io.Writer) (n int, err error) {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) Reply(resp IResponseContext) {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) OnComplete() {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) GetReceiveTime() time.Time {
	//TODO implement me
	panic("implement me")
}

func (r *RedistRequestContext) SetTimeout(parent context.Context, duration time.Duration) {
	//TODO implement me
	panic("implement me")
}
