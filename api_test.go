package gotau

import (
	"fmt"
	"github.com/stretchr/testify/require"
	"net"
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestClient_GetStreamersReturnsStreamers(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/streamers/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "[{\"id\":\"bcd0f3a5-9db9-46eb-8fb3-374ecabace47\",\"twitch_username\":\"wwsean08\",\"twitch_id\":null,\"streaming\":false,\"disabled\":false,\"created\":\"2021-05-19T20:53:31+0000\",\"updated\":\"2021-06-19T20:00:28+0000\"},{\"id\":\"5d8be520-9883-4d09-821a-3c71723e4880\",\"twitch_username\":\"GeekyCleanGaming\",\"twitch_id\":\"174097893\",\"streaming\":false,\"disabled\":false,\"created\":\"2021-05-22T22:22:51+0000\",\"updated\":\"2021-06-19T20:00:28+0000\"}]")
		require.NoError(t, err)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client := Client{
		hostname: host,
		port:     portNum,
		token:    "foo",
		hasSSL:   false,
	}

	streamers, err := client.GetStreamers()
	require.NoError(t, err)
	require.Len(t, streamers, 2)
	sean := streamers[0]
	geeky := streamers[1]

	require.NotNil(t, sean)
	require.NotNil(t, geeky)

	require.Equal(t, "bcd0f3a5-9db9-46eb-8fb3-374ecabace47", sean.ID)
	require.Equal(t, "wwsean08", sean.TwitchUsername)
	require.Zero(t, sean.TwitchID)
	require.False(t, sean.Streaming)
	require.False(t, sean.Disabled)
	require.Equal(t, 2021, sean.Created.Year())
	require.Equal(t, time.Month(5), sean.Created.Month())
	require.Equal(t, 2021, sean.Updated.Year())
	require.Equal(t, time.Month(6), sean.Updated.Month())

	require.Equal(t, "5d8be520-9883-4d09-821a-3c71723e4880", geeky.ID)
	require.Equal(t, "GeekyCleanGaming", geeky.TwitchUsername)
	require.Equal(t, "174097893", geeky.TwitchID)
	require.False(t, geeky.Streaming)
	require.False(t, geeky.Disabled)
	require.Equal(t, 2021, geeky.Created.Year())
	require.Equal(t, time.Month(5), geeky.Created.Month())
	require.Equal(t, 2021, geeky.Updated.Year())
	require.Equal(t, time.Month(6), geeky.Updated.Month())
}

func TestClient_GetStreamersReturnsError(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/streamers/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusInternalServerError)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client := Client{
		hostname: host,
		port:     portNum,
		token:    "foo",
		hasSSL:   false,
	}

	streamers, err := client.GetStreamers()
	require.Error(t, err)
	require.Nil(t, streamers)
	require.IsType(t, GenericError{}, err)
}

