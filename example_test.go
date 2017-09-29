// Idea from https://github.com/sentientmonkey/future/blob/master/example_test.go

package future_test

import (
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"regexp"

	"cclib/future"
)

func ExampleFuture() {
	f := future.NewFuture(func() (interface{}, error) {
		return http.Get("http://golang.org/")
	})

	result, err := f.Wait(nil)
	if err != nil {
		fmt.Printf("Got error: %s\n", err)
		return
	}

	response := result.(*http.Response)
	defer response.Body.Close()
	fmt.Printf("Got result: %d\n", response.StatusCode)
	// Output: Got result: 200
}

func ExamplePromise() {
	p := future.NewPromise(func() (interface{}, error) {
		return http.Get("http://golang.org/")
	})

	p = p.Then(func(value interface{}) (interface{}, error) {
		response := value.(*http.Response)
		defer response.Body.Close()
		b, err := ioutil.ReadAll(response.Body)
		if err != nil {
			return nil, err
		}

		return string(b), nil
	})

	p = p.Then(func(value interface{}) (interface{}, error) {
		body := value.(string)
		r, err := regexp.Compile("<title>(.*)</title>")
		if err != nil {
			return nil, err
		}
		match := r.FindStringSubmatch(body)

		if len(match) < 1 {
			return nil, errors.New("Title not found")
		}

		return match[1], nil
	})

	result, err := p.Wait(nil)
	if err != nil {
		fmt.Printf("Got error: %s\n", err)
		return
	}

	s := result.(string)
	fmt.Printf("Got result: %s\n", s)
	// Output: Got result: The Go Programming Language
}
