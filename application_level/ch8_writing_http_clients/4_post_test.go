package ch8writinghttpclients

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

type User struct {
	First string
	Last  string
}

func handlePostUser(t *testing.T) func(http.ResponseWriter, *http.Request) {
	return func(w http.ResponseWriter, r *http.Request) {
		defer func(r io.ReadCloser) {
			_, _ = io.Copy(io.Discard, r)
			_ = r.Close()
		}(r.Body)

		if r.Method != http.MethodPost {
			http.Error(w, "", http.StatusMethodNotAllowed)
			return
		}

		var u User
		// read json from body then decode to u(User)
		err := json.NewDecoder(r.Body).Decode(&u)

		if err != nil {
			t.Error(err)
			http.Error(w, "Decode Failed", http.StatusBadRequest)
			return
		}

		response := struct {
			Message string `json:"message"`
			User    User   `json:"user"`
		}{
			Message: "User received succesfully",
			User:    u,
		}

		w.Header().Set("Content-Type", "application/json")

		w.WriteHeader(http.StatusAccepted)

		if err := json.NewEncoder(w).Encode(response); err != nil {
			t.Error(err)
		}

	}
}

func TestPostUser(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(handlePostUser(t)))

	defer ts.Close()

	resp, err := http.Get(ts.URL)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusMethodNotAllowed {
		t.Fatalf("expected status %d; actual status %d", http.StatusMethodNotAllowed, resp.StatusCode)
	}

	// buf is our req body carry our payload
	buf := new(bytes.Buffer)
	u := User{First: "Adam", Last: "Woodbeck"}

	err = json.NewEncoder(buf).Encode(&u)

	if err != nil {
		t.Fatal(err)
	}

	// in real scenario, we send content type to handle the logic of decode/encode (choose which type, in this case is json)
	resp, err = http.Post(ts.URL, "application/json", buf)

	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusAccepted {
		t.Fatalf("expected status %d; actual status %d", http.StatusAccepted, resp.StatusCode)
	}

	var result struct {
		Message string `json:"message"`
		User    User   `json:"user"`
	}

	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		t.Fatal(err)
	}
	fmt.Println(result.Message)
	fmt.Printf("%+v\n", result.User)
	_ = resp.Body.Close()
}

// this show how you carry multi part in request body
func TestMultipartPost(t *testing.T) {
	reqBody := new(bytes.Buffer)
	w := multipart.NewWriter(reqBody)

	// iterate a map with 2 (k,v) then put each (k,v) - part to writer.
	for k, v := range map[string]string{
		"date":        time.Now().Format(time.RFC3339),
		"description": "Form values with attached files",
	} {
		// (k,v) as form
		err := w.WriteField(k, v)
		if err != nil {
			t.Fatal(err)
		}
	}

	// txt file as parts
	for i, file := range []string{
		"./files/hello.txt",
		"./files/goodbye.txt",
	} {
		filePart, err := w.CreateFormFile(fmt.Sprintf("file%d", i+1), filepath.Base(file))
		if err != nil {
			t.Fatal(err)
		}

		f, err := os.Open(file)

		if err != nil {
			t.Fatal(err)
		}

		_, err = io.Copy(filePart, f)
		_ = f.Close()

		if err != nil {
			t.Fatal(err)
		}

	}
	err := w.Close()
	if err != nil {
		t.Fatal(err)
	}
	ctx, cancel := context.WithTimeout(context.Background(), 60*time.Second)

	defer cancel()

	req, err := http.NewRequestWithContext(ctx, http.MethodPost, "https://httpbin.org/post", reqBody)

	if err != nil {
		t.Fatal(err)
	}

	req.Header.Set("Content-Type", w.FormDataContentType())

	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		t.Fatal(err)
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	b, err := io.ReadAll(resp.Body)
	if err != nil {
		t.Fatal(err)
	}

	if resp.StatusCode != http.StatusOK {
		t.Fatalf("expected status %d; actual status %d", http.StatusOK, resp.StatusCode)
	}
	t.Logf("\n%s", b)

}
