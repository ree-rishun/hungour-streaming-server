package services

import (
	"context"
	"fmt"
	"regexp"
	"strings"

	"github.com/google/generative-ai-go/genai"
	"google.golang.org/api/option"

	"hungour-streaming-server/models"
)

func GeminiRequest(ctx context.Context, message string, messages []models.Message) (string, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyB8UKIKK0PpAVvRto3P9yiXwe4CgNOx7Pg"))
	if err != nil {
		return "", err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	cs := model.StartChat()

	if err != nil {
		return "", err
	}

	var history []*genai.Content

	for _, ms := range messages {
		history = append(
			history,
			&genai.Content{
			Parts: []genai.Part{
				genai.Text(ms.Text),
			},
				Role: ms.Role,
			},
		)
	}

	cs.History = history

	resp, err := cs.SendMessage(ctx, genai.Text(message))
	if err != nil {
		return "", err
	}
	res := GeneratePlainTextResponse(resp.Candidates)

	return res, nil
}

func IsReserved(ctx context.Context, messages []models.Message) (bool, error) {
	client, err := genai.NewClient(ctx, option.WithAPIKey("AIzaSyB8UKIKK0PpAVvRto3P9yiXwe4CgNOx7Pg"))
	if err != nil {
		return false, err
	}
	defer client.Close()

	model := client.GenerativeModel("gemini-1.5-flash")

	cs := model.StartChat()

	if err != nil {
		return false, err
	}

	var history []*genai.Content

	for _, ms := range messages {
		history = append(
			history,
			&genai.Content{
				Parts: []genai.Part{
					genai.Text(ms.Text),
				},
				Role: ms.Role,
			},
		)
	}

	cs.History = history

	resp, err := cs.SendMessage(ctx, genai.Text(`
これまでの会話から今回電話したお店に対して予約は完了しましたか？
- 予約ができた場合はtrue
- 予約ができなかった場合はfalse
それ以外の文字は絶対に回答しないでください。
`))
	if err != nil {
		return false, err
	}
	res := GeneratePlainTextResponse(resp.Candidates)

	return res == "true", nil
}

func GeneratePlainTextResponse(cs []*genai.Candidate) string {
	var result string

	// Geminiのレスポンスを1つの文字列にまとめる
	for _, c := range cs {
		for _, p := range c.Content.Parts {
			result = fmt.Sprintf("[%s]:%s", result, p)
		}
	}

	// マークダウンをプレーンテキストに変換
	return stripMarkdown(result)
}

// Markdownのプレーンテキスト化関数
func stripMarkdown(input string) string {
	// 見出し (##, ### など) を削除
	re := regexp.MustCompile(`(?m)^#+\s*`)
	input = re.ReplaceAllString(input, "")

	// **太字** や *イタリック* の記号を削除
	re = regexp.MustCompile(`(\*\*|\*|__)`)
	input = re.ReplaceAllString(input, "")

	// [リンクテキスト](URL) → リンクテキスト のみにする
	re = regexp.MustCompile(`\[(.*?)\]\(.*?\)`)
	input = re.ReplaceAllString(input, "$1")

	// リスト記号 (-, *, •) を削除
	re = regexp.MustCompile(`(?m)^[-*•]\s*`)
	input = re.ReplaceAllString(input, "")

	// コードブロック (` ``` ` や ` `code` `) を削除
	re = regexp.MustCompile("(```.*?```|`[^`]*`)")
	input = re.ReplaceAllString(input, "")

	// 余分な空白や改行を整理
	input = strings.TrimSpace(input)

	return input
}
