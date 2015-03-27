package TwitterAPI

import (
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/antonholmquist/jason"
	"github.com/garyburd/go-oauth/oauth"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os/exec"
	"runtime"
	"strings"
)

//BaseUrl of all requests
const BASEURL = "https://api.twitter.com/1.1/"

//Main struct contains API key and secret
type Account struct {
	ConsumerKey    string
	ConsumerSecret string
}

var params = url.Values{}

type EndPoints struct {
	GetAccountSettings    string
	MentionsTimeline      string
	UserTimeline          string
	HomeTimeline          string
	RetweetsOfMe          string
	RetweetsByID          string
	ShowTweet             string
	Tweet                 string
	Retweet               string
	Oembed                string
	Favourite             string
	LookUp                string
	Retweeters            string
	MediaUpload           string
	ReportSpam            string
	DeleteTweet           string
	DMShow                string
	DMSent                string
	Search                string
	DirectMessages        string
	DMCreate              string
	DMDelete              string
	Following             string
	Followers             string
	PendingFollowersI     string
	PendingFollowersO     string
	FollowUser            string
	UnFollowUser          string
	FriendshipUpdate      string
	FriendshipShow        string
	FriendsList           string
	FollowersList         string
	FriendshipLookup      string
	VerifyCredentials     string
	ChangeAccountSettings string
	UpdateProfile         string
	UpdateBackgroundPic   string
	UpdatePicture         string
	BlockList             string
	BlockedIDs            string
	BlockUser             string
	UnBlockUser           string
	UsersLookup           string
	UsersShow             string
	UsersSearch           string
	UpdateBanner          string
	RemoveBanner          string
	GetUserBanner         string
	MuteUser              string
	UnmuteUser            string
	MutedIds              string
	MuteUserList          string
	FavoriteList          string
	UnFavorite            string
}

