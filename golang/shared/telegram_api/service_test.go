package telegram_api

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestGetMessageView(t *testing.T) {
	s := New()
	views, err := s.GetMessageViews("@pets132213", 12)
	assert.Nil(t, err)
	fmt.Println("### views", views)
}
