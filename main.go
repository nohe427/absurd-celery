package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/nohe427/absurd-celery/genai"
	"github.com/nohe427/absurd-celery/webprocessor"
)

func getGeminiKey() string {
	return os.Getenv("GEMINI_KEY")
}

func fullText(allPassages *webprocessor.AllPassages) string {
	var sb strings.Builder
	cs := ""
	for _, pi := range *allPassages.Passages {
		if pi.Speaker != cs {
			// Append the current speaker on its own line
			sb.WriteString(fmt.Sprintf("%v\n", pi.Speaker))
		}
		sb.WriteString(fmt.Sprintf("%v\n", pi.Text))
	}
	return sb.String()
}

func byAct(allPassages *webprocessor.AllPassages) []string {
	var sb strings.Builder
	cs := ""
	ca := 0
	byAct := make([]string, 0)
	for _, pi := range *allPassages.Passages {
		if pi.Act != ca {
			if ca != 0 {
				byAct = append(byAct, sb.String())
				sb.Reset()
			}
			ca = pi.Act
			// reset the current speaker
			cs = ""
		}
		if pi.Speaker != cs {
			// Append the current speaker on its own line
			sb.WriteString(fmt.Sprintf("%v\n", pi.Speaker))
		}
		sb.WriteString(fmt.Sprintf("%v\n", pi.Text))

	}
	return byAct
}

func main() {
	fmt.Println("Hello world")
	passages, err := webprocessor.LoadPage("https://shakespeare.mit.edu/hamlet/full.html")
	if err != nil {
		fmt.Println(err)
		return
	}
	key := getGeminiKey()
	if key == "" {
		fmt.Println("Please set your GEMINI_KEY env variable to your API Key")
		return
	}
	client := genai.NewGeminiClient(key)
	tc := client.TokenCount(fullText(passages))
	fmt.Printf("Full Text Token Count : %v\n", tc)
	if tc > 32_700 {
		fmt.Printf("Token Count Exceeds summarization Limit\n")
	}

	fmt.Printf("Parsing by Act\n")
	acts := byAct(passages)

	for i, v := range acts {
		fmt.Printf("Act %v\n", i+1)

		tc := client.TokenCount(v)
		fmt.Printf("Token Count : %v\n", tc)
		if tc > 32_600 {
			fmt.Printf("Token Count Exceeds summarization Limit\n")
		}
	}

	fmt.Printf("Summarizing by Act\n")
	actSummaries := make([]string, 0)
	for i, v := range acts {
		prompt := fmt.Sprintf("Summarize the following Shapespeare Act in modern English:\nAct %v\n%v", i+1, v)
		summary, err := client.Summarize(prompt)
		if err != nil {
			fmt.Printf("error : %v", err)
		}
		actSummaries = append(actSummaries, summary)
		fmt.Printf("Act %v Summay:\n%v\n", i+1, summary)
	}

	var sb strings.Builder
	sb.WriteString("Summarize all of the summaries of these Shakepseare acts into one narrative.\n")
	for i, v := range actSummaries {
		sb.WriteString(fmt.Sprintf("Act %v\nSumamry:%v\n", i+1, v))
	}
	finalPrompt := sb.String()
	summary, err := client.Summarize(finalPrompt)
	if err != nil {
		fmt.Printf("ERROR : %v", err)
		return
	}
	fmt.Printf("\n\n\nSummary of the entire play:\n\n%v\n\n", summary)
}