//Twitter Endpoints
var ENDPOINT = EndPoints{
	UnFavorite:            fmt.Sprintf("%favorites/destroy.json", BASEURL),
	FavoriteList:          fmt.Sprintf("%favorites/list.json", BASEURL),
	MuteUser:              fmt.Sprintf("%smutes/users/create.json", BASEURL),
	UnmuteUser:            fmt.Sprintf("%smutes/users/destroy.json", BASEURL),
	GetUserBanner:         fmt.Sprintf("%susers/profile_banner.json", BASEURL),
	RemoveBanner:          fmt.Sprintf("%saccount/remove_profile_banner.json", BASEURL),
	UpdateBanner:          fmt.Sprintf("%saccount/update_profile_banner.json", BASEURL),
	UsersSearch:           fmt.Sprintf("%susers/search.json", BASEURL),
	UsersShow:             fmt.Sprintf("%susers/show.json", BASEURL),
	UsersLookup:           fmt.Sprintf("%susers/lookup.json", BASEURL),
	UnBlockUser:           fmt.Sprintf("%sblocks/destroy.json", BASEURL),
	BlockUser:             fmt.Sprintf("%sblocks/create.json", BASEURL),
	BlockedIDs:            fmt.Sprintf("%sblocks/ids.json", BASEURL),
	BlockList:             fmt.Sprintf("%sblocks/list.json", BASEURL),
	UpdatePicture:         fmt.Sprintf("%saccount/update_profile_image.json", BASEURL),
	UpdateBackgroundPic:   fmt.Sprintf("%saccount/update_profile_background_image.json", BASEURL),
	UpdateProfile:         fmt.Sprintf("%saccount/update_profile.json", BASEURL),
	ChangeAccountSettings: fmt.Sprintf("%saccount/settings.json", BASEURL),
	VerifyCredentials:     fmt.Sprintf("%saccount/verify_credentials.json", BASEURL),
	FriendshipLookup:      fmt.Sprintf("%sfriendships/lookup.json", BASEURL),
	FriendsList:           fmt.Sprintf("%sfriends/list.json", BASEURL),
	FollowersList:         fmt.Sprintf("%sfollowers/list.json", BASEURL),
	FriendshipShow:        fmt.Sprintf("%sfriendships/show.json", BASEURL),
	FriendshipUpdate:      fmt.Sprintf("%sfriendships/update.json", BASEURL),
	UnFollowUser:          fmt.Sprintf("%sfriendships/destroy.json", BASEURL),
	FollowUser:            fmt.Sprintf("%sfriendships/create.json", BASEURL),
	PendingFollowersO:     fmt.Sprintf("%sfriendships/outgoing.json", BASEURL),
	PendingFollowersI:     fmt.Sprintf("%sfriendships/incoming.json", BASEURL),
	Followers:             fmt.Sprintf("%sfollowers/ids.json", BASEURL),
	Following:             fmt.Sprintf("%sfriends/ids.json", BASEURL),
	DMCreate:              fmt.Sprintf("%sdirect_messages/new.json", BASEURL),
	DMDelete:              fmt.Sprintf("%sdirect_messages/destroy.json", BASEURL),
	DirectMessages:        fmt.Sprintf("%sdirect_messages.json", BASEURL),
	DMShow:                fmt.Sprintf("%sdirect_messages/show.json", BASEURL),
	DMSent:                fmt.Sprintf("%sdirect_messages/sent.json", BASEURL),
	DeleteTweet:           fmt.Sprintf("%sstatuses/destroy/:id.json", BASEURL),
	ReportSpam:            fmt.Sprintf("%susers/report_spam.json", BASEURL),
	GetAccountSettings:    fmt.Sprintf("%saccount/settings.json", BASEURL),
	Favourite:             fmt.Sprintf("%sfavorites/create.json", BASEURL),
	MentionsTimeline:      fmt.Sprintf("%sstatuses/mentions_timeline.json", BASEURL),
	UserTimeline:          fmt.Sprintf("%sstatuses/user_timeline.json", BASEURL),
	HomeTimeline:          fmt.Sprintf("%sstatuses/home_timeline.json", BASEURL),
	RetweetsOfMe:          fmt.Sprintf("%sstatuses/retweets_of_me.json", BASEURL),
	RetweetsByID:          fmt.Sprintf("%sstatuses/retweets/:id.json", BASEURL),
	ShowTweet:             fmt.Sprintf("%sstatuses/show.json", BASEURL),
	Tweet:                 fmt.Sprintf("%sstatuses/update.json", BASEURL),
	Retweet:               fmt.Sprintf("%sstatuses/retweet/:id.json", BASEURL),
	Oembed:                fmt.Sprintf("%sstatuses/oembed.json", BASEURL),
	Retweeters:            fmt.Sprintf("%sstatuses/retweeters/ids.json", BASEURL),
	LookUp:                fmt.Sprintf("%sstatuses/lookup.json", BASEURL),
	MediaUpload:           fmt.Sprintf("%smedia/upload.json", BASEURL),
	Search:                fmt.Sprintf("%ssearch/tweets.json", BASEURL),
}

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

var token = &oauthClient.Credentials

