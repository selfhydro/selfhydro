package main

import (
	"fmt"
	"os"
	"time"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

func main() {
	lambda.Start(createTable)
}

func createTable() {
	session := createSession()
	createNewTable(session)
}

func createNewTable(svc *dynamodb.DynamoDB) {
	tableName := getTableName(time.Now().Add(time.Hour * 24))

	input := getTableStructure(tableName)

	_, err := svc.CreateTable(input)
	if err != nil {
		fmt.Println("Got error calling CreateTable:")
		fmt.Println(err.Error())
		os.Exit(1)
	}

	fmt.Println("Created the table", tableName)
}

func createSession() *dynamodb.DynamoDB {
	sess := session.Must(session.NewSessionWithOptions(session.Options{
		SharedConfigState: session.SharedConfigEnable,
	}))

	svc := dynamodb.New(sess)
	return svc
}

func getTableName(time time.Time) string {
	baseName := "selfhydro-state-"
	return baseName + time.Format("2006-01")
}

func getTableStructure(tableName string) *dynamodb.CreateTableInput {
	return &dynamodb.CreateTableInput{
		AttributeDefinitions: []*dynamodb.AttributeDefinition{
			{
				AttributeName: aws.String("Date"),
				AttributeType: aws.String("S"),
			},
			{
				AttributeName: aws.String("SystemID"),
				AttributeType: aws.String("S"),
			},
		},
		KeySchema: []*dynamodb.KeySchemaElement{
			{
				AttributeName: aws.String("SystemID"),
				KeyType:       aws.String("HASH"),
			},
			{
				AttributeName: aws.String("Date"),
				KeyType:       aws.String("RANGE"),
			},
		},
		ProvisionedThroughput: &dynamodb.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(2),
			WriteCapacityUnits: aws.Int64(1),
		},
		TableName: aws.String(tableName),
	}
}
