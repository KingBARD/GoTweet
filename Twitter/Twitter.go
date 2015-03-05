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
	GetAccountSettings string
	MentionsTimeline   string
	UserTimeline       string
	HomeTimeline       string
	RetweetsOfMe       string
	RetweetsByID       string
	ShowTweet          string
	Tweet              string
	Retweet            string
	Oembed             string
	Favourite          string
	LookUp             string
	Retweeters         string
	MediaUpload        string
	ReportSpam         string
	DeleteTweet        string
	DMShow             string
	DMSent             string
	Search             string
}

var ENDPOINT = EndPoints{
	DMShow:             fmt.Sprintf("%sdirect_messages/show.json", BASEURL),
	DMSent:             fmt.Sprintf("%sdirect_messages/sent.json", BASEURL),
	DeleteTweet:        fmt.Sprintf("%sstatuses/destroy/:id.json", BASEURL),
	ReportSpam:         fmt.Sprintf("%susers/report_spam.json", BASEURL),
	GetAccountSettings: fmt.Sprintf("%saccount/settings.json", BASEURL),
	Favourite:          fmt.Sprintf("%sfavorites/create.json", BASEURL),
	MentionsTimeline:   fmt.Sprintf("%sstatuses/mentions_timeline.json", BASEURL),
	UserTimeline:       fmt.Sprintf("%sstatuses/user_timeline.json", BASEURL),
	HomeTimeline:       fmt.Sprintf("%sstatuses/home_timeline.json", BASEURL),
	RetweetsOfMe:       fmt.Sprintf("%sstatuses/retweets_of_me.json", BASEURL),
	RetweetsByID:       fmt.Sprintf("%sstatuses/retweets/:id.json", BASEURL),
	ShowTweet:          fmt.Sprintf("%sstatuses/show.json", BASEURL),
	Tweet:              fmt.Sprintf("%sstatuses/update.json", BASEURL),
	Retweet:            fmt.Sprintf("%sstatuses/retweet/:id.json", BASEURL),
	Oembed:             fmt.Sprintf("%sstatuses/oembed.json", BASEURL),
	Retweeters:         fmt.Sprintf("%sstatuses/retweeters/ids.json", BASEURL),
	LookUp:             fmt.Sprintf("%sstatuses/lookup.json", BASEURL),
	MediaUpload:        fmt.Sprintf("%smedia/upload.json", BASEURL),
	Search:             fmt.Sprintf("%ssearch/tweets.json", BASEURL),
}

var oauthClient = oauth.Client{
	TemporaryCredentialRequestURI: "https://api.twitter.com/oauth/request_token",
	ResourceOwnerAuthorizationURI: "https://api.twitter.com/oauth/authenticate",
	TokenRequestURI:               "https://api.twitter.com/oauth/access_token",
}

var token = &oauthClient.Credentials

//UNTESTED
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

	resp, _ := DoRequest(ENDPOINT.GetAccountSettings, nil, "GET")

	return resp

}
func (P *Account) ShowTweet(ID string) (string, error) {
	var Params = params

	Params.Add("id", ID)

	resp, _ := DoRequest(ENDPOINT.ShowTweet, Params, "GET")

	return resp, nil
}
func (P *Account) LookUp(IDS []string) (string, error) {
	var Params = params

	ids := strings.Join(IDS, ",")

	Params.Add("id", ids)

	resp, _ := DoRequest(ENDPOINT.LookUp, Params, "GET")

	return resp, nil
}
func (P *Account) MediaUpload(FilePath string, tweet bool) (string, error) {
	filedata, _ := ioutil.ReadFile(FilePath)

	encoded := base64.StdEncoding.EncodeToString(filedata)
	var Params = params
	Params.Add("media", encoded)

	resp, _ := oauthClient.Post(http.DefaultClient, token, "https://upload.twitter.com/1.1/media/upload.json", Params)
	defer resp.Body.Close()

	body, errs := ioutil.ReadAll(resp.Body)
	if errs != nil {
		log.Fatal(errs)
	}
	j, err := jason.NewObjectFromBytes([]byte(string(body)))
	if err != nil {
		log.Fatal(err)
	}
	m, errs := j.GetString("media_id_string")
	if errs != nil {
		log.Fatal(errs)
	}
	//Json decode response get media id then parse it to the Tweet function below
	if tweet {
		P.Tweet("", "", m, false, false)
	} else {
		return m
	}

	return string(body)

}
func (P *Account) GetHomeTimeline(Count string) (string, error) {

	var Paramas = params

	if Count != "" {
		Paramas.Add("count", Count)
		resp, _ := DoRequest(ENDPOINT.MentionsTimeline, Paramas, "GET")
		return resp
	}

	resp, _ := DoRequest(ENDPOINT.MentionsTimeline, nil, "GET")

	return resp

}

func (P *Account) GetMentionsTimeline(Count string) (string, error) {
	var Params = params
	if Count != "" {
		Params.Add("count", Count)
	}

	resp, _ := DoRequest(ENDPOINT.MentionsTimeline, Params, "GET")

	return resp
}
func (P *Account) GetUserTimeline(ScreenName string, UserID string, Count string, IncludeRetweets bool) (string, error) {

	var Params = params

	if ScreenName == "" && UserID == "" {
		log.Fatal("Nigga what the fuck thinking??/??")
	}
	switch {
	case ScreenName != "":
		Params.Add("screen_name", ScreenName)
	case UserID != "":
		Params.Add("user_Id", UserID)
	}

	resp, _ := DoRequest(ENDPOINT.UserTimeline, Params, "GET")

	return resp

}
func (P *Account) FavouriteTweet(TweetID string) (string, error) {
	var Params = params

	Params.Add("id", TweetID)

	resp, _ := DoRequest(ENDPOINT.Favourite, Params, "POST")
	return resp

}

func (P *Account) UnAuth() {
	token = nil
}

func (P *Account) Retweet(TweetID string) (string, error) {

	var Params = params

	Params.Add("id", TweetID)

	resp, _ := DoRequest(strings.Replace(ENDPOINT.Retweet, ":id", TweetID, -1), Params, "POST")
	if strings.Contains(resp, "You have already retweeted this tweet.") {
		//change this to return new error etc
		log.Fatal("You can not retweet a tweet that is already retweeted")
	}

	return resp

}
func (P *Account) Tweet(Status string, ReplyStatusID string, MediaId string, PossiblySenstive bool, DisplayCoordinates bool) (string, error) {
	var Params = params

	Params.Add("status", Status)
	switch {
	case Status == "" && MediaId == "":
		return "", errors.New("Status cannot be empty")
	case MediaId != "":
		Params.Add("media_ids", MediaId)
	case ReplyStatusID != "":
		Params.Add("in_reply_to_status_id", ReplyStatusID)
	case PossiblySenstive:
		Params.Add("possibly_sensitive", "true")
	case DisplayCoordinates:
		Params.Add("display_coordinates", "true")
	}

	resp, _ := DoRequest(ENDPOINT.Tweet, Params, "POST")

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
		fmt.Println("...")
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
