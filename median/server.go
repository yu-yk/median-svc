package median

import (
	"container/heap"
	"context"

	"github.com/yu-yk/median-svc/lib"
	"github.com/yu-yk/median-svc/proto"
	"google.golang.org/protobuf/types/known/timestamppb"
)

type server struct {
	proto.UnimplementedMedianServer
	leftHeap  lib.MaxHeap
	rightHeap lib.MinHeap
	status    proto.Status
}

func NewServer() *server {
	return &server{
		leftHeap:  lib.MaxHeap{},
		rightHeap: lib.MinHeap{},
		status:    proto.Status{},
	}
}

func (s *server) PushNumber(ctx context.Context, req *proto.PushNumberRequest) (*proto.PushNumberResponse, error) {
	// 1. check the num is smaller than left's top or not
	// 2. if yes, push then num to left, else push to right
	// 3. if the length difference of left and right is >= 2,
	//    rebalance by poping one side and push that to another side
	x := req.GetNumber()
	if s.leftHeap.Len() == 0 || x < int32(s.leftHeap[0]) {
		heap.Push(&s.leftHeap, x)
	} else {
		heap.Push(&s.rightHeap, x)
	}

	if s.leftHeap.Len()-s.rightHeap.Len() >= 2 {
		y := heap.Pop(&s.leftHeap)
		heap.Push(&s.rightHeap, y)
	} else if s.rightHeap.Len()-s.leftHeap.Len() >= 2 {
		y := heap.Pop(&s.rightHeap)
		heap.Push(&s.leftHeap, y)
	}

	// calculate median and set the status
	if s.leftHeap.Len() > s.rightHeap.Len() {
		s.status.Median = float64(s.leftHeap[0])
	} else if s.rightHeap.Len() > s.leftHeap.Len() {
		s.status.Median = float64(s.rightHeap[0])
	} else if s.leftHeap.Len() != 0 && s.rightHeap.Len() != 0 {
		s.status.Median = float64((s.leftHeap[0] + s.rightHeap[0]) / 2)
	}
	s.status.Size = int32(s.leftHeap.Len() + s.rightHeap.Len())
	s.status.LastUpdated = timestamppb.Now()

	return &proto.PushNumberResponse{
		Status: &s.status,
	}, nil
}

func (s *server) GetMedian(ctx context.Context, req *proto.GetMedianRequest) (*proto.GetMedianResponse, error) {
	return &proto.GetMedianResponse{Status: &s.status}, nil
}
