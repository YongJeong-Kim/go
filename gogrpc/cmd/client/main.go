package main

import (
	"context"
	"flag"
	"gogrpc/pb"
	"gogrpc/sample"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"time"
)

func createPerson(personClient pb.PersonServiceClient) {
	person := sample.NewPerson()
	req := &pb.CreatePersonRequest{
		Person: person,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	res, err := personClient.CreatePerson(ctx, req)
	if err != nil {
		st, ok := status.FromError(err)
		if ok && st.Code() == codes.AlreadyExists {
			log.Print("already exists person")
		} else {
			log.Fatalf("cannot create person %s", err)
		}
		return
	}
	log.Print("create person:", res.PersonId)
}

func searchPerson(personClient pb.PersonServiceClient, filter *pb.Filter) {
	log.Print("search person:", filter)

	req := &pb.SearchPersonRequest{
		Filter: filter,
	}
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	stream, err := personClient.SearchPerson(ctx, req)
	if err != nil {
		log.Fatalf("cannot search person", err)
	}

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			return
		}
		if err != nil {
			log.Fatal("cannot receive response:", err)
		}

		p := res.GetPerson()
		log.Print("found:", p)
	}
}

/*func main() {
	serverAddress := flag.String("address", "", "server address")
	flag.Parse()
	log.Printf("dial server: %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server", err)
	}

	personClient := pb.NewPersonServiceClient(conn)
	// create
	//createPerson(personClient)

	// search
	//for i := 0; i < 10; i++ {
	//	createPerson(personClient)
	//}
	//filter := &pb.Filter{
	//	Brand: "LACOSTE",
	//}
	//searchPerson(personClient, filter)

	stream, err := personClient.UploadImage(context.Background())
	if err != nil {
		return
	}

	//file, err := os.Open("C:\\Users\\admin\\Desktop\\ggo.png")
	file, err := os.Open("C:\\Users\\admin\\Downloads\\Docker Desktop Installer.exe")
	if err != nil {
		log.Fatal("cannot open image file: ", err)
	}
	defer func(file *os.File) {
		err := file.Close()
		if err != nil {
			log.Fatal("close file failed")
		}
	}(file)

	err = stream.Send(&pb.UploadImageRequest{
		ChunkData: nil,
	})
	if err != nil {
		log.Fatal("cannot send image data")
	}

	reader := bufio.NewReader(file)
	buffer := make([]byte, 1024)

	for {
		n, err := reader.Read(buffer)
		if err == io.EOF {
			log.Print("read image buffer end")
			break
		}
		if err != nil {
			log.Fatal("cannot read chunk to buffer: ", err)
		}

		req := &pb.UploadImageRequest{
			ChunkData: buffer[:n],
		}

		err = stream.Send(req)
		if err != nil {
			log.Fatal("cannot send chunk to server: ", err, stream.RecvMsg(nil))
		}
	}

	res, err := stream.CloseAndRecv()
	if err != nil {
		log.Fatal("cannot receive response: ", err)
	}

	log.Printf("image uploaded with size: %d", res.GetSize())
}*/

func main() {
	serverAddress := flag.String("address", "", "server address")
	flag.Parse()
	log.Printf("dial server: %s", *serverAddress)

	conn, err := grpc.Dial(*serverAddress, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal("cannot dial server", err)
	}

	shirtClient := pb.NewShirtServiceClient(conn)
	stream, err := shirtClient.Broadcast(context.Background())
	if err != nil {
		return
	}
	req := &pb.ShirtRequest{
		Shirt: &pb.Shirt{
			Brand: "Subscribe",
			//Brand: "Publish",
			//Brand: "Only",
		},
	}

	go func() {
		err = stream.Send(req)
	}()

	for {
		res, err := stream.Recv()
		if err == io.EOF {
			log.Println("end")
			break
		}
		if err != nil {
			log.Fatal("fatal err ", err)
		}
		log.Println(res.GetShirt())
	}
}