func (P *Account) UnFavorite(ID string) (string, error) {

	var Params = params

	Params.Add("id", ID)

	resp, err := DoRequest(ENDPOINT.UnFavorite, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) FavoritesList(ScreenName, UserId, Count string) (string, error) {

	var Params = params

	switch {
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case Count != "":
		Params.Add("count", Count)
	case UserId != "":
		Params.Add("user_id", UserId)
	}

	resp, err := DoRequest(ENDPOINT.FavoriteList, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil

}

func (P *Account) UnMuteUser(ScreenName, UserId string) (string, error) {

	var Params = params

	switch {
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserId != "":
		Params.Add("user_id", UserId)
	}

	resp, err := DoRequest(ENDPOINT.UnmuteUser, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil

}

func (P *Account) MuteUser(ScreenName, UserId string) (string, error) {

	var Params = params

	switch {
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserId != "":
		Params.Add("user_id", UserId)
	}

	resp, err := DoRequest(ENDPOINT.MuteUser, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil

}

func (P *Account) GetUserBanner(ScreenName, UserId string) (string, error) {

	var Params = params

	switch {
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserId != "":
		Params.Add("user_id", UserId)
	}

	resp, err := DoRequest(ENDPOINT.GetUserBanner, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) RemoveBanner() (string, error) {

	resp, err := DoRequest(ENDPOINT.RemoveBanner, nil, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) UpdateBanner(Image, Width, Height, Offset_Left, Offset_Top string) (string, error) {

	var Params = params

	Image, _, err := ImageToBase64(Image)

	if err != nil {
		return "", err
	}

	switch {
	case Image == "":
		return "", errors.New("Image cannot be empty")
	case Width != "":
		Params.Add("width", Width)
	case Height != "":
		Params.Add("height", Height)
	case Offset_Left != "":
		Params.Add("offset_left", Offset_Left)
	case Offset_Top != "":
		Params.Add("offset_top", Offset_Top)
	}

	resp, err := DoRequest(ENDPOINT.UpdateBanner, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil

}

func (P *Account) UsersSearch(Q, Page string) (string, error) {
	var Params = params
	Params.Add("q", Q)
	Params.Add("page", Page)

	resp, err := DoRequest(ENDPOINT.UsersSearch, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) UsersShow(ScreenName, UserId string) (string, error) {

	var Params = params

	switch {
	case ScreenName == "" && UserId == "":
		return "", errors.New("ScreenName and UserId cannot both be empty")
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserId != "":
		Params.Add("user_id", UserId)
	}

	resp, err := DoRequest(ENDPOINT.UsersShow, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) UserLookUp(ScreenName, UserId string) (string, error) {

	var Params = params

	switch {
	case ScreenName == "" && UserId == "":
		return "", errors.New("ScreenName and UserId cannot both be empty")
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserId != "":
		Params.Add("user_id", UserId)
	}

	resp, err := DoRequest(ENDPOINT.UsersLookup, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) UnBlockUser(ScreenName, UserId string) (string, error) {

	var Params = params

	switch {
	case ScreenName == "" && UserId == "":
		return "", errors.New("ScreenName and UserId cannot both be empty")
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserId != "":
		Params.Add("user_id", UserId)
	}

	resp, err := DoRequest(ENDPOINT.UnBlockUser, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) BlockUser(ScreenName, UserId string) (string, error) {

	var Params = params

	switch {
	case ScreenName == "" && UserId == "":
		return "", errors.New("ScreenName and UserId cannot both be empty")
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserId != "":
		Params.Add("user_id", UserId)
	}

	resp, err := DoRequest(ENDPOINT.BlockUser, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil

}

// func (P *Account) BlockIds() (string, error) {

// }

//Only supports 5000 user objects currently (Need to sort out cursors)
func (P *Account) BlockList() (string, error) {

	resp, err := DoRequest(ENDPOINT.BlockList, nil, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) ChangeProfilePicture(FileName string) (string, error) {

	var Params = params

	F, _, err := ImageToBase64(FileName)

	if err != nil {
		return "", err
	}

	Params.Add("image", F)

	resp, err := DoRequest(ENDPOINT.UpdatePicture, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) RemoveBackgroundPicture() (string, error) {
	var Params = params

	Params.Add("use", "false")

	resp, err := DoRequest(ENDPOINT.UpdateBackgroundPic, Params, "POST")

	if err != nil {
		return "", nil
	}

	return resp, nil
}

func (P *Account) UpdateBackgroundPicture(FilePath string, Tile bool) (string, error) {

	var Params = params

	f, i, err := ImageToBase64(FilePath)

	if err != nil {
		return "", err
	}
	//if file size is greater than 800kb we return an error because twitter will not accept anything > 800kbs.
	if i > 800 {
		return "", errors.New("Images may not be largers than 800kbs")
	}

	fmt.Println(f)

	Params.Add("image", f)
	Params.Add("use", "1")

	if Tile == true {
		Params.Add("tile", "true")
	}

	resp, err := DoRequest(ENDPOINT.UpdateBackgroundPic, Params, "POST")

	if err != nil {
		return "", nil
	}

	return resp, nil
}
func (P *Account) UpdateProfile(Options map[string]string) (string, error) {

	Par := []string{"name", "url", "location", "description", "profile_link_color"}

	var Params = params

	Op := make(map[string]string)

	for key, value := range Options {
		for _, ele := range Par {
			if strings.Contains(key, ele) {
				Op[key] = value
			}
		}
	}

	for k, v := range Op {
		Params.Add(k, v)
	}
	resp, err := DoRequest(ENDPOINT.UpdatePicture, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) ChangeAccountSettings(Options map[string]string) (string, error) {

	var Params = params

	Pref := []string{"sleep_time_enabled", "trend_location_woeid", "start_sleep_time", "end_sleep_time", "time_zone", "lang"}

	Op := make(map[string]string)

	for key, value := range Options {
		for _, ele := range Pref {
			if strings.Contains(key, ele) {
				Op[key] = value
			}
		}
	}

	for k, v := range Op {
		Params.Add(k, v)
	}

	resp, err := DoRequest(ENDPOINT.ChangeAccountSettings, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) VerifyCredential() (string, error) {

	resp, err := DoRequest(ENDPOINT.VerifyCredentials, nil, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) FriendshipShow(ScreenName, TargetScreenName string) (string, error) {

	var Params = params

	Params.Add("source_screen_name", ScreenName)

	Params.Add("target_screen_name", TargetScreenName)

	resp, err := DoRequest(ENDPOINT.FriendshipShow, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) UnFollowUser(ScreenName, UserId string) (string, error) {

	var Params = params

	switch {

	case UserId == "" && ScreenName == "":
		return "", errors.New("UserID and ScreenName cannot both be empty")
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserId != "":
		Params.Add("user_id", UserId)
	}

	resp, err := DoRequest(ENDPOINT.UnFollowUser, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) FollowUser(ScreenName, UserId string) (string, error) {

	var Params = params

	switch {
	case UserId == "" && ScreenName == "":
		return "", errors.New("UserID and ScreenName cannot both be empty")
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserId != "":
		Params.Add("user_id", UserId)
	}

	resp, err := DoRequest(ENDPOINT.FollowUser, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil

}

func (P *Account) PendingFollowersOutgoing(Cursor string) (string, error) {

	var Params = params

	if Cursor != "" {
		Params.Add("cursor", Cursor)
	}

	resp, err := DoRequest(ENDPOINT.PendingFollowersO, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}
func (P *Account) PendingFollowersIncoming(Cursor string) (string, error) {

	var Params = params

	if Cursor != "" {

		Params.Add("cursor", Cursor)
	}

	resp, err := DoRequest(ENDPOINT.PendingFollowersI, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}
func (P *Account) FollowersList(UserID, ScreenName, Cursor, Count string) (string, error) {

	var Params = params

	switch {
	case UserID == "" && ScreenName == "":
		return "", errors.New("UserID and ScreenName cannot be both empty")
	case UserID != "":
		Params.Add("user_id", UserID)
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case Cursor != "":
		Params.Add("cursor", Cursor)
	case Count != "":
		Params.Add("count", Count)
	}

	resp, err := DoRequest(ENDPOINT.Followers, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil

}

func (P *Account) FollowingList(UserID, ScreenName, Cursor, Count string) (string, error) {

	var Params = params

	switch {
	case UserID == "" && ScreenName == "":
		return "", errors.New("UserID and ScreenName cannot be both empty")
	case UserID != "":
		Params.Add("user_id", UserID)
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case Cursor != "":
		Params.Add("cursor", Cursor)
	case Count != "":
		Params.Add("count", Count)
	}

	resp, err := DoRequest(ENDPOINT.Following, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) DMDelete(ID string) (string, error) {

	var Params = params

	Params.Add("id", ID)

	resp, err := DoRequest(ENDPOINT.DMDelete, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) DMCreate(UserID, ScreenName, Text string) (string, error) {

	var Params = params

	Params.Add("text", Text)

	switch {
	case UserID == "" && ScreenName == "":
		return "", errors.New("UserID and ScreenName cannot be both empty")
	case UserID != "":
		Params.Add("user_id", UserID)
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	}

	resp, err := DoRequest(ENDPOINT.DMCreate, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil

}

func (P *Account) DirectMessages(Count, SkipStatus string) (string, error) {
	var Params = params

	if Count != "" {
		Params.Add("count", Count)
	}

	resp, err := DoRequest(ENDPOINT.DirectMessages, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) Search(Query, GeoCode string) (string, error) {
	var Params = params

	Params.Add("q", Query)

	if GeoCode != "" {
		Params.Add("geocode", GeoCode)
	}

	resp, err := DoRequest(ENDPOINT.Search, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil

}

//UNTESTED
func (P *Account) DirectMessageShow(ID string) (string, error) {
	var Params = params
	Params.Add("id", ID)

	resp, err := DoRequest(ENDPOINT.DMShow, Params, "GET")

	if err != nil {
		return "", nil
	}

	return resp, nil
}

//UNTESTED
func (P *Account) DirectMessageSent(Page, Count string) (string, error) {

	var Params = params

	switch {
	case Page != "":
		Params.Add("page", Page)
	case Count != "":
		Params.Add("id", Count)
	}

	resp, err := DoRequest(ENDPOINT.DMSent, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil

}

//UNTESTED
func (P *Account) ReportForSpam(ID string) (string, error) {

	var Params = params

	Params.Add("id", ID)

	resp, err := DoRequest(ENDPOINT.ReportSpam, Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil

}

func (P *Account) DeleteTweet(ID string) (string, error) {

	var Params = params

	Params.Add("id", ID)

	resp, err := DoRequest(strings.Replace(ENDPOINT.DeleteTweet, ":id", ID, -1), Params, "POST")

	if err != nil {
		return "", err
	}

	return resp, nil

}

func (P *Account) Retweeters(ID string) (string, error) {
	//Cursor doesn't work?
	var Params = params

	Params.Add("id", ID)

	resp, err := DoRequest(ENDPOINT.Retweeters, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) RetweetsByID(ID, Count string) (string, error) {
	var Params = params
	Params.Add("id", ID)

	if Count != "" {
		Params.Add("count", Count)
	}

	resp, err := DoRequest(strings.Replace(ENDPOINT.RetweetsByID, ":id", ID, -1), Params, "GET")

	if err != nil {
		return "", err
	}
	return resp, nil
}
func (P *Account) RetweetsOfMe(Count string) (string, error) {
	var Params = params

	if Count != "" {
		Params.Add("count", Count)
	}

	resp, err := DoRequest(ENDPOINT.RetweetsOfMe, Params, "GET")

	if err != nil {
		return "", err
	}
	return resp, nil

}

func (P *Account) Oembed(ID, URL string) (string, error) {
	var Params = params

	switch {
	case ID == "" && URL == "":
		log.Fatal("RIP")
	case ID != "":
		Params.Add("id", ID)
	case URL != "":
		Params.Add("url", URL)
	}

	resp, err := DoRequest(ENDPOINT.Oembed, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func (P *Account) GetAccountSettings() (string, error) {

	resp, err := DoRequest(ENDPOINT.GetAccountSettings, nil, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil

}
func (P *Account) ShowTweet(ID string) (string, error) {
	var Params = params

	Params.Add("id", ID)

	resp, err := DoRequest(ENDPOINT.ShowTweet, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}
func (P *Account) LookUp(IDS []string) (string, error) {
	var Params = params

	ids := strings.Join(IDS, ",")

	Params.Add("id", ids)

	resp, err := DoRequest(ENDPOINT.LookUp, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}

func ImageToBase64(FilePath string) (string, int, error) {

	filedata, err := ioutil.ReadFile(FilePath)

	if err != nil {
		return "", 0, err
	}

	encoded := base64.StdEncoding.EncodeToString(filedata)

	return encoded, len(filedata) / 1000, nil
}
func (P *Account) MediaUpload(FilePath string, tweet bool) (string, error) {

	encoded, _, _ := ImageToBase64(FilePath)
	var Params = params

	Params.Add("media", encoded)

	resp, _ := oauthClient.Post(http.DefaultClient, token, "https://upload.twitter.com/1.1/media/upload.json", Params)
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	j, err := jason.NewObjectFromBytes([]byte(string(body)))

	if err != nil {
		return "", err
	}

	m, err := j.GetString("media_id_string")

	if err != nil {
		return "", err
	}

	if tweet {
		P.Tweet("", "", m, false, false)
	} else {
		return m, nil
	}

	return string(body), nil

}
func (P *Account) GetHomeTimeline(Count string) (string, error) {

	var Paramas = params

	if Count != "" {
		Paramas.Add("count", Count)
		resp, err := DoRequest(ENDPOINT.MentionsTimeline, Paramas, "GET")

		if err != nil {
			return "", err
		}

		return resp, nil
	}

	resp, err := DoRequest(ENDPOINT.MentionsTimeline, nil, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil

}

func (P *Account) GetMentionsTimeline(Count string) (string, error) {

	var Params = params

	if Count != "" {
		Params.Add("count", Count)
	}

	resp, err := DoRequest(ENDPOINT.MentionsTimeline, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil
}
func (P *Account) GetUserTimeline(ScreenName string, UserID string, Count string, IncludeRetweets bool) (string, error) {

	var Params = params

	if ScreenName == "" && UserID == "" {
		return "", errors.New("Screenname and UserID can't be both empty")
	}

	switch {
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserID != "":
		Params.Add("user_Id", UserID)
	}

	resp, err := DoRequest(ENDPOINT.UserTimeline, Params, "GET")

	if err != nil {
		return "", err
	}

	return resp, nil

}
func (P *Account) FavouriteTweet(TweetID string) (string, error) {

	var Params = params

	Params.Add("id", TweetID)

	resp, err := DoRequest(ENDPOINT.Favourite, Params, "POST")

	if err != nil {
		return "", err
	}
	return resp, nil

}

func (P *Account) UnAuth() {
	token = nil
}

func (P *Account) Retweet(TweetID string) (string, error) {

	var Params = params

	Params.Add("id", TweetID)

	resp, err := DoRequest(strings.Replace(ENDPOINT.Retweet, ":id", TweetID, -1), Params, "POST")

	if err != nil {
		return "", err
	}

	if strings.Contains(resp, "You have already retweeted this tweet.") {
		return "", errors.New("You can not retweet a tweet that is already retweeted")
	}

	return resp, nil

}

func (P *Account) Tweet(Status string, ReplyStatusID string, MediaId string, PossiblySenstive bool, DisplayCoordinates bool) (string, error) {
	var Params = params

	Params.Add("status", Status)
	switch {
	case Status == "" && MediaId == "":
		return "", errors.New("Status and MediaID cannot both be empty")
	case MediaId != "":
		Params.Add("media_ids", MediaId)
	case ReplyStatusID != "":
		Params.Add("in_reply_to_status_id", ReplyStatusID)
	case PossiblySenstive:
		Params.Add("possibly_sensitive", "true")
	case DisplayCoordinates:
		Params.Add("display_coordinates", "true")
	}

	resp, err := DoRequest(ENDPOINT.Tweet, Params, "POST")

	if err != nil {
		return "", err
	}

	if strings.Contains(resp, "errors") {
		return "", errors.New(resp)
	}
	return resp, nil
}

func (P *Account) Auth() (string, error) {

	oauthClient.Credentials.Token = P.ConsumerKey
	oauthClient.Credentials.Secret = P.ConsumerSecret

	tempcred, errors := oauthClient.RequestTemporaryCredentials(http.DefaultClient, "oob", nil)

	if errors != nil {
		return "", errors
	}

	test := oauthClient.AuthorizationURL(tempcred, nil)

	fmt.Printf("Paste the PIN code: ")

	switch runtime.GOOS {
	case "linux":
		exec.Command("xdg-open", test).Start()
	case "windows":
		exec.Command("cmd", "/c", "start", test).Start()
	case "darwin":
		exec.Command("open", test).Start()
	default:
		fmt.Println("Error opening the link try manually opening the link: ", test)
	}

	var code string
	fmt.Scanln(&code)

	tokenCred, _, err := oauthClient.RequestToken(http.DefaultClient, tempcred, code)

	if err != nil {
		return "", err
	}

	token = tokenCred

	return "", nil
}
func DoRequest(Endpoint string, Params url.Values, Method string) (string, error) {

	switch Method {
	case "POST":
		resp, err := oauthClient.Post(http.DefaultClient, token, Endpoint, Params)

		if err != nil {
			return "", err
		}

		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", err
		}

		return string(body), nil
	case "GET":
		resp, err := oauthClient.Get(http.DefaultClient, token, Endpoint, Params)

		if err != nil {
			return "", err
		}
		defer resp.Body.Close()

		body, err := ioutil.ReadAll(resp.Body)

		if err != nil {
			return "", err
		}

		return string(body), nil
	}

	return "", errors.New("You must supply either a GET or POST method.")
}

func (P *Account) TweetURLtoID(link string) string {

	a := strings.Split(link, "/")[5]

	return a

}
