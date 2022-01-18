package service

import (
	"errors"
	"gogrpcpubsub/pb"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/peer"
	"io"
	"log"
	"sync"
)

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
	//id := uuid.NewString()
	//server.addClient(id, stream)
	//defer server.removeClient(id)
	md, _ := metadata.FromIncomingContext(stream.Context())
	p, _ := peer.FromContext(stream.Context())
	log.Print(md)
	log.Print(p)

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
		id := msg.GetId()
		server.addClient(id, stream)
		log.Printf("server recv id: %s", id)
		//defer server.removeClient(id)
		//toID := msg.GetTo()

		//client, err := server.getClient(toID)
		//if err != nil {
		//	return errors.New("not found client")
		//}
		//res := &pb.SendUserResponse{
		//	Id:      id,
		//	From:    msg.GetFrom(),
		//	Content: msg.GetContent(),
		//}
		//err = client.Send(res)
		//if err != nil {
		//	return err
		//}

		for _, st := range server.getClients() {
			res := &pb.SendUserResponse{
				Id:      id,
				From:    msg.GetFrom(),
				Content: msg.GetContent(),
			}
			err := st.Send(res)
			if err != nil {
				return err
			}
		}
	}
	return nil
}
