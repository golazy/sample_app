package posts

import "testing"

func TestServiceListAndGet(t *testing.T) {
	service, err := New()
	if err != nil {
		t.Fatal(err)
	}

	first := service.List()
	second := service.List()
	if len(first) != 2 || first[0].Param != "hello-golazy" {
		t.Fatalf("unexpected posts: %#v", first)
	}
	first[0].Title = "changed"
	if second[0].Title == "changed" {
		t.Fatal("List returned shared slice storage")
	}

	post, ok := service.Get("embedded-content")
	if !ok || post.Title != "Embedded Content" {
		t.Fatalf("unexpected post: %#v, %v", post, ok)
	}
	if _, ok := service.Get("missing"); ok {
		t.Fatal("missing post was found")
	}
}
