package main

import (
	"fmt"
	"log"
	"os"

	"github.com/joho/godotenv"
	exporter "github.com/khuenqdev/cognito_exporter"
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env in local env `" + err.Error() + "`")
		panic("config file error")
	}

	region, appId, accessKey, secretKey, userPoolId := os.Getenv("REGION"), os.Getenv("APP_ID"), os.Getenv("ACCESS_KEY"), os.Getenv("SECRET_KEY"), os.Getenv("USER_POOL_ID")
	fmt.Printf("REGION:%s\nAPP_ID:%s\nACCESS_KEY:%s\nSECRET_KEY:%s\nUSER_POOL_ID:%s\n", region, appId, accessKey, secretKey, userPoolId)

	e := exporter.NewCognitoExporter(exporter.CognitoExporterConfig{
		Region:     region,
		AccessKey:  accessKey,
		SecretKey:  secretKey,
		AppId:      appId,
		UserPoolId: userPoolId,
	}, getSampleAttributes())

	e.ExportToCSV()
}

// getSampleAttributes returns the list of required attributes for CSV header
func getSampleAttributes() []string {
	// NOTE: Attributes presented here must follow AWS import template, which is different between user pools
	return []string{
		"profile",
		"address",
		"birthdate",
		"gender",
		"preferred_username",
		"updated_at",
		"website",
		"picture",
		"phone_number",
		"phone_number_verified",
		"zoneinfo",
		"locale",
		"email",
		"email_verified",
		"given_name",
		"family_name",
		"middle_name",
		"name",
		"nickname",
		"cognito:mfa_enabled",
		"cognito:username",
	}
}
