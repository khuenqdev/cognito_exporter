package cognitoexporter

import (
	"context"
	"encoding/csv"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
)

const (
	TIMESTAMP_FORMAT = "20060102_150405"
	DEFAULT_LIMIT    = 60
)

type ICognitoExporter interface {
	ExportToCSV()
}

type CognitoExporter struct {
	Config             CognitoExporterConfig
	RequiredAttributes []string
}

type CognitoExporterConfig struct {
	Region     string
	AccessKey  string
	SecretKey  string
	AppId      string
	UserPoolId string
}

func NewCognitoExporter(config CognitoExporterConfig, requiredAttributes []string) ICognitoExporter {
	return CognitoExporter{config, requiredAttributes}
}

func (e CognitoExporter) ExportToCSV() {
	// Create Cognito client
	client := cognito.New(cognito.Options{
		Region:           e.Config.Region,
		Credentials:      aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(e.Config.AccessKey, e.Config.SecretKey, "")),
		AppID:            e.Config.AppId,
		RetryMaxAttempts: 5,
	})

	// Create a context for export
	ctx := context.Background()

	// Initialize CSV output
	timestamp := time.Now().Format(TIMESTAMP_FORMAT)
	filename := fmt.Sprintf("cognito_users_%s.csv", timestamp)
	file, err := os.Create(filename)
	if err != nil {
		log.Fatalf("Cannot create file: %v", err)
		return
	}
	defer file.Close()

	writer := csv.NewWriter(file)
	defer writer.Flush()

	// Write headers to CSV
	header := e.RequiredAttributes
	if err := writer.Write(header); err != nil {
		log.Fatalf("Error writing header to CSV: %v\n", err)
		return
	}

	// Get paginated user list object
	paginator := cognito.NewListUsersPaginator(client, &cognito.ListUsersInput{
		UserPoolId: aws.String(e.Config.UserPoolId),
		Limit:      aws.Int32(DEFAULT_LIMIT),
	})

	// Track number of users processed
	userCount := 0
	pageCount := 0

	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			log.Fatalf("Failed to list users: %v\n", err)
			return
		}
		pageCount++

		log.Printf("Exporting users of page %d\n", pageCount)

		for _, user := range output.Users {
			// Prepare row data
			row := make([]string, len(header))

			// Get other attributes
			for i, attrName := range header {
				row[i] = e.getAttribute(user.Attributes, attrName)
			}

			// Write the row to CSV
			if err := writer.Write(row); err != nil {
				log.Printf("Error writing user %s to CSV: %v\n", *user.Username, err)
				continue
			}

			userCount++
		}
	}

	log.Println("Total number of users exported: ", userCount)
}

// getAttribute safely gets attribute value from user attributes
func (e CognitoExporter) getAttribute(attributes []types.AttributeType, attrName string) string {
	for _, attr := range attributes {
		if *attr.Name == attrName {
			return *attr.Value
		}
	}
	return ""
}
