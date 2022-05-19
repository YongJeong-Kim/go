package service

import (
	"context"
	"gogrpc/pb"
	"io"
	"log"
	"time"
)

type ShirtServer struct {
	pb.UnimplementedShirtServiceServer
	rr        chan chan int
	publishCh chan interface{}
	subCh     chan chan interface{}
	unsubCh   chan chan interface{}
	stopCh    chan struct{}
	mySubCh   chan pb.ShirtService_BroadcastServer
}

func NewShirtServer() *ShirtServer {
	return &ShirtServer{
		rr:        make(chan chan int),
		publishCh: make(chan interface{}),
		subCh:     make(chan chan interface{}),
		unsubCh:   make(chan chan interface{}),
		stopCh:    make(chan struct{}),
		mySubCh:   make(chan pb.ShirtService_BroadcastServer),
	}
}

func (s *ShirtServer) Publish(msgCh interface{}) {
	s.publishCh <- msgCh
}

func (s *ShirtServer) Subscribe() chan interface{} {
	subCh := make(chan interface{})
	s.subCh <- subCh
	return subCh
}

func (s *ShirtServer) Unsubscribe(msgCh chan interface{}) {
	s.unsubCh <- msgCh
	//close(msgCh)
}

func (s *ShirtServer) Broadcast(stream pb.ShirtService_BroadcastServer) error {
	for {
		res, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			break
		}
		switch res.GetShirt().GetBrand() {
		case "mySub":

		case "Subscribe":
			if err := stream.Send(&pb.ShirtResponse{
				Shirt: &pb.Shirt{
					Brand: "subs",
				},
			}); err != nil {
				log.Println("subs event failed. ", err)
			}
			go func(st pb.ShirtService_BroadcastServer) {
				select {
				case <-st.Context().Done():
					err := st.Context().Err()
					switch err {
					case context.Canceled:
						log.Println("sub canceled")
						//s.errCh <- status.Errorf(codes.Canceled, "canceledddd")
						//s.Unsubscribe(sCh)
					default:
						log.Println("done")
					}
				}
			}(stream)

			go func() {
				log.Println("in case subs")
				subCh := s.Subscribe()

				/*go func(st pb.ShirtService_BroadcastServer, sCh chan interface{}) {
					select {
					case <-st.Context().Done():
						err := st.Context().Err()
						switch err {
						case context.Canceled:
							log.Println("sub canceled")
							//s.errCh <- status.Errorf(codes.Canceled, "canceledddd")
							s.Unsubscribe(sCh)
						default:
							log.Println("done")
						}
					}
				}(stream, subCh)*/

				st := <-subCh
				log.Println("in case subs before send")
				if st != nil {
					//for {
					stream.Send(&pb.ShirtResponse{
						Shirt: &pb.Shirt{
							Brand: st.(string),
						},
					},
					)
					time.Sleep(500 * time.Millisecond)
					//}
				}
			}()
		case "Publish":
			go func(st pb.ShirtService_BroadcastServer) {
				select {
				case <-st.Context().Done():
					err := st.Context().Err()
					switch err {
					case context.Canceled:
						log.Println("pub canceled")
						//s.errCh <- status.Errorf(codes.Canceled, "canceledddd")
					default:
						log.Println("done")
					}
				}
			}(stream)
			s.Publish(res.GetShirt().GetBrand())
		case "Unsubscribe":

		case "Only":
			s.rr <- make(chan int)
		}
	}
	return nil
}

//func (s *ShirtServer) b(stream pb.ShirtService_BroadcastServer) error {
func (s *ShirtServer) ChannelReceiver() error {
	subs := map[chan interface{}]struct{}{}
	//sss := map[pb.ShirtService_BroadcastServer]struct{}{}
	for {
		select {
		case sCh := <-s.subCh:
			subs[sCh] = struct{}{}
			log.Println(subs)
		case unsCh := <-s.unsubCh:
			delete(subs, unsCh)
			//close(unsCh)
			log.Println("unsub: ", subs, unsCh)
		case <-s.stopCh:
			for subCh := range subs {
				close(subCh)
			}
		case msg := <-s.publishCh:
			for msgCh := range subs {
				/*	log.Println("sec11: ", msg)
					log.Println("sec22: ", msgCh)*/

				select {
				case msgCh <- msg:
					log.Println("sec: ", msgCh, msg)
					/*(<-msgCh).(pb.ShirtService_BroadcastServer).Send(&pb.ShirtResponse{
						Shirt: &pb.Shirt{
							Brand: msg.(string),
						},
					},
					)*/
					log.Println(subs)
				default:
					log.Println("dfeault")
				}
			}
			/*case <-stream.Context().Done():
			err := stream.Context().Err()
			switch err {
			case context.Canceled:
				log.Println("canceled")
				//s.errCh <- status.Errorf(codes.Canceled, "canceledddd")
				return status.Errorf(codes.Canceled, "canceledddd")
			default:
				log.Println("done")
				return nil
			}*/
			/*case r := <-s.rr:
			_ = stream.Send(&pb.ShirtResponse{
				Shirt: &pb.Shirt{
					Brand: "asdasd",
				},
			})
			log.Println(r)*/
		}
	}
}
