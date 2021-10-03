package seed

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/dish.io/internal/domain"
	"github.com/dish.io/internal/services/user"
	"os"
)

type FileData struct {
	Users []domain.User `json:"users"`
}

func HandleSeed(ctx context.Context, seedCmd *flag.FlagSet, filename *string, store *user.Service) {
	_ = seedCmd.Parse(os.Args[2:])

	// Opening seed file
	seedFile, err := os.OpenFile(*filename, os.O_RDONLY, os.ModePerm)
	if err != nil {
		fmt.Printf("Error reading the file: %s", err.Error())
		os.Exit(0)
	}
	//Decoding data from json file

	fileData := FileData{}
	if err := json.NewDecoder(seedFile).Decode(&fileData); err != nil {
		fmt.Printf("Error decoding file: %s", err.Error())
		os.Exit(0)
	}
	indentedData, _ := json.MarshalIndent(fileData, "", "    ")
	fmt.Println(string(indentedData))

	//Inserting data into db

	numWorkers := 5
	pipe := make(chan domain.User, numWorkers)
	result := make(chan string, len(fileData.Users))
	// Generating workers
	for i := 0; i < numWorkers; i++ {
		go func(user <-chan domain.User, output chan string) {
			u := <-user
			u.PasswordHash = u.Username
			_, err := store.CreateUser(ctx, u.Email, u.Username, u.Username)
			if err != nil {
				fmt.Printf("Error creating user: %s", err.Error())
				os.Exit(0)
			}
			result <- fmt.Sprintf("User %s is created", u.Username)
		}(pipe, result)
	}
	for _, v := range fileData.Users {
		pipe <- v
	}
	for i := 0; i < len(fileData.Users); i++ {
		fmt.Println(<-result)
	}
	close(pipe)
	close(result)
	fmt.Printf("Data inserted with sucess: %d lignes\n", len(fileData.Users))
	os.Exit(1)
}
