package main

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"io/ioutil"
	"log"

	pb "github.com/shaineminkyaw/microservice/pb"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
)

func main() {

	// not tls config
	var opts []grpc.DialOption

	// opts = append(opts, grpc.WithInsecure())
	// conn, err := grpc.Dial("localhost:9099", opts...)
	// if err != nil {
	// 	log.Fatal("error on connection :", err)
	// }

	//include tls config

	tlsConfig, err := LoadTlsCredentials()
	if err != nil {
		log.Fatalf("error on tls config %v", err)
	}

	opts = append(opts, grpc.WithTransportCredentials(tlsConfig))
	conn, err := grpc.Dial("localhost:9099", opts...)
	if err != nil {
		log.Fatal("error on connection :", err)
	}

	defer conn.Close()
	clientRegister := pb.NewUserServiceClient(conn)

	code := &pb.RequestVerifyCode{
		Email: "smk9@gmail.com",
	}
	getCode, err := clientRegister.GetVerifyCode(context.Background(), code)
	if err != nil {
		log.Fatalf("error on generate code %v ", err)
	}

	user := &pb.UserRequest{
		Email:      "smk9@gmail.com",
		Password:   "123456",
		VerifyCode: getCode.Code,
		NationId:   "12/LMN(N)154655",
		GenderType: 2,
		City:       "taunggyi",
	}
	userResp, err := clientRegister.UserRegister(context.Background(), user)
	if err != nil {
		log.Fatalf("error on user register %v ", err)
	}

	log.Printf("User Data : %v", userResp)

}

func LoadTlsCredentials() (credentials.TransportCredentials, error) {
	//Load certificate of the CA who signed server certificate

	pemServerCA, err := ioutil.ReadFile("./cert/ca-cert.pem")
	if err != nil {
		return nil, err
	}

	certPool := x509.NewCertPool()
	if !certPool.AppendCertsFromPEM(pemServerCA) {
		return nil, err
	}

	//Create the credentials and return it
	config := &tls.Config{
		RootCAs: certPool,
	}

	return credentials.NewTLS(config), nil

}
