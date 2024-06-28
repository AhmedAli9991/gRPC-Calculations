package handlers

import (
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"github.com/ahmedalialphasquad123/calculationService/database"
	"github.com/ahmedalialphasquad123/calculationService/utils"

	pb "github.com/ahmedalialphasquad123/calculationService/proto"
	"google.golang.org/protobuf/types/known/structpb"
)

type Transaction struct {
	ID       int
	Seller   string
	Buyer    string
	Material string
	Product  string
	Status   string
}

func HandleAggregationObject(req *pb.AggregationObjectRequest) (*pb.CalculationResponse, error) {
	log.Printf("request: %v", req.Material)
	db, err := database.Connect()
	if err != nil {
		log.Printf("Error connecting to database: %v", err)
		return nil, err
	}
	defer db.Close()

	conditions := []string{fmt.Sprintf("business_entity_id = %d", req.BusinessEntityId)}

	if req.StartDate != "" {
		date, err := utils.TimeFormat(req.StartDate)
		if err != nil {
			log.Printf("Error formating date: %v", err)
			return nil, err
		}
		conditions = append(conditions, fmt.Sprintf("transaction_date >= '%s'", date))
	}

	if req.EndDate != "" {
		date, err := utils.TimeFormat(req.EndDate)
		if err != nil {
			log.Printf("Error formating date: %v", err)
			return nil, err
		}
		conditions = append(conditions, fmt.Sprintf("transaction_date <= '%s'", date))
	}

	if len(req.Seller) > 0 {
		conditions = append(conditions, fmt.Sprintf("seller IN ('%s')", strings.Join(req.Seller, "','")))
	}

	if len(req.Buyer) > 0 {
		conditions = append(conditions, fmt.Sprintf("buyer IN ('%s')", strings.Join(req.Buyer, "','")))
	}

	var materialArray []string
	if len(req.Material) > 0 {
		if strings.Contains(req.Material[0], "|") {
			materialArray = strings.Split(req.Material[0], "|")
		} else {
			materialArray = req.Material
		}
	}

	if len(materialArray) > 0 {
		materialCondition := fmt.Sprintf("material IN ('%s')", strings.Join(materialArray, "','"))
		conditions = append(conditions, fmt.Sprintf("(product_component_link.material IN (%s) OR verified_transaction.material IN (%s))", materialCondition, materialCondition))
	}

	if req.Status != "" {
		conditions = append(conditions, fmt.Sprintf("verified = '%s'", req.Status))
	}

	if req.Search != "" {
		searchCondition := fmt.Sprintf("(identifier LIKE '%%%s%%' OR seller LIKE '%%%s%%' OR buyer LIKE '%%%s%%' OR material LIKE '%%%s%%')", req.Search, req.Search, req.Search, req.Search)
		conditions = append(conditions, searchCondition)
	}

	query := fmt.Sprintf("SELECT * FROM transaction WHERE %s", strings.Join(conditions, " AND "))

	if req.Skip > 0 {
		query += fmt.Sprintf(" OFFSET %d", req.Skip)
	}

	if req.Take > 0 {
		query += fmt.Sprintf(" LIMIT %d", req.Take)
	}

	rows, err := db.Query(query)
	log.Printf("Query : %v", query)
	if err != nil {
		log.Printf("Error querying database: %v", err)
		return nil, err
	}
	defer rows.Close()

	var transactions []Transaction
	for rows.Next() {
		var transaction Transaction
		if err := rows.Scan(&transaction.ID, &transaction.Seller, &transaction.Buyer, &transaction.Material, &transaction.Product, &transaction.Status); err != nil {
			log.Printf("Error scanning row: %v", err)
			return nil, err
		}
		transactions = append(transactions, transaction)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over rows: %v", err)
		return nil, err
	}

	if len(transactions) == 0 {
		return nil, fmt.Errorf("no transactions found")
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
	// return transactions, nil
}

func HandleRequest(req *pb.CalculationRequest) (*pb.CalculationResponse, error) {
	responseData := &structpb.Struct{
		Fields: map[string]*structpb.Value{
			"status": structpb.NewStringValue("success"),
		},
	}
	return &pb.CalculationResponse{Data: responseData}, nil
}
