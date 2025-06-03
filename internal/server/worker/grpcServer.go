package worker

import (
	"context"
	"fmt"
	"net"

	"github.com/Alexandrfield/Smaug/internal/common"
	"github.com/Alexandrfield/Smaug/internal/protobufproto"
	"google.golang.org/grpc"
)

type GRPCServerPart struct {
	logicServer *ServerSafe
	protobufproto.GRPServiceServer
}

func (grpcServ *GRPCServerPart) Registration(ctx context.Context, cred *protobufproto.Cred) (*protobufproto.Err,
	error) {
	login := cred.GetLogin()
	password := cred.GetPassword()
	err := grpcServ.logicServer.Registration(login, password)
	return nil, err
}

func (grpcServ *GRPCServerPart) Login(ctx context.Context, cred *protobufproto.Cred) (*protobufproto.Err,
	error) {
	login := cred.GetLogin()
	password := cred.GetPassword()
	err := grpcServ.logicServer.Login(login, password)
	return nil, err
}

func (grpcServ *GRPCServerPart) SaveData(ctx context.Context, note *protobufproto.Note) (*protobufproto.Err,
	error) {
	login := note.GetLogin()
	data := note.GetData()
	err := grpcServ.logicServer.AddData(login, string(data))
	return nil, err
}
func (grpcServ *GRPCServerPart) GetData(ctx context.Context, note *protobufproto.Note) (*protobufproto.SavedNotes,
	error) {
	var notes protobufproto.SavedNotes
	login := note.GetLogin()
	data := note.GetData()
	existsData, err := grpcServ.logicServer.GetData(login, string(data))
	var v []*protobufproto.Note
	for _, val := range existsData {
		var note protobufproto.Note
		note.SetLogin(login)
		note.SetData([]byte(val))
		v = append(v, &note)
	}
	notes.SetNotes(v)
	return &notes, err
}

func (grpcServ *GRPCServerPart) GetAllData(ctx context.Context, note *protobufproto.Note) (*protobufproto.SavedNotes,
	error) {
	var notes protobufproto.SavedNotes
	login := note.GetLogin()
	data := note.GetData()
	existsData, err := grpcServ.logicServer.GetAllData(login, string(data))
	var v []*protobufproto.Note
	for _, val := range existsData {
		var note protobufproto.Note
		note.SetLogin(login)
		note.SetData([]byte(val))
		v = append(v, &note)
	}
	notes.SetNotes(v)
	return &notes, err
}

func StartGRPCServer(logger common.Logger, servLogic *ServerSafe, port string) {
	grpcWorker := GRPCServerPart{logicServer: servLogic}
	listener, err := net.Listen("tcp", fmt.Sprintf(":%s", port))
	if err != nil {
		logger.Fatalf("Error start tcp server for grpc: %v", err)
	}

	grpcServer := grpc.NewServer()
	protobufproto.RegisterGRPServiceServer(grpcServer, &grpcWorker)

	logger.Infof("start grpc server.")
	if err := grpcServer.Serve(listener); err != nil {
		logger.Fatalf("Error start tcp server for grpc: %v", err)
	}
}
