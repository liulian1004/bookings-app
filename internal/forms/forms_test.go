package forms

import (
	"net/http/httptest"
	"net/url"
	"testing"
)

func TestForm_valid(t *testing.T) {
	//create a new request
	r := httptest.NewRequest("POST","/whatever", nil) 
	//create a new form
	form := New(r.PostForm)

	isValid := form.Valid()
	if !isValid{
		t.Error("go invalid form when should have been valid")
	}
}

func TestForm_Required(t *testing.T) {
	//create a new request
	r := httptest.NewRequest("POST","/whatever", nil) 
	//create a new form
	form := New(r.PostForm)
	form.Required("a","b","c")
	if form.Valid() {
		t.Error("form shows valid when required fields missing")
	} 
	postedData := url.Values{}
	postedData.Add("a","a")
	postedData.Add("b","a")
	postedData.Add("c","a")

	r = httptest.NewRequest("POST","/whatever", nil) 
	r.PostForm = postedData
	form = New(r.PostForm)
	form.Required("a","b","c")
	if !form.Valid() {
		t.Error("show does not have required fields when it does")
	} 
}

func TestForm_Has(t *testing.T) {
	//create a new request
	r := httptest.NewRequest("POST","/whatever", nil) 
	//create a new form
	form := New(r.PostForm)

	has := form.Has("whatever")

	if has {
		t.Error("forms shows has field when it does not has")
	}

	postedDate := url.Values{}
	postedDate.Add("a","a")
	form = New(postedDate)
	has = form.Has("a")

	if !has {
		t.Error("show form does not has when it should ")
	}
}

func TestForm_MinLength(t *testing.T) {
	//create a new request
	r := httptest.NewRequest("POST","/whatever", nil) 
	//create a new form
	form := New(r.PostForm)
	form.MinLength("x",10)

	if form.Valid() {
		t.Error("forms shows minLength for non-existing field")
	}
	isError := form.Errors.Get("x")
	if isError == "" {
		t.Error(("should has an error but not get one"))
	}

	postedValue := url.Values{}
	postedValue.Add("some_fields","some_value")
	form = New(postedValue)
	form.MinLength("some_fields",100)
	if form.Valid() {
		t.Error("forms shows minLength of 100 when data is shorter")
	}

	postedValue = url.Values{}
	postedValue.Add("another_field","abc")
	form = New(postedValue)
	form.MinLength("another_field",1)
	if !form.Valid() {
		t.Error("forms shows minLength of 1 when data is longer")
	}
	isError = form.Errors.Get("another_field")
	if isError != "" {
		t.Error(("should not have an error but  get one"))
	}
	

}

func TestForm_IsEmail(t *testing.T) {
	postedValue := url.Values{}
	form := New(postedValue)

	form.IsEmail("x")

	if form.Valid(){
		t.Error("form shows validated email for non existed filed")
	}

	postedValue = url.Values{}
	postedValue.Add("email","abc@abc.com")
	form = New(postedValue)
	form.IsEmail("email")
	if !form.Valid() {
		t.Error("got an invalided email when it has validated email filed")
	}

	postedValue = url.Values{}
	postedValue.Add("email","abc")
	form = New(postedValue)
	form.IsEmail("email")
	if form.Valid() {
		t.Error("got an validated email when it has invalidated email filed")
	}


}

