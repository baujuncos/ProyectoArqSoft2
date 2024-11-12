package repositories

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/stevenferrer/solr-go"
	courseDTO "search-api/dto"
	"time"
)

type SolrConfig struct {
	Host       string // Solr host
	Port       string // Solr port
	Collection string // Solr collection name
}

type Solr struct {
	Client     *solr.JSONClient
	Collection string
}

// NewSolr initializes a new Solr client
func NewSolr(config SolrConfig) Solr {
	// Construct the BaseURL using the provided host and port
	baseURL := fmt.Sprintf("http://%s:%s", config.Host, config.Port)
	client := solr.NewJSONClient(baseURL)

	return Solr{
		Client:     client,
		Collection: config.Collection,
	}
}

// Index adds a new hotel document to the Solr collection
func (searchEngine Solr) Index(ctx context.Context, course courseDTO.CourseDto) (string, error) {
	// Prepare the document for Solr
	doc := map[string]interface{}{
		"course_id":    course.Course_id,
		"nombre":       course.Nombre,
		"profesor_id":  course.Profesor_id,
		"categoria":    course.Categoria,
		"descripcion":  course.Descripcion,
		"valoracion":   course.Valoracion,
		"duracion":     course.Duracion,
		"requisitos":   course.Requisitos,
		"url_image":    course.Url_image,
		"fecha_inicio": course.Fecha_inicio,
	}

	// Prepare the index request
	indexRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Index the document in Solr
	body, err := json.Marshal(indexRequest)
	if err != nil {
		return "", fmt.Errorf("error marshaling hotel document: %w", err)
	}

	// Index the document in Solr
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return "", fmt.Errorf("error indexing hotel: %w", err)
	}
	if resp.Error != nil {
		return "", fmt.Errorf("failed to index hotel: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return "", fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return course.Course_id, nil
}

// Update modifies an existing hotel document in the Solr collection
func (searchEngine Solr) Update(ctx context.Context, course courseDTO.CourseDto) error {
	// Prepare the document for Solr
	doc := map[string]interface{}{
		"course_id":    course.Course_id,
		"nombre":       course.Nombre,
		"profesor_id":  course.Profesor_id,
		"categoria":    course.Categoria,
		"descripcion":  course.Descripcion,
		"valoracion":   course.Valoracion,
		"duracion":     course.Duracion,
		"requisitos":   course.Requisitos,
		"url_image":    course.Url_image,
		"fecha_inicio": course.Fecha_inicio,
	}

	// Prepare the update request
	updateRequest := map[string]interface{}{
		"add": []interface{}{doc}, // Use "add" with a list of documents
	}

	// Update the document in Solr
	body, err := json.Marshal(updateRequest)
	if err != nil {
		return fmt.Errorf("error marshaling hotel document: %w", err)
	}

	// Execute the update request using the Update method
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error updating hotel: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to update hotel: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return nil
}

func (searchEngine Solr) Delete(ctx context.Context, id string) error {
	// Prepare the delete request
	docToDelete := map[string]interface{}{
		"delete": map[string]interface{}{
			"id": id,
		},
	}

	// Update the document in Solr
	body, err := json.Marshal(docToDelete)
	if err != nil {
		return fmt.Errorf("error marshaling hotel document: %w", err)
	}

	// Execute the delete request using the Update method
	resp, err := searchEngine.Client.Update(ctx, searchEngine.Collection, solr.JSON, bytes.NewReader(body))
	if err != nil {
		return fmt.Errorf("error deleting hotel: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("failed to index hotel: %v", resp.Error)
	}

	// Commit the changes
	if err := searchEngine.Client.Commit(ctx, searchEngine.Collection); err != nil {
		return fmt.Errorf("error committing changes to Solr: %w", err)
	}

	return nil
}

func (searchEngine Solr) Search(ctx context.Context, query string, limit int, offset int) ([]courseDTO.CourseDto, error) {
	// Prepare the Solr query with limit and offset
	solrQuery := fmt.Sprintf("q=(nombre:%s)&rows=%d&start=%d", query, limit, offset)

	// Execute the search request
	resp, err := searchEngine.Client.Query(ctx, searchEngine.Collection, solr.NewQuery(solrQuery))
	if err != nil {
		return nil, fmt.Errorf("error executing search query: %w", err)
	}
	if resp.Error != nil {
		return nil, fmt.Errorf("failed to execute search query: %v", resp.Error)
	}

	// Parse the response and extract hotel documents
	var coursesList []courseDTO.CourseDto
	for _, doc := range resp.Response.Documents {
		// Safely extract course fields with type assertions
		course := courseDTO.CourseDto{
			Course_id:    getStringField(doc, "course_id"),
			Nombre:       getStringField(doc, "nombre"),
			Profesor_id:  getStringField(doc, "profesor_id"),
			Categoria:    getStringField(doc, "categoria"),
			Descripcion:  getStringField(doc, "descripcion"),
			Valoracion:   getFloatField(doc, "valoracion"),
			Duracion:     getIntField(doc, "duracion"),
			Requisitos:   getStringField(doc, "requisitos"),
			Url_image:    getStringField(doc, "url_image"),
			Fecha_inicio: getTimeField(doc, "fecha_inicio"),
		}
		coursesList = append(coursesList, course)
	}
	return coursesList, nil
}

// Helper function to safely get int fields from the document
func getIntField(doc map[string]interface{}, field string) int {
	if val, ok := doc[field].(int); ok {
		return val
	}
	if val, ok := doc[field].(float64); ok { // Solr puede devolver números como float64
		return int(val)
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if intVal, ok := val[0].(int); ok {
			return intVal
		}
		if floatVal, ok := val[0].(float64); ok { // Manejo de float64 en listas
			return int(floatVal)
		}
	}
	return 0
}

// Helper function to safely get time fields from the document
func getTimeField(doc map[string]interface{}, field string) time.Time {
	if val, ok := doc[field].(string); ok {
		parsedTime, err := time.Parse(time.RFC3339, val)
		if err == nil {
			return parsedTime
		}
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if strVal, ok := val[0].(string); ok {
			parsedTime, err := time.Parse(time.RFC3339, strVal)
			if err == nil {
				return parsedTime
			}
		}
	}
	return time.Time{} // Valor por defecto si no se puede parsear
}

// Helper function to safely get string fields from the document
func getStringField(doc map[string]interface{}, field string) string {
	if val, ok := doc[field].(string); ok {
		return val
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if strVal, ok := val[0].(string); ok {
			return strVal
		}
	}
	return ""
}

// Helper function to safely get float64 fields from the document
func getFloatField(doc map[string]interface{}, field string) float64 {
	if val, ok := doc[field].(float64); ok {
		return val
	}
	if val, ok := doc[field].([]interface{}); ok && len(val) > 0 {
		if floatVal, ok := val[0].(float64); ok {
			return floatVal
		}
	}
	return 0.0
}
