package cli

import (
	"io"
	proto "junodb_lite/pkg/a_proto"
)

type RecordInfo struct {
	version      uint32
	creationTime uint32
	timeToLive   uint32
	originatorId proto.RequestId
}

func (r *RecordInfo) GetVersion() uint32 {
	//TODO implement me
	panic("implement me")
}

func (r *RecordInfo) GetCreationTime() uint32 {
	//TODO implement me
	panic("implement me")
}

func (r *RecordInfo) GetTimeToLive() uint32 {
	//TODO implement me
	panic("implement me")
}

func (r *RecordInfo) PrettyPrint(w io.Writer) {
	//TODO implement me
	panic("implement me")
}

func (r *RecordInfo) SetFromOpMsg(m *proto.OperationalMessage) {

}

func (r *RecordInfo) SetRequestWithUpdateCond(request *proto.OperationalMessage) {

}
