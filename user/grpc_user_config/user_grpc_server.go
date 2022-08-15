package grpc_user_config

import (
	"crypto/tls"

	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/reflection"

	"github.com/shaineminkyaw/microservice/pb"
	"github.com/shaineminkyaw/microservice/user/config"
	"github.com/shaineminkyaw/microservice/user/ds"
)

type Server struct {
	pb.UnimplementedUserServiceServer
	Database ds.DataSource
}

func NewUserGRPC() (*Server, error) {
	db := ds.AuthConnectToDB()
	return &Server{
		Database: *db,
	}, nil
}

//TLS
func LoadTLSCredentials() (credentials.TransportCredentials, error) {
	//Load server certificate and private key
	serverCert, err := tls.LoadX509KeyPair("./cert/server-cert.pem", "./cert/server-key.pem")
	if err != nil {
		return nil, err
	}

	//Create the credentials and return it
	config := &tls.Config{
		Certificates: []tls.Certificate{serverCert},
		ClientAuth:   tls.NoClientCert,
	}

	return credentials.NewTLS(config), nil
}

func RunGrpcServer() {
	sourceServer, err := NewUserGRPC()
	if err != nil {
		log.Fatalf("error on source server : %v", err)
	}

	tlsCredentials, err := LoadTLSCredentials()
	if err != nil {
		log.Fatalf("cannot load tls credentials %v", err)
	}

	grpcServer := grpc.NewServer(
		grpc.Creds(tlsCredentials),
	)

	pb.RegisterUserServiceServer(grpcServer, sourceServer)
	reflection.Register(grpcServer)

	listener, err := net.Listen("tcp", config.GRPC.AuthGrpc)
	if err != nil {
		log.Fatalf("error on not working listener : %v", err)
	}

	log.Printf("starting GRPC server : %v ", config.GRPC.AuthGrpc)

	err = grpcServer.Serve(listener)
	if err != nil {
		log.Fatalf("error on not working grpc server : %v", err)
	}
}

// func RunGatewayServer() {
// 	sourceServer, err := NewUserGRPC()
// 	if err != nil {
// 		log.Fatalf("error on source server : %v", err)
// 	}

// 	grpcMux := runtime.NewServeMux()
// 	ctx, cancel := context.WithCancel(context.Background())
// 	defer cancel()
// 	err = pb.RegisterUserServiceHandlerServer(ctx, grpcMux, sourceServer)
// 	if err != nil {
// 		log.Fatalf("cannot register handler server %v", err)
// 	}

// 	mux := http.NewServeMux()
// 	mux.Handle("/", grpcMux)

// 	listener, err := net.Listen("tcp", config.App.AppAddress)
// 	if err != nil {
// 		log.Fatalf("not working listener %v", err)
// 	}

// 	log.Printf("started HTTP gateway server %v ", config.App.AppAddress)
// 	err = http.Serve(listener, mux)
// 	if err != nil {
// 		log.Fatalf("error on http gateway server %v", err)
// 	}
// }
