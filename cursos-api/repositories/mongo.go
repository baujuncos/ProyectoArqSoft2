package repositories

import (
	"context"
	coursesDAO "cursos-api/dao"
	"fmt"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
	"go.mongodb.org/mongo-driver/mongo/options"
	"log"
)

type MongoConfig struct {
	Host       string
	Port       string
	Username   string
	Password   string
	Database   string
	Collection string
}

type Mongo struct {
	client     *mongo.Client
	database   string
	collection string
}

const (
	connectionURI = "mongodb://%s:%s"
)

func NewMongo(config MongoConfig) Mongo {
	credentials := options.Credential{
		Username: config.Username,
		Password: config.Password,
	}

	ctx := context.Background()
	uri := fmt.Sprintf(connectionURI, config.Host, config.Port)
	cfg := options.Client().ApplyURI(uri).SetAuth(credentials)

	client, err := mongo.Connect(ctx, cfg)
	if err != nil {
		log.Panicf("error connecting to mongo DB: %v", err)
	}

	return Mongo{
		client:     client,
		database:   config.Database,
		collection: config.Collection,
	}
}

func (repository Mongo) GetCourseByID(ctx context.Context, id string) (coursesDAO.Course, error) {
	// Get from MongoDB
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return coursesDAO.Course{}, fmt.Errorf("error converting id to mongo ID: %w", err)
	}
	result := repository.client.Database(repository.database).Collection(repository.collection).FindOne(ctx, bson.M{"_id": objectID})
	if result.Err() != nil {
		return coursesDAO.Course{}, fmt.Errorf("error finding document: %w", result.Err())
	}

	// Convert document to DAO
	var courseDAO coursesDAO.Course
	if err := result.Decode(&courseDAO); err != nil {
		return coursesDAO.Course{}, fmt.Errorf("error decoding result: %w", err)
	}
	return courseDAO, nil
}

func (repository Mongo) Create(ctx context.Context, course coursesDAO.Course) (string, error) {
	// Insert into mongo
	result, err := repository.client.Database(repository.database).Collection(repository.collection).InsertOne(ctx, course)
	if err != nil {
		return "", fmt.Errorf("error creating document: %w", err)
	}

	// Get inserted ID
	objectID, ok := result.InsertedID.(primitive.ObjectID)
	if !ok {
		return "", fmt.Errorf("error converting mongo ID to object ID")
	}

	return objectID.Hex(), nil
}

func (repository Mongo) Update(ctx context.Context, course coursesDAO.Course) error {
	// Convert course ID to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(course.Course_id)
	if err != nil {
		return fmt.Errorf("error converting id to mongo ID: %w", err)
	}

	// Create an update document
	update := bson.M{}

	// Only set the fields that are not empty or their default value
	if course.Nombre != "" {
		update["nombre"] = course.Nombre
	}
	if course.Categoria != "" {
		update["categoria"] = course.Categoria
	}
	if !course.Fecha_inicio.IsZero() {
		update["fecha_inicio"] = course.Fecha_inicio
	}
	if course.Duracion != 0 { //asumimos que 0 es el valor por defecto
		update["duracion"] = course.Duracion
	}
	if course.Descripcion != "" {
		update["descripcion"] = course.Descripcion
	}
	if course.Requisitos != "" {
		update["requisitos"] = course.Requisitos
	}
	if course.Valoracion != 0 {
		update["valoracion"] = course.Valoracion
	}
	if course.Url_image != "" {
		update["url_image"] = course.Url_image
	}

	// Update the document in MongoDB
	if len(update) == 0 {
		return fmt.Errorf("no fields to update for course ID %s", course.Course_id)
	}

	filter := bson.M{"_id": objectID}
	result, err := repository.client.Database(repository.database).Collection(repository.collection).UpdateOne(ctx, filter, bson.M{"$set": update})
	if err != nil {
		return fmt.Errorf("error updating document: %w", err)
	}
	if result.MatchedCount == 0 {
		return fmt.Errorf("no document found with ID %s", course.Course_id)
	}

	return nil
}

func (repository Mongo) Delete(ctx context.Context, id string) error {
	// Convert course ID to MongoDB ObjectID
	objectID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return fmt.Errorf("error converting id to mongo ID: %w", err)
	}

	// Delete the document from MongoDB
	filter := bson.M{"_id": objectID}
	result, err := repository.client.Database(repository.database).Collection(repository.collection).DeleteOne(ctx, filter)
	if err != nil {
		return fmt.Errorf("error deleting document: %w", err)
	}
	if result.DeletedCount == 0 {
		return fmt.Errorf("no document found with ID %s", id)
	}

	return nil
}

func (repository Mongo) GetCourses(ctx context.Context) (coursesDAO.Courses, error) {
	cursor, err := repository.client.Database(repository.database).Collection(repository.collection).Find(ctx, bson.M{})
	if err != nil {
		return nil, fmt.Errorf("error getting documents: %w", err)
	}
	defer func(cursor *mongo.Cursor, ctx context.Context) {
		err := cursor.Close(ctx)
		if err != nil {

		}
	}(cursor, ctx)

	var cursos coursesDAO.Courses
	if err := cursor.All(ctx, &cursos); err != nil {
		return nil, fmt.Errorf("error decoding documents: %w", err)
	}
	return cursos, nil
}