func TestClient_GetLatestStreamForStreamer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/streamers/5d8be520-9883-4d09-821a-3c71723e4880/streams/latest/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"id\":\"a409e338-3ff4-4fba-a0aa-5b4fd9d67381\",\"stream_id\":\"42219045165\",\"user_id\":\"174097893\",\"user_login\":\"geekycleangaming\",\"user_name\":\"GeekyCleanGaming\",\"game_id\":\"511224\",\"game_name\":\"Apex Legends\",\"type\":\"live\",\"title\":\"Grinding Ranked For Funsies 🔫 !support !merch !discord\",\"viewer_count\":0,\"started_at\":\"2021-06-05T01:09:40+0000\",\"ended_at\":\"2021-06-19T20:00:28+0000\",\"language\":\"en\",\"thumbnail_url\":\"https://static-cdn.jtvnw.net/previews-ttv/live_user_geekycleangaming-{width}x{height}.jpg\",\"tag_ids\":\"['6ea6bca4-4712-4ab9-a906-e3336a9d8039']\",\"is_mature\":false}")
		require.NoError(t, err)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client := Client{
		hostname: host,
		port:     portNum,
		token:    "foo",
		hasSSL:   false,
	}

	stream, err := client.GetLatestStreamForStreamer("5d8be520-9883-4d09-821a-3c71723e4880")
	require.NoError(t, err)
	require.NotNil(t, stream)
	require.Equal(t, "a409e338-3ff4-4fba-a0aa-5b4fd9d67381", stream.ID)
	require.Equal(t, "42219045165", stream.StreamID)
	require.Equal(t, "174097893", stream.UserID)
	require.Equal(t, "geekycleangaming", stream.UserLogin)
	require.Equal(t, "GeekyCleanGaming", stream.UserName)
	require.Equal(t, "511224", stream.GameID)
	require.Equal(t, "Apex Legends", stream.GameName)
	require.Equal(t, "live", stream.Type)
	require.Equal(t, "Grinding Ranked For Funsies 🔫 !support !merch !discord", stream.Title)
	require.Zero(t, stream.ViewerCount)
	require.Equal(t, 2021, stream.StartedAt.Year())
	require.Equal(t, time.Month(6), stream.StartedAt.Month())
	require.Equal(t, 5, stream.StartedAt.Day())
	require.Equal(t, 2021, stream.EndedAt.Year())
	require.Equal(t, time.Month(6), stream.EndedAt.Month())
	require.Equal(t, 19, stream.EndedAt.Day())
	require.Equal(t, "en", stream.Language)
	require.Equal(t, "https://static-cdn.jtvnw.net/previews-ttv/live_user_geekycleangaming-{width}x{height}.jpg", stream.ThumbnailUrl)
	require.False(t, stream.IsMature)
	require.Len(t, stream.TagIDs, 1)
	require.Equal(t, "6ea6bca4-4712-4ab9-a906-e3336a9d8039", stream.TagIDs[0])
}

func TestClient_GetLatestStreamForStreamerReturnsError(t *testing.T) {
	c := Client{}
	stream, err := c.GetLatestStreamForStreamer("")
	require.ErrorIs(t, err, BadRequestError{"invalid request, ID can't be blank"})
	require.Nil(t, stream)

	stream, err = c.GetLatestStreamForStreamer("    ")
	require.ErrorIs(t, err, BadRequestError{"invalid request, ID can't be blank"})
	require.Nil(t, stream)

	stream, err = c.GetLatestStreamForStreamer("	")
	require.ErrorIs(t, err, BadRequestError{"invalid request, ID can't be blank"})
	require.Nil(t, stream)
}

func TestClient_FollowStreamerOnTau(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/streamers/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "POST", r.Method)
		w.WriteHeader(http.StatusCreated)
		_, err := fmt.Fprint(w, "{\"id\":\"417e786c-2e48-4371-97bc-e782ab44f524\",\"twitch_username\":\"Freyline\",\"twitch_id\":\"208887405\",\"streaming\":false,\"disabled\":false,\"created\":\"2021-06-23T00:35:41+0000\",\"updated\":\"2021-06-23T00:35:41+0000\"}")
		require.NoError(t, err)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client := Client{
		hostname: host,
		port:     portNum,
		token:    "foo",
		hasSSL:   false,
	}

	streamer, err := client.FollowStreamerOnTau("Freyline")
	require.NoError(t, err)
	require.NotNil(t, streamer)
	require.Equal(t, "417e786c-2e48-4371-97bc-e782ab44f524", streamer.ID)
	require.Equal(t, "Freyline", streamer.TwitchUsername)
	require.Equal(t, "208887405", streamer.TwitchID)
	require.False(t, streamer.Streaming)
	require.False(t, streamer.Disabled)
	require.Equal(t, 2021, streamer.Created.Year())
	require.Equal(t, 2021, streamer.Updated.Year())
}

func TestClient_FollowStreamerOnTauReturnsError(t *testing.T) {
	c := Client{}
	streamer, err := c.FollowStreamerOnTau("")
	require.ErrorIs(t, err, BadRequestError{"invalid request, username can't be blank"})
	require.Nil(t, streamer)

	streamer, err = c.FollowStreamerOnTau("    ")
	require.ErrorIs(t, err, BadRequestError{"invalid request, username can't be blank"})
	require.Nil(t, streamer)

	streamer, err = c.FollowStreamerOnTau("	")
	require.ErrorIs(t, err, BadRequestError{"invalid request, username can't be blank"})
	require.Nil(t, streamer)
}

