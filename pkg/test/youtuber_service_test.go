package test

import (
	"testing"

	ourErrors "gitlab.com/TheDonDope/youtube-frequenter/pkg/util/errors"
	"google.golang.org/api/youtube/v3"
)

func TestChannelsList(t *testing.T) {
	service := YouTuberMock{}

	expectedResponse := &youtube.ChannelListResponse{}

	actualResponse, actualResponseError := service.ChannelsList(nil, "nil", "nil")
	if actualResponseError != nil {
		t.Fatal(actualResponseError)
	}
	if actualResponse.Etag != expectedResponse.Etag {
		t.Errorf(ourErrors.GetFormattedFailMessage("YouTuber#ChannelsList", "expectedResponse", "actualResponse"))
	}
}
