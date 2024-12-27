package main

import (
	"context"
	"log"
	"fmt"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	cognito "github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider"
	"github.com/aws/aws-sdk-go-v2/service/cognitoidentityprovider/types"
	"github.com/joho/godotenv"
)

const (
	REGION = ""
	ACCESS_KEY = ""
	SECRET_KEY = ""
	APP_ID = ""
)

func main() {
	err := godotenv.Load("config.env")
	if err != nil {
		log.Fatal("Error loading config.env in local env `" + err.Error() + "`")
		panic("config file error")
	}

	fmt.Println("==> Creating Cognito client")
	client := cognito.New(cognito.Options{
		Region: REGION,
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(ACCESS_KEY, SECRET_KEY, "")),
		AppID: APP_ID,
		RetryMaxAttempts: 5,
	})

	fmt.Println("==> Listing user pools")
	ctx := context.Background()
	var pools []types.UserPoolDescriptionType
	paginator := cognito.NewListUserPoolsPaginator(client, &cognito.ListUserPoolsInput{MaxResults: aws.Int32(10)})
	for paginator.HasMorePages() {
		output, err := paginator.NextPage(ctx)
		if err != nil {
			log.Printf("Couldn't get user pools. Here's why: %v\n", err)
			panic(err)
		} else {
			pools = append(pools, output.UserPools...)
		}
	}
	if len(pools) == 0 {
		fmt.Println("You don't have any user pools!")
	} else {
		for _, pool := range pools {
			fmt.Printf("\t%v: %v\n", *pool.Name, *pool.Id)
		}
	}
}
