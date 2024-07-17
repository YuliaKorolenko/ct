package ads

type AdFilter struct {
	Title     string
	AuthorID  int64
	Published bool
	Time      string
}

func NewFilter() AdFilter {
	return AdFilter{Title: "", AuthorID: -1, Published: true, Time: ""}
}

func (f *AdFilter) UpdateTitle(title string) {
	f.Title = title
}

func (f *AdFilter) UpdateAuthor(authorId int64) {
	f.AuthorID = authorId
}

func (f *AdFilter) UpdatePublished(published string) {
	f.Published = published == "true"
}

func (f *AdFilter) UpdateTime(time string) {
	f.Time = time
}
