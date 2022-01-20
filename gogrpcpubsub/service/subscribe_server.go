package service

import (
	"context"
	"github.com/google/uuid"
	"gogrpcpubsub/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"log"
)

type SubsServer struct {
	pb.UnimplementedSubscribeServiceServer
	broadcast chan *BroadcastPayload
}

type BroadcastPayload struct {
	ID    string
	Conn  chan *pb.SubscribeResponse
	Event string
}

func NewSubsServer() *SubsServer {
	return &SubsServer{
		broadcast: make(chan *BroadcastPayload),
	}
}

func (subsServer *SubsServer) Subscribe(req *pb.SubscribeRequest, stream pb.SubscribeService_SubscribeServer) error {
	log.Print("in Subscribe server call")
	id := uuid.NewString()
	conn := make(chan *pb.SubscribeResponse)

	go func() {
		//reqID := req.GetId()
		//res := &pb.SubscribeResponse{
		//	Id: fmt.Sprintf("your request id: %s", reqID),
		//}
		//conn <- res
		subsServer.broadcast <- &BroadcastPayload{
			ID:    id,
			Conn:  conn,
			Event: "remove connection",
		}
	}()

	go func() {
		<-subsServer.broadcast
		for v := range subsServer.broadcast {
			log.Print(v)
			switch v.Event {
			case "add connection":
				log.Print("add connection in goroutine")
			case "remove connection":
				log.Print("remove connection in goroutine")
			case "receive response":
				log.Print("receive response in goroutine")
			}
		}
	}()

	for {
		select {
		case <-stream.Context().Done():
			switch stream.Context().Err() {
			case context.Canceled:
				log.Print("client canceled")
			case context.DeadlineExceeded:
				log.Print("client deadline exceeded")
			}
			subsServer.broadcast <- &BroadcastPayload{
				ID:    id,
				Event: "receive response",
			}
		case response := <-conn:
			if st, ok := status.FromError(stream.Send(response)); ok {
				log.Print("response in if")
				switch st.Code() {
				case codes.OK:
					log.Print("codes ok")
					return nil
				case codes.Unavailable, codes.Canceled, codes.DeadlineExceeded:
					log.Print("codes unavailable, cancel, deadline")
					return nil
				default:
					log.Print("default case")
					return nil
				}
			}
		default:
			log.Print("default in for select")
			return nil
		}
	}
}
