// the is a rework for easy understanding & use, not built from scratch!
package main

import (
	"context"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/machinebox/progress"
	"golang.org/x/net/html"
)

// StreamUrls download link response struct
type StreamUrls struct {
	Video []Video
}

// Response download link response struct
type Response struct {
	AssetType  string     `json:"asset_type"`
	StreamUrls StreamUrls `json:"stream_urls"`
}

// Video videos response struct
type Video struct {
	File, Type, Label string
}

// CourseResponse videos response struct
type CourseResponse struct {
	Results []Course
}

// CourseDetail videos response struct
type CourseDetail struct {
	Results []CourseContent
}

// CourseContent detail
type CourseContent struct {
	Class       string `json:"_class"`
	ID          int
	Title       string
	Asset       Asset
	ObjectIndex int `json:"object_index"`
}

// Asset detail
type Asset struct {
	Class               string `json:"_class"`
	ID                  int
	AssetType           string `json:"asset_type"`
	Filename            string
	SupplementaryAssets []SupplementaryAssets `json:"supplementary_assets"`
	Title               string
	ObjectIndex         int
}

// SupplementaryAssets detail
type SupplementaryAssets struct {
	Class     string `json:"_class"`
	ID        int
	AssetType string `json:"asset_type"`
	Filename  string
}

// Course videso response struct
type Course struct {
	ID         int
	Title, URL string
}

// Udemy struct
type Udemy struct {
	AccessToken       string
	CourseURL         string
	SelectedCourseID  string
	Start             int
	End               int
	Resolution        string
	SessionMaxAttempt int
	CurrentAttempt    int
	DownloadPath      string
}

// Udemy URLs
const (
	GetCoursesURL      = "https://www.udemy.com/api-2.0/users/me/subscribed-courses/?ordering=-last_accessed&fields[course]=@min,title,id&page=1&page_size=100"
	GetDownloadURL     = "https://www.udemy.com/api-2.0/assets/{{assetID}}?fields[asset]=@min,status,asset_type,time_estimation,stream_urls&fields"
	GetCourseDetailURL = "https://www.udemy.com/api-2.0/courses/{{courseID}}/subscriber-curriculum-items/?page_size=1400&fields[lecture]=title,object_index,asset,supplementary_assets&fields[chapter]=title,object_index&fields[asset]=filename,asset_type&caching_intent=True"
)

func main() {
	// Prints the package info message
	Info()

	accessToken := flag.String("access-token", "false", "Authentication Token")
	CourseURL := flag.String("course-url", "false", "Course URL")
	Start := flag.Int("start", 0, "Start Lecture Id")
	End := flag.Int("end", 0, "End Lecture Id")
	Resolution := flag.String("resolution", "false", "Video Resolution")
	DownloadPath := flag.String("download-location", "false", "Download Path")
	flag.Parse()

	u := Udemy{
		AccessToken:       "Bearer " + *accessToken,
		CourseURL:         *CourseURL,
		SelectedCourseID:  "false",
		Start:             *Start,
		End:               *End,
		Resolution:        *Resolution,
		SessionMaxAttempt: 3,
		CurrentAttempt:    0,
		DownloadPath:      *DownloadPath,
	}

	_, err := u.AuthenticateToken()

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	_, err = u.GetCourses()

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	courseAssests, err := u.GetCourseDetail()

	if err != nil {
		fmt.Println(err)
		os.Exit(0)
	}

	u.startDownloading(courseAssests)

	// Wait for user input before exiting the program
	fmt.Println("✅  All downloads completed. Press Enter to exit.")
	fmt.Scanln()
}

// Info Package info message
func Info() {
	info := `
                                    ...
	                   //...the revive projekt!
	     *****************************************************
	     **                                                 **
	     **        ██╗░░░██╗██████╗░██╗░░░░░██╗░░██╗        **
	     **        ██║░░░██║██╔══██╗██║░░░░░╚██╗██╔╝        **
	     **        ██║░░░██║██║░░██║██║░░░░░░╚███╔╝░        **
	     **        ██║░░░██║██║░░██║██║░░░░░░██╔██╗░        **
	     **        ╚██████╔╝██████╔╝███████╗██╔╝╚██╗        **
	     **        ░╚═════╝░╚═════╝░╚══════╝╚═╝░░╚═╝        **
	     **       Yet Another Udemy Course Downloader       **
	     *****************************************************
                            Name : UdlX
                         Version : 1.0b
                          Author : mr•vybes (iykex)
                          Github : https://github.com/iykex
  `
	fmt.Println(info)
}

