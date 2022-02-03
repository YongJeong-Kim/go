package service

import (
	"context"
	"errors"
	"fmt"
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
	Request  *pb.SubscribeRequest
	Response *pb.SubscribeResponse
	Event    string
	stream   pb.SubscribeService_SubscribeBidiServer
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
		Response: nil,
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
	res := make(chan *pb.SubscribeResponse, 1)
	payload := make(chan *BidiPayload, 1)

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

	//go subsServer.SelectEvent(payload, res, errCh)
	go func() {
		//log.Print("subserver bidi payload")
		//<-subsServer.bidiPayload
		p := <-payload
		log.Print(p)
		//for v := range subsServer.bidiPayload {
		switch p.Event {
		case "subscribe":
			log.Print("subscribe")

			id := p.Request.GetId()
			subsServer.addClient(id, stream)
			log.Printf("subscribe %v", subsServer.client)
			res <- &pb.SubscribeResponse{
				Id:      p.Request.GetId(),
				From:    p.Request.GetFrom(),
				Content: p.Request.GetContent(),
			}
		case "unsubscribe":
			log.Print("unsubscribe")
			id := p.Request.GetId()
			log.Printf("unsubscribe %v", subsServer.client)
			client, err := subsServer.getClient(id)
			if err != nil {
				log.Print("get client error. ", err)
			}
			err = client.Send(&pb.SubscribeResponse{
				Id: fmt.Sprintf("%s was unsubscribed", id),
			})
			if err != nil {
				log.Print("client send error from server. ", err)
			}
			subsServer.removeClient(id)
			log.Print("current clients: ", subsServer.client)
			//if err != nil {
			//	res <- &pb.SubscribeResponse{
			//		Id: "err is not nil",
			//	}
			//}
			res <- &pb.SubscribeResponse{
				Id: fmt.Sprintf("%s was unsubscribed", id),
			}
			//id := req.GetId()
			log.Printf("all client: %v", subsServer.client)
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

		log.Print("set bidi payload event:", req.GetEvent())

		p := NewBidiPayload()
		p.Request = &pb.SubscribeRequest{
			Id:      req.GetId(),
			From:    req.GetFrom(),
			Content: req.GetContent(),
			Event:   req.GetEvent(),
			To:      req.GetTo(),
		}
		p.Event = req.GetEvent()
		payload <- p
	}
}

func (subsServer *SubsServer) sendMessage(stream pb.SubscribeService_SubscribeServer, errCh chan error) {
	//for {
	//	stream.Send()
	//}
}

//func (subsServer *SubsServer) SelectEvent(payload chan *BidiPayload, response chan *pb.SubscribeResponse, errCh chan error) {
//	p := <-payload
//	log.Print(p)
//
//	switch p.Event {
//	case "subscribe":
//		log.Print("subscribe")
//		id := p.Request.GetId()
//		subsServer.addClient(id, p.stream)
//		log.Printf("subscribe %v", subsServer.client)
//		response <- &pb.SubscribeResponse{
//			Id:      p.Request.GetId(),
//			From:    p.Request.GetFrom(),
//			Content: p.Request.GetContent(),
//		}
//	case "unsubscribe":
//		log.Print("unsubscribe")
//		id := p.Request.GetId()
//		subsServer.removeClient(id)
//		log.Printf("unsubscribe %v", subsServer.client)
//		client, err := subsServer.getClient(id)
//		if err != nil {
//			response <- &pb.SubscribeResponse{
//				Id: "err is not nil",
//			}
//		}
//		response <- &pb.SubscribeResponse{
//			Id: id,
//		}
//		//id := req.GetId()
//		//log.Printf("remove client: %v", id)
//		//subsServer.removeClient(id)
//	case "sendToUser":
//		log.Print("sendToUser")
//		//log.Printf("send to user, to: %v, from: %v", req.GetTo(), req.GetFrom())
//		//to := req.GetTo()
//		//client, err := subsServer.getClient(to)
//		//if err != nil {
//		//	log.Print("get client failed.", err)
//		//	break
//		//}
//		//err = client.Send(&pb.SubscribeResponse{
//		//	From:    req.GetFrom(),
//		//	Content: req.GetContent(),
//		//})
//		//if err != nil {
//		//	log.Print("server send err", err)
//		//	break
//		//}
//	case "send to all":
//		log.Print("send to all")
//	default:
//		log.Print("default event, no action")
//		//log.Print("default event, no action", req)
//		errCh <- errors.New("default event, no action")
//	}
//}
