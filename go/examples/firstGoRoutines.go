
func getPlaylistIdByChannelIdOrCustomUrlAndPlaylistName(service *youtube.Service, part string, idType string, idValue string, playlistName string) {
  call := service.Channels.List(part)

  if idValue != "" {
    if idType=="customUrl" {
      call = call.ForUsername(idValue)
    } else if idType=="channelId"
      call = call.Id(idValue)
    }
  }

  response, err := call.Do()
  handleError(err, "")

  item = response.items[0]
  var playlistId string = ""
  if playlistName=="uploads" {
    playlistId = item.contentDetails.relatedPlaylists.uploads
  } else if playlistName=="favorites" {
    playlist = item.contentDetails.relatedPlaylists.favorites
  }
  fmt.Println(item.Id, ": ", playlist)
  }
}

/* ^^^^^^^^ da eingebaut
func channelsListById(service *youtube.Service, part string, id string) {
        call := service.Channels.List(part)
        if id != "" {
                call = call.Id(id)
        }
        response, err := call.Do()
        handleError(err, "")
        printChannelsListResults(response)
}