// BytesToMegaBytes convert bytes to mb
func BytesToMegaBytes(n int64) float64 {
	mb := float64(n / (1000 * 1024))
	return math.Floor(mb*100) / 100
}

// NewRequest to create new request for udemy
func (u Udemy) NewRequest(method, url string) *http.Response {
	client := &http.Client{}
	req, err := http.NewRequest(method, url, nil)

	if err != nil {
		fmt.Println(err)
	}

	req.Header.Add("Authorization", u.AccessToken)
	res, err := client.Do(req)

	return res
}

// AuthenticateToken Authnticate user provided token
func (u *Udemy) AuthenticateToken() (bool, error) {

	if u.AccessToken == "Bearer false" {
		u.getAuthenticationToken()
	}

	fmt.Println("✳️   : Authenticating Access Token...")
	resp := u.NewRequest("HEAD", GetCoursesURL)
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return false, errors.New("❌  : Invalid Authentication Token")
	}

	fmt.Println("➕  : Succesfully Authenticated ✅")
	return true, nil
}

// GetCourses to get all the courses details
func (u Udemy) GetCourses() (bool, error) {
	// Return if CourseURL is present
	if u.CourseURL != "false" {
		return false, nil
	}

	fmt.Println("✳️   : Fetching courses...")

	resp := u.NewRequest("GET", GetCoursesURL)
	defer resp.Body.Close()

	if resp.StatusCode > 299 {
		return false, errors.New("❌  : Error fetching courses, try to open link in your browser \n" + GetCoursesURL)
	}

	body, _ := ioutil.ReadAll(resp.Body)

	// fmt.Println(body)
	var response CourseResponse
	json.Unmarshal(body, &response)

	fmt.Println("➕  : Courses [ID]")

	for i := range response.Results {
		fmt.Printf("   ID[%v] : %v \n", response.Results[i].ID, response.Results[i].Title)
	}

	return true, nil
}

// GetCourseDetail to get course detail
func (u *Udemy) GetCourseDetail() ([]Asset, error) {

	if u.CourseURL != "false" {
		u.ParseHTMLAndGetCourseID()
	}

	if u.SelectedCourseID == "false" {
		u.getCourseID()
	}

	fmt.Println("✳️  : Fetching course lectures...")

	url := strings.Replace(GetCourseDetailURL, "{{courseID}}", u.SelectedCourseID, 1)
	res := u.NewRequest("GET", url)
	defer res.Body.Close()

	if res.StatusCode == 404 {
		fmt.Println("❌  : Invalid course ID")
		u.SelectedCourseID = "false"

		return u.GetCourseDetail()
	}

	if res.StatusCode > 299 {
		return nil, errors.New("❌  : Error fetching course lectures, try to open link in your browser \n" + url)
	}

	body, _ := ioutil.ReadAll(res.Body)

	var response CourseDetail
	json.Unmarshal(body, &response)

	finalAsset := make([]Asset, len(response.Results))

	fmt.Println("➕  : Lectures [No.]")

	for i := range response.Results {
		if response.Results[i].Class == "lecture" && (response.Results[i].Asset.AssetType == "Video" || response.Results[i].Asset.AssetType == "File") {
			fmt.Printf("   No.[%v] : %v[%v] \n", response.Results[i].ObjectIndex, response.Results[i].Title, response.Results[i].Asset.AssetType)
			response.Results[i].Asset.Title = response.Results[i].Title
			response.Results[i].Asset.ObjectIndex = response.Results[i].ObjectIndex
			finalAsset[response.Results[i].ObjectIndex] = response.Results[i].Asset
		}
	}

	return finalAsset, nil
}

// getCourseID get course id from user input
func (u *Udemy) getCourseID() {
	fmt.Print("❓  : Enter the course ID which you want to download: ")
	var courseID string
	fmt.Scanln(&courseID)

	u.SelectedCourseID = courseID
}

// getAuthenticationToken get lectures id which needs to download
func (u *Udemy) getAuthenticationToken() {
	var token string
	fmt.Print("⚠️  You must have your Udemy Access Token Ready ⚠️\n")

	fmt.Print("➡️  Enter your udemy access token: ")
	fmt.Scanln(&token)

	u.AccessToken = "Bearer " + token
}

