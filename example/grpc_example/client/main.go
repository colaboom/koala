package main

const (
	address     = "localhost:8080"
	defaultName = "colaboom"
)

/*func main() {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal("did not connect, err :%v", err)
	}
	defer conn.Close()
	c := pb.NewHelloServiceClient(conn)

	name := defaultName
	if len(os.Args) > 1 {
		name = os.Args[1]
	}
	for {
		resp, err := c.SayHello(context.Background(), &pb.HelloRequest{Name: name})
		if err != nil {
			log.Printf("could not greet :%v", err)
			continue
		}
		log.Printf("Greeting: %s", resp.Reply)
		//time.Sleep(time.Millisecond * 10)
	}
}*/
