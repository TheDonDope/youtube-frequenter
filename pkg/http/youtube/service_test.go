package youtube

import (
	"testing"

	"github.com/TheDonDope/youtube-frequenter/pkg/mocks"
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
		// TODO: rewrite
		// t.Errorf(errors.Fail("YouTubeService#ChannelsList", "want", "got"))
	}
}
