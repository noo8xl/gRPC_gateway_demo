package cache_test

import (
	"testing"
	"time"

	reviewPb "github.com/noo8xl/anvil-api/main/reviews"
)

var (
	reviewId uint64 = 1

	reviewCommentsList = &reviewPb.ReviewCommentsList{
		Comments: []*reviewPb.ReviewCommentResponse{
			{
				Comment: &reviewPb.ReviewComment{
					ReviewId:  1,
					PostedBy:  4,
					Title:     "test title",
					Body:      "test comment",
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				},
			},
			{
				Comment: &reviewPb.ReviewComment{
					ReviewId:  2,
					PostedBy:  4,
					Title:     "test title 2",
					Body:      "test comment 2",
					CreatedAt: time.Now().Format(time.RFC3339),
					UpdatedAt: time.Now().Format(time.RFC3339),
				},
			},
		},
	}

	reviewDetails = &reviewPb.ReviewResponse{
		Review: &reviewPb.Review{
			ReviewId:   1,
			CustomerId: 4,
			ReviewerId: 5,
		},
		Details: &reviewPb.ReviewDetails{
			Rank:      4.3,
			Title:     "test title",
			Body:      "test body",
			OrderId:   1,
			CreatedAt: time.Now().Format(time.RFC3339),
			UpdatedAt: time.Now().Format(time.RFC3339),
			Helpful:   0,
			Useless:   0,
		},
	}
)

// ###########################################
// ############## comments area ##############
// ###########################################

func TestSetReviewCommentsList(t *testing.T) {

	err := svc.SetReviewCommentsList(reviewId, reviewCommentsList)
	if err != nil {
		t.Errorf("error setting review comments list: Expected nil, got %v", err)
	}
}

func TestGetReviewCommentsList(t *testing.T) {

	reviewCommentsList, err := svc.GetReviewCommentsList(reviewId)
	if err != nil {
		t.Errorf("error getting review comments list: Expected nil, got %v", err)
	}

	if reviewCommentsList != nil {
		if len(reviewCommentsList.Comments) == 0 {
			t.Errorf("error getting review comments list: Expected a list of comments, got nil")
		}
	}
}

// ###########################################
// ############## details area ##############
// ###########################################

func TestSetReviewDetails(t *testing.T) {
	err := svc.SetReviewDetails(reviewId, reviewDetails)
	if err != nil {
		t.Errorf("error setting review details: Expected nil, got %v", err)
	}
}

func TestGetReviewDetails(t *testing.T) {
	details, err := svc.GetReviewDetails(reviewId)
	if err != nil {
		t.Errorf("error getting review details: Expected nil, got %v", err)
	}

	if details != nil {
		if details.Review == nil {
			t.Errorf("error getting review details: Expected a review, got nil")
		}
	}
}

func TestClearReviewCommentsList(t *testing.T) {
	err := svc.ClearReviewCommentsList(reviewId)
	if err != nil {
		t.Errorf("error clearing review comments list: Expected nil, got %v", err)
	}
}

func TestClearReviewDetails(t *testing.T) {
	err := svc.ClearReviewDetails(reviewId)
	if err != nil {
		t.Errorf("error clearing review details: Expected nil, got %v", err)
	}
}
