package mongo

import (
	"context"
	"errors"
	"fmt"
	"strings"
	"time"

	"akiba/backend/internal/domain"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type UserRepository struct {
	collection *mongo.Collection
	timeout    time.Duration
}

func NewUserRepository(db *mongo.Database, timeout time.Duration) *UserRepository {
	return &UserRepository{collection: db.Collection("users"), timeout: timeout}
}

func (r *UserRepository) EnsureIndexes(ctx context.Context) error {
	models := []mongo.IndexModel{
		{Keys: bson.D{{Key: "emailLower", Value: 1}}, Options: options.Index().SetName("uniq_emailLower").SetUnique(true)},
		{Keys: bson.D{{Key: "phoneE164", Value: 1}}, Options: options.Index().SetName("uniq_phoneE164").SetUnique(true)},
		{Keys: bson.D{{Key: "usernameLower", Value: 1}}, Options: options.Index().SetName("uniq_usernameLower").SetUnique(true)},
	}
	_, err := r.collection.Indexes().CreateMany(ctx, models)
	return err
}

func (r *UserRepository) Create(ctx context.Context, user *domain.User) error {
	cctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	doc := bson.M{"emailLower": user.EmailLower, "phoneE164": user.PhoneE164, "usernameLower": user.UsernameLower, "passwordHash": user.PasswordHash, "status": user.Status, "createdAt": user.CreatedAt, "updatedAt": user.UpdatedAt}
	res, err := r.collection.InsertOne(cctx, doc)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			return domain.ErrUserExists
		}
		return err
	}
	id, ok := res.InsertedID.(primitive.ObjectID)
	if !ok {
		return fmt.Errorf("invalid inserted id")
	}
	user.ID = id.Hex()
	return nil
}

func (r *UserRepository) GetByID(ctx context.Context, id string) (*domain.User, error) {
	objID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	cctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	var out struct {
		ID            primitive.ObjectID `bson:"_id"`
		EmailLower    string             `bson:"emailLower"`
		PhoneE164     string             `bson:"phoneE164"`
		UsernameLower string             `bson:"usernameLower"`
		PasswordHash  string             `bson:"passwordHash"`
		Status        domain.UserStatus  `bson:"status"`
		CreatedAt     time.Time          `bson:"createdAt"`
		UpdatedAt     time.Time          `bson:"updatedAt"`
	}
	err = r.collection.FindOne(cctx, bson.M{"_id": objID}).Decode(&out)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domain.User{ID: out.ID.Hex(), EmailLower: out.EmailLower, PhoneE164: out.PhoneE164, UsernameLower: out.UsernameLower, PasswordHash: out.PasswordHash, Status: out.Status, CreatedAt: out.CreatedAt.UTC(), UpdatedAt: out.UpdatedAt.UTC()}, nil
}

func (r *UserRepository) GetByLogin(ctx context.Context, login string) (*domain.User, error) {
	cctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()
	filter := bson.M{}
	if strings.Contains(login, "@") {
		filter["emailLower"] = strings.ToLower(login)
	} else if strings.HasPrefix(login, "+") {
		filter["phoneE164"] = login
	} else {
		filter["usernameLower"] = strings.ToLower(login)
	}
	var out struct {
		ID            primitive.ObjectID `bson:"_id"`
		EmailLower    string             `bson:"emailLower"`
		PhoneE164     string             `bson:"phoneE164"`
		UsernameLower string             `bson:"usernameLower"`
		PasswordHash  string             `bson:"passwordHash"`
		Status        domain.UserStatus  `bson:"status"`
		CreatedAt     time.Time          `bson:"createdAt"`
		UpdatedAt     time.Time          `bson:"updatedAt"`
	}
	err := r.collection.FindOne(cctx, filter).Decode(&out)
	if errors.Is(err, mongo.ErrNoDocuments) {
		return nil, domain.ErrUserNotFound
	}
	if err != nil {
		return nil, err
	}
	return &domain.User{ID: out.ID.Hex(), EmailLower: out.EmailLower, PhoneE164: out.PhoneE164, UsernameLower: out.UsernameLower, PasswordHash: out.PasswordHash, Status: out.Status, CreatedAt: out.CreatedAt.UTC(), UpdatedAt: out.UpdatedAt.UTC()}, nil
}
