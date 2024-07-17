package grpc

import (
	context "context"
	"errors"
	"fmt"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/logging"
	"github.com/grpc-ecosystem/go-grpc-middleware/v2/interceptors/recovery"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/known/emptypb"
	"homework9/internal/app"
	"homework9/internal/entities/ads"
	"homework9/internal/entities/user"
	"log"
	"os"
)

type AdService struct {
	UnimplementedAdServiceServer
	app app.App
}

func NewAdService(app app.App) *AdService {
	return &AdService{
		UnimplementedAdServiceServer{},
		app,
	}
}

func InterceptorLogger(l *log.Logger) logging.Logger {
	return logging.LoggerFunc(func(_ context.Context, lvl logging.Level, msg string, fields ...any) {
		switch lvl {
		case logging.LevelDebug:
			msg = fmt.Sprintf("DEBUG :%v", msg)
		case logging.LevelInfo:
			msg = fmt.Sprintf("INFO :%v", msg)
		case logging.LevelWarn:
			msg = fmt.Sprintf("WARN :%v", msg)
		case logging.LevelError:
			msg = fmt.Sprintf("ERROR :%v", msg)
		default:
			panic(fmt.Sprintf("unknown level %v", lvl))
		}
		l.Println(append([]any{"msg", msg}, fields...))
	})
}

func NewServerWithLogger() *grpc.Server {
	logger := log.New(os.Stderr, "", log.Ldate|log.Ltime|log.Lshortfile)

	panicLog := func(p any) (err error) {
		return status.Errorf(codes.Unknown, "panic triggered: %v", p)
	}

	opt := logging.WithLogOnEvents(logging.StartCall, logging.FinishCall)

	optPanic := recovery.WithRecoveryHandler(panicLog)

	srv := grpc.NewServer(
		grpc.ChainUnaryInterceptor(
			logging.UnaryServerInterceptor(InterceptorLogger(logger), opt),
			recovery.UnaryServerInterceptor(optPanic)),
	)
	return srv
}

func (a AdService) CreateAd(ctx context.Context, request *CreateAdRequest) (*AdResponse, error) {
	ad, err := a.app.CreateAd(request.Title, request.Text, request.UserId)

	if errors.As(err, &app.NoSuchUserError{Err: app.ErrorNoSuchUser}) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied. You are not owner")
	}

	if errors.As(err, &app.ValidationApError{Err: app.ErrorValidation}) {
		return nil, status.Error(codes.InvalidArgument, "Invalid Argument")
	}

	if err != nil {
		return nil, status.Error(codes.Canceled, "Server error")
	}

	return toAdResponse(ad), nil
}

func (a AdService) ChangeAdStatus(ctx context.Context, request *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := a.app.ChangeAdStatus(request.AdId, request.UserId, request.Published)

	if errors.As(err, &app.PermissionDeniedError{Err: app.ErrorNoRights}) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied. You are not owner")
	}

	if errors.As(err, &app.ValidationApError{Err: app.ErrorValidation}) {
		return nil, status.Error(codes.InvalidArgument, "Invalid Argument")
	}

	if err != nil {
		return nil, status.Error(codes.Canceled, "Server error")
	}

	return toAdResponse(ad), nil
}

func (a AdService) UpdateAd(ctx context.Context, request *UpdateAdRequest) (*AdResponse, error) {
	ad, err := a.app.UpdateAd(request.AdId, request.UserId, request.Title, request.Text)
	if errors.As(err, &app.PermissionDeniedError{Err: app.ErrorNoRights}) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied. You are not owner")
	}

	if errors.As(err, &app.ValidationApError{Err: app.ErrorValidation}) {
		return nil, status.Error(codes.InvalidArgument, "Invalid Argument")
	}

	if err != nil {
		return nil, status.Error(codes.Canceled, "Server error")
	}
	return toAdResponse(ad), nil
}

func (a AdService) ListAds(ctx context.Context, empty *emptypb.Empty) (*ListAdResponse, error) {

	answer, err := a.app.GetAdsList(ads.NewFilter())

	if err != nil {
		return nil, status.Error(codes.Canceled, "Server error")
	}
	return toAdList(answer), nil
}

func (a AdService) CreateUser(ctx context.Context, request *CreateUserRequest) (*UserResponse, error) {
	u, err := a.app.CreateUser(request.Name, request.Email)
	if err != nil {
		return nil, status.Error(codes.Canceled, "Server error")
	}
	return toUserResponse(u), nil
}

func (a AdService) GetUser(ctx context.Context, request *GetUserRequest) (*UserResponse, error) {
	answer, err := a.app.GetUser(request.Id)

	if errors.As(err, &app.NoSuchUserError{Err: app.ErrorNoSuchUser}) {
		return nil, status.Error(codes.InvalidArgument, "No such user")
	}

	if err != nil {
		return nil, status.Error(codes.Canceled, "Server error")
	}

	return toUserResponse(answer), nil
}

func (a AdService) DeleteUser(ctx context.Context, request *DeleteUserRequest) (*emptypb.Empty, error) {
	err := a.app.DeleteUser(request.Id)

	if errors.As(err, &app.NoSuchUserError{Err: app.ErrorNoSuchUser}) {
		return nil, status.Error(codes.InvalidArgument, "No such user")
	}

	if err != nil {
		return nil, status.Error(codes.Canceled, "Server error")
	}

	return &emptypb.Empty{}, nil
}

func (a AdService) DeleteAd(ctx context.Context, request *DeleteAdRequest) (*emptypb.Empty, error) {
	err := a.app.DeleteAd(request.AdId, request.AuthorId)

	if errors.As(err, &app.NoSuchUserError{Err: app.ErrorNoSuchUser}) {
		return nil, status.Error(codes.InvalidArgument, "No such user")
	}

	if errors.As(err, &app.NoSuchAdError{Err: app.ErrorNoAd}) {
		return nil, status.Error(codes.InvalidArgument, "No such add")
	}

	if err != nil {
		return nil, status.Error(codes.Canceled, "Server error")
	}

	return &emptypb.Empty{}, nil
}

func (a AdService) GetAd(ctx context.Context, request *GetAdRequest) (*AdResponse, error) {
	answer, err := a.app.GetAd(request.AdId)

	if errors.As(err, &app.NoSuchAdError{Err: app.ErrorNoAd}) {
		return nil, status.Error(codes.InvalidArgument, "No such add")
	}

	if errors.As(err, &app.NoSuchUserError{Err: app.ErrorNoSuchUser}) {
		return nil, status.Error(codes.InvalidArgument, "No such user")
	}

	if err != nil {
		return nil, status.Error(codes.Canceled, "Server error")
	}

	return toAdResponse(answer), nil
}

func toAdResponse(ad ads.Ad) *AdResponse {
	return &AdResponse{Id: ad.ID, AuthorId: ad.AuthorID, Title: ad.Title, Text: ad.Text, Published: ad.Published}
}

func toAdList(ads []ads.Ad) *ListAdResponse {
	adsList := make([]*AdResponse, 0)
	for _, v := range ads {
		adsList = append(adsList, toAdResponse(v))
	}
	return &ListAdResponse{List: adsList}
}

func toUserResponse(u user.User) *UserResponse {
	return &UserResponse{Name: u.Nickname, Id: u.ID, Email: u.Email}
}
