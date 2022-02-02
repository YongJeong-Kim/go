package service

import (
	"context"
	"errors"
	"github.com/google/uuid"
	"gogrpcpubsub/pb"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"sync"
)

type SubsServer struct {
	pb.UnimplementedSubscribeServiceServer
	broadcast   chan *BroadcastPayload
	client      map[string]pb.SubscribeService_SubscribeBidiServer
	event       string
	bidiPayload chan *BidiPayload
	mu          sync.RWMutex
}

type BroadcastPayload struct {
	ID       string
	Response chan *pb.SubscribeResponse
	Event    string
}

type BidiPayload struct {
	Response chan *pb.SubscribeResponse
	Event    string
}

func NewSubsServer() *SubsServer {
	return &SubsServer{
		broadcast:   make(chan *BroadcastPayload),
		client:      make(map[string]pb.SubscribeService_SubscribeBidiServer),
		bidiPayload: make(chan *BidiPayload),
	}
}

func NewBidiPayload() *BidiPayload {
	return &BidiPayload{
		Response: make(chan *pb.SubscribeResponse),
		Event:    "",
	}
}

func (subsServer *SubsServer) addClient(id string, stream pb.SubscribeService_SubscribeBidiServer) {
	subsServer.mu.Lock()
	defer subsServer.mu.Unlock()
	subsServer.client[id] = stream
}

func (subsServer *SubsServer) removeClient(id string) {
	subsServer.mu.Lock()
	defer subsServer.mu.Unlock()
	delete(subsServer.client, id)
}

func (subsServer *SubsServer) getClient(id string) (pb.SubscribeService_SubscribeBidiServer, error) {
	subsServer.mu.RLock()
	defer subsServer.mu.RUnlock()
	if c, ok := subsServer.client[id]; ok {
		return c, nil
	}
	return nil, errors.New("not found client from getClient(id)")
}

func (subsServer *SubsServer) Subscribe(req *pb.SubscribeRequest, stream pb.SubscribeService_SubscribeServer) error {
	log.Print("in Subscribe server call")
	id := uuid.NewString()
	res := make(chan *pb.SubscribeResponse)

	go func() {
		//reqID := req.GetId()
		//res := &pb.SubscribeResponse{
		//	Id: fmt.Sprintf("your request id: %s", reqID),
		//}
		//conn <- res
		subsServer.broadcast <- &BroadcastPayload{
			ID:       id,
			Response: res,
			Event:    "remove connection",
		}
	}()

	go func() {
		<-subsServer.broadcast
		for v := range subsServer.broadcast {
			log.Print(v)
			switch v.Event {
			case "subscribe":
				log.Print("add connection in goroutine")
			case "unsubscribe":
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
		case response := <-res:
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

func (subsServer *SubsServer) SubscribeBidi(stream pb.SubscribeService_SubscribeBidiServer) error {
	errCh := make(chan error)
	go subsServer.receiveMessage(stream, errCh)
	go subsServer.sendMessage(stream, errCh)

	return <-errCh
}

func (subsServer *SubsServer) receiveMessage(stream pb.SubscribeService_SubscribeBidiServer, errCh chan error) {
	res := make(chan *pb.SubscribeResponse)
	payload := make(chan *BidiPayload)

	go func() {
		for {
			select {
			case <-stream.Context().Done():
				switch stream.Context().Err() {
				case context.Canceled:
					log.Print("client canceled")
					errCh <- errors.New("client canceled")
				case context.DeadlineExceeded:
					log.Print("client deadline exceeded")
					errCh <- errors.New("client deadline exceeded")
				default:
					errCh <- errors.New("select case default")
				}
			case response := <-res:
				log.Print("case response,", response)
				if st, ok := status.FromError(stream.Send(response)); ok {
					log.Print("response in if")
					switch st.Code() {
					case codes.OK:
						log.Print("codes ok")
					case codes.Unavailable, codes.Canceled, codes.DeadlineExceeded:
						log.Print("codes unavailable, cancel, deadline")
					default:
						log.Print("default case")
					}
				}
			}
		}
	}()

	go func() {
		//log.Print("subserver bidi payload")
		//<-subsServer.bidiPayload
		p := <-payload
		log.Print(p)
		//for v := range subsServer.bidiPayload {
		switch p.Event {
		case "subscribe":
			log.Print("subscribe")
			//id := req.GetId()
			//log.Printf("add client: %v", id)
			//subsServer.addClient(id, stream)
			//res <- &pb.SubscribeResponse{
			//	Id: "res res res",
			//}
		case "unsubscribe":
			log.Print("unsubscribe")
			//id := req.GetId()
			//log.Printf("remove client: %v", id)
			//subsServer.removeClient(id)
		case "sendToUser":
			log.Print("sendToUser")
			//log.Printf("send to user, to: %v, from: %v", req.GetTo(), req.GetFrom())
			//to := req.GetTo()
			//client, err := subsServer.getClient(to)
			//if err != nil {
			//	log.Print("get client failed.", err)
			//	break
			//}
			//err = client.Send(&pb.SubscribeResponse{
			//	From:    req.GetFrom(),
			//	Content: req.GetContent(),
			//})
			//if err != nil {
			//	log.Print("server send err", err)
			//	break
			//}
		case "send to all":
			log.Print("send to all")
		default:
			log.Print("default event, no action")
			//log.Print("default event, no action", req)
			errCh <- errors.New("default event, no action")
		}
	}()

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			log.Print("receive from server err is not nil", err)
			errCh <- err
		}

		//go func() {
		//event <- req.GetEvent()
		log.Print("set bidi payload event:", req.GetEvent())
		res <- &pb.SubscribeResponse{
			Id: "aa",
		}
		//subsServer.bidiPayload <- &BidiPayload{
		//	Response: res,
		//	Event:    req.GetEvent(),
		//}
		//p := NewBidiPayload()
		//p.Event = "subscribe"
		//payload <- p
		//}()

	}

}

func (subsServer *SubsServer) sendMessage(stream pb.SubscribeService_SubscribeServer, errCh chan error) {
	//for {
	//	stream.Send()
	//}
}
