package service

import (
	"errors"
	"fmt"
	"gogrpcpubsub/pb"
	"google.golang.org/grpc/metadata"
	"io"
	"log"
	"sync"
)

type broadcastPayload struct {

}

type Server struct {
	pb.UnimplementedMsgServiceServer
	mu     sync.RWMutex
	client map[string]pb.MsgService_SendToUserServer
}

func NewMsgServer() *Server {
	return &Server{
		mu:     sync.RWMutex{},
		client: make(map[string]pb.MsgService_SendToUserServer),
	}
}

func (server *Server) addClient(id string, stream pb.MsgService_SendToUserServer) {
	server.mu.Lock()
	defer server.mu.Unlock()
	if _, ok := server.client[id]; ok {
		// no action
	} else {
		server.client[id] = stream
	}
}

func (server *Server) removeClient(id string) {
	server.mu.Lock()
	defer server.mu.Unlock()
	delete(server.client, id)
}

func (server *Server) getClients() map[string]pb.MsgService_SendToUserServer {
	server.mu.RLock()
	defer server.mu.RUnlock()
	return server.client
}

func (server *Server) getClient(id string) (pb.MsgService_SendToUserServer, error) {
	server.mu.RLock()
	defer server.mu.RUnlock()
	if stream, ok := server.client[id]; ok {
		return stream, nil
	}
	return nil, errors.New("not found client")
}

func (server *Server) SendToUser(stream pb.MsgService_SendToUserServer) error {
	md, ok := metadata.FromIncomingContext(stream.Context())
	if !ok {
		return errors.New("incoming context error")
	}

	id := md.Get("id")
	if len(id) == 1 {
		server.addClient(id[0], stream)
		defer func() {
			log.Println("before clients: ", server.client)
			server.removeClient(id[0])
			log.Println("after clients: ", server.client)
		}()
	}
	log.Print(server.getClients())
	//p, _ := peer.FromContext(stream.Context())
	//log.Print(p)

	to := md.Get("to")
	from := md.Get("from")
	if len(to) == 1 {

	}

	for {
		req, err := stream.Recv()
		if err == io.EOF {
			log.Print("no more data")
			break
		}
		if err != nil {
			return err
		}

		msg := req.GetMsg()
		log.Printf("server recv id: %s", id)
		log.Printf("id: %s, to: %s, from: %s, content: %s", msg.GetId(), msg.GetTo(), msg.GetFrom(), msg.GetContent())

		// send
		if len(to) == 1 && len(from) == 1 {
			st, err := server.getClient(to[0])
			if err != nil {
				return err
			}
			res := &pb.SendUserResponse{
				Id:      "",
				From:    msg.GetFrom(),
				Content: fmt.Sprintf("%s, I'm %s", msg.GetContent(), msg.GetFrom()),
			}
			err = st.Send(res)
			if err != nil {
				return err
			}
			break
		} else if len(to) == 0 && len(from) == 1 { // send all
			res := &pb.SendUserResponse{
				From: msg.GetFrom(),
				Content: fmt.Sprintf("%s, I'm %s", msg.GetContent(), msg.GetFrom()),
			}
			for _, st := range server.getClients() {
				err = st.Send(res)
				if err != nil {
					return err
				}
			}
			break
		}
	}
	return nil
}

func (server *Server) Subscribe(req *pb.SubscribeRequest, stream pb.SubscribeService_SubscribeServer) error {

}