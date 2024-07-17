package ads

import (
	"time"
)

type Ad struct {
	ID        int64
	Title     string `validate:"min:1;max:100"`
	Text      string `validate:"min:1;max:100"`
	AuthorID  int64
	Published bool
	Time      time.Time
}

func New(id int64, title string, text string, authorId int64) Ad {
	return Ad{ID: id, Title: title, Text: text, AuthorID: authorId, Published: false, Time: time.Now().UTC()}
}

func Update(id int64, title string, text string, authorId int64, published bool, timeCreate time.Time) Ad {
	return Ad{ID: id, Title: title, Text: text, AuthorID: authorId, Published: published, Time: timeCreate}
}
