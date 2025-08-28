package types

type SnekFile struct {
	CourseSlug    string           `json:"course_slug"`
	UpdateTime    int64            `json:"update_time"`
	CourseContent CheckoutResponse `json:"course_content"`
}

type CheckoutResponse []struct {
	Name string `json:"name"`
	Url  string `json:"url"`
	Hash string `json:"hash"`
}
