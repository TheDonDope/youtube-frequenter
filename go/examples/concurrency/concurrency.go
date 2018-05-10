package main

import (
	"fmt"
	"math/rand"
	"time"
)

var (
	// Web1 Replicae
	Web1 = fakeSearch("web1")
	// Web2 Replicae
	Web2 = fakeSearch("web2")
	// Image1 Replicae
	Image1 = fakeSearch("image1")
	// Image2 Replicae
	Image2 = fakeSearch("image2")
	// Video1 Replicae
	Video1 = fakeSearch("video1")
	// Video2 Replicae
	Video2 = fakeSearch("video2")
)

// Result datatype
type Result struct {
	resultValue string
}

// Search function declaration
type Search func(query string) Result

func fakeSearch(kind string) Search {
	return func(query string) Result {
		time.Sleep(time.Duration(rand.Intn(100)) * time.Millisecond)
		return Result{fmt.Sprintf("%s result for %q\n", kind, query)}
	}
}

func FirstResponder(query string, replicas ...Search) Result {
	resultChannel := make(chan Result)
	searchReplica := func(index int) { resultChannel <- replicas[index](query) }
	for replicaIndex := range replicas {
		go searchReplica(replicaIndex)
	}
	return <-resultChannel
}

func Google(query string) (resultsFromChannel []Result) {
	resultChannel := make(chan Result)
	go func() { resultChannel <- FirstResponder(query, Web1, Web2) }()
	go func() { resultChannel <- FirstResponder(query, Image1, Image2) }()
	go func() { resultChannel <- FirstResponder(query, Video1, Video2) }()

	timeout := time.After(80 * time.Millisecond)
	for i := 0; i < 3; i++ {
		select {
		case resultFromChannel := <-resultChannel:
			resultsFromChannel = append(resultsFromChannel, resultFromChannel)
		case <-timeout:
			fmt.Println("timed out")
			return
		}
	}
	return
}

func main() {
	rand.Seed(time.Now().UnixNano())
	start := time.Now()
	result := FirstResponder("golang", fakeSearch("replica 1"), fakeSearch("replica 2"))
	results := Google("golang")
	elapsed := time.Since(start)
	fmt.Println(result)
	fmt.Println(results)
	fmt.Println(elapsed)
}
