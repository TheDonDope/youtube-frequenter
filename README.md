# Youtube Tinfoil Expose

Exfoliate!

## Example input data

### `channelID`

- pewdiepie: `UC-lHJZR3Gqxm24_Vd_AJ5Yw`

### `customURL`

- wwwKenFMde

### `playlistID`

- Breaking Bad Soundtrack: `PLkUhiWFEpuqoqQDugpSZE-sFiiHwz1FBU`

## Vocabulary

- `channelId`:
  - used in form [youtube.com/channel/UC52XYgEExV9VG6Rt-6vnzVA](https://www.youtube.com/channel/UC52XYgEExV9VG6Rt-6vnzVA)
- `customUrl`:
  - used in form [youtube.com/user/destiny](https://www.youtube.com/user/destiny)

## API Calls

-`getPlaylistIdByChannelIdOrCustomUrlAndPlaylistName`
  - Input from user: `<channelId>` or `<customUrl>` and `<playlistName>`
  - Inputparts:
    - `snippet` (for `title`)
    - `contentDetails` (for `relatedPlaylists` -> `uploads` Playlist (All videos the channel has upped))
  - Output:
    - Id of the requested playlist
  - Usage:
    - Get Id of Uploads and Favourites Playlist

- `getVideoIdsByPlaylistId`:
  - Input from user: `<playlistId>` and `<maxSearchResults>`
  - Inputparts:
    - `contentDetails` (for `videoId`)
  - Output:
    - Array of `videoId`s
  - Usage:
    - Get all videoIds of the Uploads and Favourites Playlist

- `getChannelIdsOfCommentersByVideoId`
  - Input from user: `<videoId>`
  - Inputparts: `snippet` (for `topLevelComment` -> `authorChannelId` -> `value`)
  - Output:
    - Array of `channelId`s of the video commenters
  - Usage
    - Get all channel ids of people who commented on a video

## Utility Functions

- `getStatisticsByChannelIdOrCustomUrl`:
  - Input from user: `<channelId>` or `<customUrl>`
  - Inputparts:
    - `statistics`
  - Output:

## Workflow Example

1. Get PlaylistId of Uploads Playlist
2. Get all VideoIds of uploaded videos
3. Loop through all videos and comments to get channelIds of commentors
4. Loop through all channelIds and get the id of the favourites playlist
5. Loop through all videos of the favourites playlist
6. Get the channelIds/Name of the favourited video
