package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

type UpdatePostPayload struct {
	Title   *string `json:"title"`
	Content *string `json:"content"`
}

func updatePost(postID int, p UpdatePostPayload, wg *sync.WaitGroup) {
	defer wg.Done()
	url := fmt.Sprintf("http://localhost:8080/v1/posts/%d", postID)
	b, _ := json.Marshal(p)
	req, err := http.NewRequest(http.MethodPatch, url, bytes.NewBuffer(b))
	if err != nil {
		panic(err)
	}

	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("error sending request:", err)
		return
	}
	defer resp.Body.Close()
	fmt.Println("update post response status: ", resp.Status)
}

func main() {
	var wg sync.WaitGroup

	postID := 7
	//simulate user A and user B updating the same post concurrently
	wg.Add(2)
	title := "NEW TITLE FROM USER A"
	content := "NEW CONTENT FROM USER B"
	// title := "title"
	// content := "content"
	// updatePost(postID, UpdatePostPayload{
	// 	Title:   &title,
	// 	Content: &content,
	// }, &wg)
	go updatePost(postID, UpdatePostPayload{
		Title: &title,
	}, &wg)
	go updatePost(postID, UpdatePostPayload{

		Content: &content,
	}, &wg)
	wg.Wait()
}
