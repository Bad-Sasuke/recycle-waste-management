package repositories

import (
	"context"
	"os"
	"recycle-waste-management-backend/src/domain/datasources"
	"recycle-waste-management-backend/src/domain/models"
	"time"

	"github.com/google/uuid"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type IReviewRepository interface {
	CreateReview(ctx context.Context, review *models.ReviewModel) error
	GetReviewByCustomerRequestID(ctx context.Context, customerRequestID string) (*models.ReviewModel, error)
	GetReviewsByShopID(ctx context.Context, shopID string, page, pageSize int) ([]models.ReviewModel, int64, error)
	CheckReviewExists(ctx context.Context, customerRequestID string) (bool, error)
	GetShopRatingStats(ctx context.Context, shopID string) (float64, int64, error)
}

type reviewRepository struct {
	collection *mongo.Collection
}

func NewReviewRepository(db *datasources.MongoDB) IReviewRepository {
	return &reviewRepository{
		collection: db.MongoDB.Database(os.Getenv("DATABASE_NAME")).Collection("reviews"),
	}
}

func (r *reviewRepository) CreateReview(ctx context.Context, review *models.ReviewModel) error {
	review.ReviewID = uuid.New().String()
	review.CreatedAt = time.Now()
	review.UpdatedAt = time.Now()

	_, err := r.collection.InsertOne(ctx, review)
	return err
}

func (r *reviewRepository) GetReviewByCustomerRequestID(ctx context.Context, customerRequestID string) (*models.ReviewModel, error) {
	var review models.ReviewModel
	filter := bson.M{"customer_request_id": customerRequestID}

	err := r.collection.FindOne(ctx, filter).Decode(&review)
	if err == mongo.ErrNoDocuments {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	return &review, nil
}

func (r *reviewRepository) GetReviewsByShopID(ctx context.Context, shopID string, page, pageSize int) ([]models.ReviewModel, int64, error) {
	// Calculate skip
	skip := (page - 1) * pageSize

	// Filter for non-skipped reviews only
	filter := bson.M{
		"shop_id":    shopID,
		"is_skipped": false,
	}

	// Count total documents
	totalCount, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return nil, 0, err
	}

	// Find with pagination, sorted by created_at descending
	opts := options.Find().
		SetSkip(int64(skip)).
		SetLimit(int64(pageSize)).
		SetSort(bson.D{{Key: "created_at", Value: -1}})

	cursor, err := r.collection.Find(ctx, filter, opts)
	if err != nil {
		return nil, 0, err
	}
	defer cursor.Close(ctx)

	var reviews []models.ReviewModel
	if err := cursor.All(ctx, &reviews); err != nil {
		return nil, 0, err
	}

	return reviews, totalCount, nil
}

func (r *reviewRepository) CheckReviewExists(ctx context.Context, customerRequestID string) (bool, error) {
	filter := bson.M{"customer_request_id": customerRequestID}
	count, err := r.collection.CountDocuments(ctx, filter)
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

func (r *reviewRepository) GetShopRatingStats(ctx context.Context, shopID string) (float64, int64, error) {
	matchStage := bson.D{{Key: "$match", Value: bson.M{"shop_id": shopID, "is_skipped": false}}}
	groupStage := bson.D{{Key: "$group", Value: bson.M{
		"_id":           nil,
		"averageRating": bson.M{"$avg": "$rating"},
		"totalReviews":  bson.M{"$sum": 1},
	}}}

	cursor, err := r.collection.Aggregate(ctx, mongo.Pipeline{matchStage, groupStage})
	if err != nil {
		return 0, 0, err
	}
	defer cursor.Close(ctx)

	var results []bson.M

	if err := cursor.All(ctx, &results); err != nil {
		return 0, 0, err
	}

	if len(results) > 0 {
		result := results[0]
		var avgRating float64
		var totalReviews int64

		// Safely convert averageRating
		if val, ok := result["averageRating"]; ok {
			switch v := val.(type) {
			case int32:
				avgRating = float64(v)
			case int64:
				avgRating = float64(v)
			case float64:
				avgRating = v
			case float32:
				avgRating = float64(v)
			default:
				avgRating = 0
			}
		}

		// Safely convert totalReviews
		if val, ok := result["totalReviews"]; ok {
			switch v := val.(type) {
			case int32:
				totalReviews = int64(v)
			case int64:
				totalReviews = v
			case float64:
				totalReviews = int64(v)
			case float32:
				totalReviews = int64(v)
			default:
				totalReviews = 0
			}
		}

		return avgRating, totalReviews, nil
	}

	return 0, 0, nil
}