func TestClient_GetStreamsForStreamer(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		require.Equal(t, "/api/v1/streamers/34593db1-1228-40e0-bc3b-2c14b58b1f64/streams/", r.URL.Path)
		require.Equal(t, "Token foo", r.Header.Get("Authorization"))
		require.Equal(t, "GET", r.Method)
		w.WriteHeader(http.StatusOK)
		_, err := fmt.Fprint(w, "{\"count\":1,\"next\":null,\"previous\":null,\"results\":[{\"id\":\"dd780465-6454-40ed-98a6-c863a9638d09\",\"stream_id\":\"42502956941\",\"user_id\":\"501585826\",\"user_login\":\"4davidblue\",\"user_name\":\"4davidblue\",\"game_id\":\"509658\",\"game_name\":\"Just Chatting\",\"type\":\"live\",\"title\":\"IT TAKES TWOsday (w CiaoJordyn) | !donate !gif !freesub !dnd\",\"viewer_count\":0,\"started_at\":\"2021-06-23T01:31:54+0000\",\"ended_at\":null,\"language\":\"en\",\"thumbnail_url\":\"https://static-cdn.jtvnw.net/previews-ttv/live_user_4davidblue-{width}x{height}.jpg\",\"tag_ids\":null,\"is_mature\":false}]}")
		require.NoError(t, err)
	}))
	defer ts.Close()

	url := strings.TrimPrefix(ts.URL, "http://")
	host, port, err := net.SplitHostPort(url)
	require.NoError(t, err)
	portNum, err := strconv.Atoi(port)
	require.NoError(t, err)

	client := Client{
		hostname: host,
		port:     portNum,
		token:    "foo",
		hasSSL:   false,
	}

	streams, err := client.GetStreamsForStreamer("34593db1-1228-40e0-bc3b-2c14b58b1f64", 10)
	require.NoError(t, err)
	require.NotNil(t, streams)
	require.Len(t, streams, 1)
	require.Equal(t, "dd780465-6454-40ed-98a6-c863a9638d09", streams[0].ID)
	require.Equal(t, "42502956941", streams[0].StreamID)
	require.Equal(t, "501585826", streams[0].UserID)
	require.Equal(t, "4davidblue", streams[0].UserLogin)
	require.Equal(t, "4davidblue", streams[0].UserName)
	require.Equal(t, "509658", streams[0].GameID)
	require.Equal(t, "Just Chatting", streams[0].GameName)
	require.Equal(t, "live", streams[0].Type)
	require.Equal(t, "IT TAKES TWOsday (w CiaoJordyn) | !donate !gif !freesub !dnd", streams[0].Title)
	require.Zero(t, streams[0].ViewerCount)
	require.Equal(t, 2021, streams[0].StartedAt.Year())
	require.Nil(t, streams[0].EndedAt)
	require.Equal(t, "en", streams[0].Language)
	require.Equal(t, "https://static-cdn.jtvnw.net/previews-ttv/live_user_4davidblue-{width}x{height}.jpg", streams[0].ThumbnailUrl)
	require.Nil(t, streams[0].TagIDs)
	require.False(t, streams[0].IsMature)
}

func TestClient_GetStreamsForStreamerReturnsError(t *testing.T) {
	c := Client{}
	streams, err := c.GetStreamsForStreamer("", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, streamer id can't be blank"})
	require.Nil(t, streams)

	streams, err = c.GetStreamsForStreamer("    ", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, streamer id can't be blank"})
	require.Nil(t, streams)

	streams, err = c.GetStreamsForStreamer("	", 0)
	require.ErrorIs(t, err, BadRequestError{"invalid request, streamer id can't be blank"})
	require.Nil(t, streams)
}
