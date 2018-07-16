package main;

import (
    pb "cloudtogo.local/biz/test-grpc/grpc"
    "golang.org/x/net/context"
    "net"
    "log"
    "google.golang.org/grpc"
    "fmt"
    "time"
    "math/rand"
)

const (
    port = ":50051"
)

type server struct{}

func (this *server) SayHello(ctx context.Context, in *pb.HelloRequest) (*pb.HelloReply, error) {
    log.Print(">>>> ", in.Name, "; ", in.Sex)
    return &pb.HelloReply{Message: "Hello " + in.Name}, nil
}

func (this *server) SingleStream(request *pb.HelloRequest, stream pb.Greeter_SingleStreamServer) error {
    index := 0
    log.Print(">>>> ", request.Name)
    ctx, cancel := context.WithCancel(context.Background())
    defer cancel()
    for {
        index++
        select {
        case <-ctx.Done():
            if ctx.Err() == context.Canceled {
                log.Print("ERROR Canceled: " + ctx.Err().Error())
                //
            } else if ctx.Err() == context.DeadlineExceeded {
                log.Print("ERROR DeadlineExceeded: " + ctx.Err().Error())
                //
            } else {
                log.Print("ERROR: " + ctx.Err().Error())
            }
        default:
            message := fmt.Sprintf("index_%d_%s", index, request.Name)
            log.Print("send: ", message)
            if ctx.Err() != nil {
                log.Print("ERROR >> : " + ctx.Err().Error())
                break
            }
            if err := stream.Send(&pb.HelloReply{Message:message}); err != nil {
                // eg: rpc error: code = Unavailable desc = transport is closing
                log.Print("ERROR::: " + err.Error())
                return err
            }
            time.Sleep(time.Duration(rand.Intn(10) * 100) * time.Millisecond)
        }
    }
    //for {
    //    index++
    //    message := fmt.Sprintf("index_%d", index)
    //    log.Print("send: ", message)
    //    //ctx.Done()
    //    if err := stream.Send(&pb.HelloReply{Message:message}); err != nil {
    //        log.Print("ERROR: " + err.Error())
    //        return err
    //    }
    //    time.Sleep(time.Duration(rand.Intn(10) * 100) * time.Millisecond)
    //}
    return nil
}
func (this *server)  Chat(stream pb.Greeter_ChatServer) error {
    // TODO
    return nil
}

func main() {
    lis, err := net.Listen("tcp", port)
    if err != nil {
        log.Fatal("failed to listen: %v", err)
    }
    s := grpc.NewServer()
    pb.RegisterGreeterServer(s, &server{})
    log.Print("start ...")
    s.Serve(lis)
}