// getLecturesIDs get lectures id which needs to download
func (u *Udemy) getLecturesIDs() {
	var start, end int
	fmt.Print("❓  : Enter the Lecture No. to start download from: ")
	fmt.Scanln(&start)

	fmt.Print("❓  : Enter the Lecture No. you want end download: ")
	fmt.Scanln(&end)

	u.Start = start
	u.End = end
}

// getVideoResolution get resolution which need to download
func (u *Udemy) getVideoResolution() {
	var resolution string
	fmt.Print("❓   : Enter the video Resolution(360|480|720|1080): ")
	fmt.Scanln(&resolution)
	u.Resolution = resolution
}

// GetDownloadLink to get the video download link
func (u *Udemy) GetDownloadLink(asset Asset) error {
	u.CurrentAttempt = u.CurrentAttempt + 1
	url := strings.Replace(GetDownloadURL, "{{assetID}}", strconv.Itoa(asset.ID), 1)

	res := u.NewRequest("GET", url)
	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var response Response
	json.Unmarshal(body, &response)

	var videosUrls = response.StreamUrls.Video
	for i := range videosUrls {
		if videosUrls[i].Label == u.Resolution {
			return u.Download(videosUrls[i].File, asset)
		}
	}

	fmt.Printf("❌  Don't have any valid download link for resolution %v, try with different resolution. \n", u.Resolution)

	if u.SessionMaxAttempt >= u.CurrentAttempt {
		u.getVideoResolution()
		u.GetDownloadLink(asset)
	} else {
		fmt.Println("❌  Max attempt exceeded, please try again.")
		os.Exit(0)
	}
	return nil
}

func (u *Udemy) startDownloading(courseAsset []Asset) {
	if u.Start == 0 || u.End == 0 {
		u.getLecturesIDs()
	}

	if u.Resolution == "false" {
		u.getVideoResolution()
	}

	for l := u.Start; l <= u.End; l++ {
		if courseAsset[l].ID != 0 {
			u.GetDownloadLink(courseAsset[l])
		}
	}
}

// Download to download files and vidoes
func (u *Udemy) Download(downloadURL string, asset Asset) error {

	out, err := os.Create(strconv.Itoa(asset.ObjectIndex) + ". " + asset.Title + ".mp4")
	if err != nil {
		return errors.New("❌  Error creating a new file, try to download from link" + downloadURL)
	}
	defer out.Close()

	resp, err := http.Head(downloadURL)
	if err != nil {
		fmt.Print(err)
		return err
	}

	size, err := strconv.ParseInt(resp.Header.Get("Content-Length"), 10, 64)
	defer resp.Body.Close()

	res, err := http.Get(downloadURL)
	if err != nil {
		return errors.New("❌  Error in downloading video, try to download from link" + downloadURL)
	}
	defer res.Body.Close()

	r := progress.NewReader(res.Body)

	// Start a goroutine printing progress
	go func() {
		ctx := context.Background()
		progressChan := progress.NewTicker(ctx, r, size, 1*time.Second)
		for p := range progressChan {
			fmt.Printf("\r   🔻 Downloading : %v(%.2f MB/%.2f MB)", asset.Title, BytesToMegaBytes(p.N()), BytesToMegaBytes(p.Size()))
		}
		s := strconv.FormatFloat(BytesToMegaBytes(size), 'f', -1, 64)
		fmt.Println("   - Download Finished : " + asset.Title + "(" + s + "MB)")
	}()

	var _, copyError = io.Copy(out, r)

	if copyError != nil {
		return copyError
	}
	return nil
}

// ParseHTMLAndGetCourseID it will parse the html content and get course id
func (u *Udemy) ParseHTMLAndGetCourseID() {
	res := u.NewRequest("GET", u.CourseURL)
	defer res.Body.Close()

	body, _ := ioutil.ReadAll(res.Body)

	bodyString := string(body)
	doc, err := html.Parse(strings.NewReader(bodyString))

	if err != nil {
		log.Fatal(err)
	}

	var f func(*html.Node)
	f = func(n *html.Node) {
		if n.Type == html.ElementNode && n.Data == "body" {
			for _, a := range n.Attr {
				if a.Key == "data-clp-course-id" {
					u.SelectedCourseID = a.Val
					break
				}
			}
		}
		for c := n.FirstChild; c != nil; c = c.NextSibling {
			f(c)
		}
	}

	f(doc)
}
