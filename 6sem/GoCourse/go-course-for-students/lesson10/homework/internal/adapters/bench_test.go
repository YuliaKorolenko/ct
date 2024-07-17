package adapters

import (
	"github.com/stretchr/testify/assert"
	"homework10/internal/adapters/adrepo"
	"homework10/internal/adapters/usrepo"
	"homework10/internal/entities/ads"
	"homework10/internal/entities/user"
	"strings"
	"testing"
	"time"
)

var BenchSink int // make sure compiler cannot optimize away benchmarks

var (
	firstTitle  = strings.Repeat("t", 50)
	secondTitle = strings.Repeat("h", 50)
	thirdTitle  = strings.Repeat("hello", 11)
	fourthTitle = strings.Repeat("hello", 500)
	tests       = [...]ads.AdFilter{
		ads.CreateFilter(firstTitle, 0, false, time.Now().String()),
		ads.CreateFilter(secondTitle, 0, true, time.Now().String()),
		ads.CreateFilter(firstTitle, 0, true, time.Now().String()),
		ads.CreateFilter("hei", 0, true, time.Now().String()),
		ads.CreateFilter(thirdTitle, 0, true, time.Now().String()),
		ads.CreateFilter(secondTitle, 0, false, time.Now().String()),
		ads.CreateFilter(fourthTitle, 0, true, time.Now().String()),
		ads.CreateFilter(thirdTitle, 0, true, time.Now().String()),
	}
)

func BenchmarkAd(b *testing.B) {
	usrepo := usrepo.NewUserRepo()
	repo := adrepo.New()
	user, err := usrepo.CreateUser(user.NewUser(-1, "name", "email"))
	assert.NoError(b, err)
	for _, val := range tests {
		_, err := repo.CreateAd(ads.New(-1, val.Title, "text", user.ID))
		assert.NoError(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cur := i % len(tests)
		post, err := repo.GetAdsList(ads.CreateFilter(tests[cur].Title, user.ID, tests[cur].Published, tests[cur].Time))
		assert.NoError(b, err)
		BenchSink += len(post)
	}
}

func BenchmarkAdSlow(b *testing.B) {
	usrepo := usrepo.NewUserRepo()
	repo := adrepo.New()
	user, err := usrepo.CreateUser(user.NewUser(-1, "name", "email"))
	assert.NoError(b, err)
	for _, val := range tests {
		_, err := repo.CreateAd(ads.New(-1, val.Title, "text", user.ID))
		assert.NoError(b, err)
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		cur := i % len(tests)
		post, err := repo.GetAdsListSlow(ads.CreateFilter(tests[cur].Title, user.ID, tests[cur].Published, tests[cur].Time))
		assert.NoError(b, err)
		BenchSink += len(post)
	}
}
