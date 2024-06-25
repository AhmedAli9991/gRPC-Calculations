package handlers

import (
	"database/sql"
	"encoding/json"
	"log"

	"github.com/ahmedalialphasquad123/calculationService/database"
	pb "github.com/ahmedalialphasquad123/calculationService/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

func HandleAggregationObject(req *pb.CalculationRequest) (*pb.CalculationResponse, error) {
	log.Printf("request: %v", req)

	db, err := database.Connect()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}
	defer db.Close()

	rows, err := db.Query(`SELECT * FROM transaction`)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return nil, err
	}
	defer rows.Close()

	transactions, err := processRows(rows)
	if err != nil {
		return nil, err
	}

	transJSON, err := json.Marshal(transactions)
	if err != nil {
		log.Printf("Error marshaling JSON: %v", err)
		return nil, err
	}
	// log.Printf("JSON: %v", *rows)

	// log.Printf("JSON: %v", transactions)

	return &pb.CalculationResponse{
		Data: &structpb.Struct{
			Fields: map[string]*structpb.Value{
				"status": structpb.NewStringValue("success"),
				"data":   structpb.NewStringValue(string(transJSON)),
			},
		},
	}, nil
}

func processRows(rows *sql.Rows) ([]map[string]interface{}, error) {
	var transactions []map[string]interface{}
	columns, err := rows.Columns()
	if err != nil {
		return nil, err
	}

	for rows.Next() {
		values := make([]interface{}, len(columns))
		valuePtrs := make([]interface{}, len(columns))
		for i := range values {
			valuePtrs[i] = &values[i]
		}

		if err := rows.Scan(valuePtrs...); err != nil {
			return nil, err
		}

		entry := make(map[string]interface{})
		for i, col := range columns {
			entry[col] = values[i]
		}
		transactions = append(transactions, entry)
	}
	if err := rows.Err(); err != nil {
		return nil, err
	}

	return transactions, nil
}

func HandleRequest(req *pb.CalculationRequest) (*pb.CalculationResponse, error) {
	responseData := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"status": structpb.NewStringValue("success"),
		},
	}
	return &pb.CalculationResponse{Data: responseData}, nil
}
