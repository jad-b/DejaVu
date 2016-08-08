package main

import (
	"fmt"
	"log"

	tasks "google.golang.org/api/tasks/v1"
)

func main() {
	client := loadClient("./client_secret.json", tasks.TasksReadonlyScope)

	/* Demonstrate task operations */
	srv, err := tasks.New(client)
	if err != nil {
		log.Fatalf("Unable to retrieve tasks Client %v", err)
	}

	r, err := srv.Tasklists.List().MaxResults(10).Do()
	if err != nil {
		log.Fatalf("Unable to retrieve task lists.", err)
	}

	fmt.Println("Task Lists:")
	if len(r.Items) > 0 {
		for _, i := range r.Items {
			fmt.Printf("%s (%s)\n", i.Title, i.Id)
		}
	} else {
		fmt.Print("No task lists found.")
	}
}
