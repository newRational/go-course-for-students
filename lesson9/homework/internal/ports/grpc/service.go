package grpc

import (
	"context"
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"homework9/internal/ads"
	"homework9/internal/app"
)

type Server struct {
	app app.App
}

func NewService(a app.App) *Server {
	s := &Server{
		app: a,
	}
	return s
}

func (s *Server) CreateAd(ctx context.Context, req *CreateAdRequest) (*AdResponse, error) {
	ad, err := s.app.CreateAd(ctx, req.Title, req.Text, req.UserId)
	if errors.Is(err, app.ErrBadRequest) {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		UserId:    ad.UserID,
		Published: ad.Published,
	}, nil
}

func (s *Server) GetAd(ctx context.Context, req *GetAdRequest) (*AdResponse, error) {
	ad, err := s.app.AdByID(ctx, req.Id)
	if errors.Is(err, app.ErrBadRequest) {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		UserId:    ad.UserID,
		Published: ad.Published,
	}, nil
}

func (s *Server) ListAds(ctx context.Context, req *ListAdsRequest) (*ListAdResponse, error) {
	adverts, err := s.app.AdsByPattern(ctx, createAdPattern(req))
	if errors.Is(err, app.ErrBadRequest) {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	var list []*AdResponse
	for i := range adverts {
		list = append(list, &AdResponse{
			Id:        adverts[i].ID,
			Title:     adverts[i].Title,
			Text:      adverts[i].Text,
			UserId:    adverts[i].UserID,
			Published: adverts[i].Published,
		})
	}

	return &ListAdResponse{
		List: list,
	}, nil
}

func (s *Server) UpdateAd(ctx context.Context, req *UpdateAdRequest) (*AdResponse, error) {
	ad, err := s.app.UpdateAd(ctx, req.AdId, req.UserId, req.Title, req.Text)
	if errors.Is(err, app.ErrBadRequest) {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	} else if errors.Is(err, app.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		UserId:    ad.UserID,
		Published: ad.Published,
	}, nil
}

func (s *Server) ChangeAdStatus(ctx context.Context, req *ChangeAdStatusRequest) (*AdResponse, error) {
	ad, err := s.app.ChangeAdStatus(ctx, req.AdId, req.UserId, req.Published)
	if errors.Is(err, app.ErrBadRequest) {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	} else if errors.Is(err, app.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		UserId:    ad.UserID,
		Published: ad.Published,
	}, nil
}

func (s *Server) DeleteAd(ctx context.Context, req *DeleteAdRequest) (*AdResponse, error) {
	ad, err := s.app.DeleteAd(ctx, req.AdId, req.UserId)
	if errors.Is(err, app.ErrBadRequest) {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	} else if errors.Is(err, app.ErrForbidden) {
		return nil, status.Error(codes.PermissionDenied, "Permission denied")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &AdResponse{
		Id:        ad.ID,
		Title:     ad.Title,
		Text:      ad.Text,
		UserId:    ad.UserID,
		Published: ad.Published,
	}, nil
}

func (s *Server) CreateUser(ctx context.Context, req *CreateUserRequest) (*UserResponse, error) {
	u, err := s.app.CreateUser(ctx, req.Nickname, req.Email)
	if errors.Is(err, app.ErrBadRequest) {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &UserResponse{
		Id:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
	}, nil
}

func (s *Server) GetUser(ctx context.Context, req *GetUserRequest) (*UserResponse, error) {
	u, err := s.app.UserByID(ctx, req.Id)
	if errors.Is(err, app.ErrBadRequest) {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &UserResponse{
		Id:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
	}, nil
}

func (s *Server) UpdateUser(ctx context.Context, req *UpdateUserRequest) (*UserResponse, error) {
	u, err := s.app.UpdateUser(ctx, req.Id, req.Nickname, req.Email)
	if errors.Is(err, app.ErrBadRequest) {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &UserResponse{
		Id:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
	}, nil
}

func (s *Server) DeleteUser(ctx context.Context, req *DeleteUserRequest) (*UserResponse, error) {
	u, err := s.app.DeleteUser(ctx, req.Id)
	if errors.Is(err, app.ErrBadRequest) {
		return nil, status.Error(codes.InvalidArgument, "Invalid argument")
	} else if err != nil {
		return nil, status.Error(codes.Internal, "Internal server error")
	}

	return &UserResponse{
		Id:       u.ID,
		Nickname: u.Nickname,
		Email:    u.Email,
	}, nil
}

// Метод для генерации шаблона для выборки объявлений
func createAdPattern(req *ListAdsRequest) *ads.Pattern {
	f := ads.NewPattern()

	if req.Title != nil {
		f.Title = *req.Title
	}
	if req.Created != nil {
		f.Created = req.Created.AsTime().UTC()
	}
	if req.UserId != nil {
		f.UserID = *req.UserId
	}
	if req.Published != nil {
		f.Published = *req.Published
	}

	return f
}
