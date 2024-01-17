package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_Valid(t *testing.T) {
	r := httptest.NewRequest("POST", "/whatever", nil)
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid {
		t.Error("got invalid when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {

	// r := httptest.NewRequest("POST", "/whatever", nil)
	postedData := url.Values{}
	// form := New(r.PostForm)
	form := New(postedData)
	form.Required("a", "b", "c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	}
	//-------------------------------------------------------
	postedData = url.Values{}
	postedData.Add("a", "a")
	postedData.Add("b", "a")
	postedData.Add("c", "a")
	// r = httptest.NewRequest("POST", "/whatever", nil)
	// r.PostForm = postedData
	// form = New(r.PostForm)
	form = New(postedData)
	form.Required("a", "b", "c")
	if !form.Valid() {
		t.Error("shows does not have required fields when it does")
	}
}

func TestForm_Has(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	has := form.Has("whatever")
	if has {
		t.Error("form show has fieldd when it does not")
	}

	postedData = url.Values{}
	postedData.Add("a", "a")
	form = New(postedData)
	has = form.Has("a")
	if !has {
		t.Error("shows form does not have field when it should")
	}
}

func TestForm_MinLength(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.MinLength("x", 10)
	if form.Valid() {
		t.Error("form shows min length for non-existent field")
	}
	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error("should have an error but did not get one")
	}

	postedData = url.Values{}
	postedData.Add("someField", "some value")
	form = New(postedData)
	form.MinLength("someField", 100)
	if form.Valid() {
		t.Error("shows min length 100 met when data is shorter")
	}

	postedData = url.Values{}
	postedData.Add("other_field", "abc123")
	form = New(postedData)
	form.MinLength("other_field", 1)
	if !form.Valid() {
		t.Error("shows min length of 1 is not met when it does")
	}
	isError = form.Errors.Get("other_field")
	if isError != "" {
		t.Error("should NOT have an error but did got one")
	}
}

func TestForm_IsEmail(t *testing.T) {
	postedData := url.Values{}
	form := New(postedData)
	form.IsEmail("x")
	if form.Valid() {
		t.Error("form shows valid email for non-existent field")
	}

	postedData = url.Values{}
	postedData.Add("email", "me@at.com")
	form = New(postedData)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got and invalid email when we should not have")
	}

	postedData = url.Values{}
	postedData.Add("email", "x")
	form = New(postedData)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("got valid for invalid email")
	}
}
