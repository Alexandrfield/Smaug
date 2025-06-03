package worker

import (
	"context"
	"fmt"
	"log"
	"time"

	"github.com/Alexandrfield/Smaug/internal/protobufproto"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

type GRPCClient struct {
	serverAddr string
}

func (c *GRPCClient) Registration(login string, password string) error {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	grpcClient := protobufproto.NewGRPServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var cred protobufproto.Cred
	cred.SetLogin(login)
	cred.SetPassword([]byte(password))
	_, err = grpcClient.Registration(ctx, &cred)
	if err != nil {
		return fmt.Errorf("registration filed. err:%w", err)
	}
	return nil
}

func (c *GRPCClient) Login(login string, password []byte) error {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	grpcClient := protobufproto.NewGRPServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var cred protobufproto.Cred
	cred.SetLogin(login)
	cred.SetPassword(password)
	_, err = grpcClient.Login(ctx, &cred)
	if err != nil {
		return fmt.Errorf("login filed. err:%w", err)
	}
	return nil
}

func (c *GRPCClient) SaveData(login string, dat []byte) error {
	conn, err := grpc.Dial(c.serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	grpcClient := protobufproto.NewGRPServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var note protobufproto.Note
	note.SetLogin(login)
	note.SetData(dat)
	_, err = grpcClient.SaveData(ctx, &note)
	if err != nil {
		return fmt.Errorf("save data filed. err:%w", err)
	}
	return nil
}
func (c *GRPCClient) GetData(login string, dat []byte) ([][]byte, error) {
	var res [][]byte
	conn, err := grpc.Dial(c.serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	grpcClient := protobufproto.NewGRPServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var note protobufproto.Note
	note.SetLogin(login)
	note.SetData(dat)
	savedNotes, err := grpcClient.GetData(ctx, &note)
	notes := savedNotes.GetNotes()
	for _, v := range notes {
		dat := v.GetData()
		res = append(res, dat)
	}
	if err != nil {
		return res, fmt.Errorf("getData filed. err:%w", err)
	}
	return res, nil
}
func (c *GRPCClient) GetAllData(login string, dat []byte) ([][]byte, error) {
	var res [][]byte
	conn, err := grpc.Dial(c.serverAddr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()
	grpcClient := protobufproto.NewGRPServiceClient(conn)
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	var note protobufproto.Note
	note.SetLogin(login)
	note.SetData(dat)
	savedNotes, err := grpcClient.GetAllData(ctx, &note)
	notes := savedNotes.GetNotes()
	for _, v := range notes {
		dat := v.GetData()
		res = append(res, dat)
	}
	if err != nil {
		return res, fmt.Errorf("getAllData filed. err:%w", err)
	}
	return res, nil
}
