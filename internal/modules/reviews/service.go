package reviews

import (
	"context"

	apperrors "github.com/DavidMovas/Movies-Reviews/internal/error"

	"github.com/DavidMovas/Movies-Reviews/internal/log"
)

type Service struct {
	repo *Repository
}

func NewService(repo *Repository) *Service {
	return &Service{
		repo: repo,
	}
}

func (s *Service) GetReviewsByMovieID(ctx context.Context, movieID int, offset int, limit int, sort, order string) ([]*Review, int, error) {
	return s.repo.GetReviewsByMovieID(ctx, movieID, offset, limit, sort, order)
}

func (s *Service) GetReviewsByUserID(ctx context.Context, userID int, offset int, limit int, sort, order string) ([]*Review, int, error) {
	return s.repo.GetReviewsByUserID(ctx, userID, offset, limit, sort, order)
}

func (s *Service) GetReviewByID(ctx context.Context, reviewID int) (*Review, error) {
	return s.repo.GetReviewByID(ctx, reviewID)
}

func (s *Service) CreateReview(ctx context.Context, req *CreateReviewRequest) (*Review, error) {
	review, err := s.repo.CreateReview(ctx, req)
	if err != nil {
		return nil, err
	}

	log.FromContext(ctx).Info("review created", "review_id", review.ID)

	return review, nil
}

func (s *Service) UpdateReview(ctx context.Context, reviewID int, req *UpdateReviewRequest) (*Review, error) {
	review, err := s.repo.UpdateReview(ctx, reviewID, req)
	if err != nil {
		return nil, err
	}

	if review.UserID != req.UserID {
		return nil, apperrors.Forbidden("insufficient permissions")
	}

	log.FromContext(ctx).Info("review updated", "review_id", review.ID)

	return review, nil
}

func (s *Service) DeleteReview(ctx context.Context, reviewID int) error {
	if err := s.repo.DeleteReview(ctx, reviewID); err != nil {
		return err
	}

	log.FromContext(ctx).Info("review deleted", "review_id", reviewID)
	return nil
}
