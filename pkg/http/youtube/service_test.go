package youtube

import (
	"testing"

	"gitlab.com/TheDonDope/gocha/v3/pkg/errors"
	"gitlab.com/TheDonDope/youtube-frequenter/pkg/mocks"
	youtubeV3 "google.golang.org/api/youtube/v3"
)

func TestChannelsList(t *testing.T) {
	s := mocks.YouTubeService{}
	want := &youtubeV3.ChannelListResponse{}
	got, err := s.ChannelsList("nil", "nil")
	if err != nil {
		t.Fatal(err)
	}
	if got.Etag != want.Etag {
		t.Errorf(errors.Fail("YouTubeService#ChannelsList", "want", "got"))
	}
}
