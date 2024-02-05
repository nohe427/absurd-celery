package webprocessor

import (
	"net/http"
	"strconv"
	"strings"

	"golang.org/x/net/html"
)

type PassageInfo struct {
	Act     int
	Scene   int
	Passage int
	Speaker string
	Text    string
}

type AllPassages struct {
	Passages *[]PassageInfo
}

var lastSpeaker string = ""

func LoadPage(url string) (*AllPassages, error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	// bytes, err := io.ReadAll(resp.Body)
	doc, err := html.Parse(resp.Body)
	if err != nil {
		return nil, err
	}
	passageList := make([]PassageInfo, 0)
	allPassages := &AllPassages{Passages: &passageList}
	processNode(doc, allPassages)
	// for _, v := range *allPassages.Passages {
	// 	fmt.Printf("%v\n", v.Text)
	// }
	return allPassages, nil
}

func IsSpeaker(value string) bool {
	return strings.HasPrefix(value, "speech")
}

func convertToInt(s string) int {
	i, err := strconv.Atoi(s)
	if err != nil {
		return -1
	}
	return i
}

func IsPassage(value string) (bool, *PassageInfo) {
	ss := strings.Split(value, ".")
	if len(ss) != 3 {
		return false, nil
	}

	return true, &PassageInfo{Act: convertToInt(ss[0]), Scene: convertToInt(ss[0]), Passage: convertToInt(ss[0]), Speaker: lastSpeaker}
}

func processNode(n *html.Node, passages *AllPassages) {
	if n.Data == "a" {
		interestingValue := ""
		for _, i := range n.Attr {
			// determine current speaker
			if i.Key == "name" {
				interestingValue = i.Val
				break
			}
		}
		if interestingValue != "" {
			if IsSpeaker(interestingValue) {
				lastSpeaker = n.FirstChild.FirstChild.Data
				// fmt.Printf("Speaker : %v\n", lastSpeaker)
			}
			isPassage, passage := IsPassage(interestingValue)
			if isPassage {
				passage.Text = n.FirstChild.Data
				// fmt.Printf("%v\n", passage.Text)
				*passages.Passages = append(*passages.Passages, *passage)
			}
		}
	}
	// fmt.Println(n.Data)
	if n.Type == html.ElementNode {
		// fmt.Println(n.Data)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		processNode(c, passages)
	}
}
