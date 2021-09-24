package result

import (
	"fmt"
	"github.com/PeterYangs/article-spider/v2/form"
)

type result struct {
	form *form.Form
}

func NewResult(form *form.Form) *result {

	return &result{form: form}
}

func (r *result) Work() {

	for s := range r.form.Storage {

		fmt.Println(s)
	}

}
