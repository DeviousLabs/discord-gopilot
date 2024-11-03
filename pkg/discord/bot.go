package discord

import (
	"bytes"
	"fmt"
	"os"
	"regexp"
	"strings"

	"github.com/DeviousLabs/discord-gopilot/pkg/ai"

	"github.com/bwmarrin/discordgo"
	"golang.org/x/net/html"
)

type Bot struct {
	Session      *discordgo.Session
	MessageQueue *Queue
}

type MessageReference struct {
	ChannelID string
	Content   string
	GuildID   string
	MessageID string
}

func NewBot(token string) (*Bot, error) {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		return nil, err
	}
	session.Identify.Intents = discordgo.IntentsGuildMessages | discordgo.IntentMessageContent
	bot := &Bot{
		Session: session,
		// MAX Queue size of 255, cannot be larger unless Queue is changed from a uint8
		MessageQueue: NewQueue(255),
	}
	session.AddHandler(bot.messageCreate)
	return bot, nil
}

// Set a personality for the bot
func Personality() string {
	persona := os.Getenv("PERSONA")
	if persona == "" {
		fmt.Println("PERSONA is not set...")
		fmt.Println("Using default persona...")
	}
	switch persona {
	case "developer":
		return "Imagine you are a senior software developer with over 10 years of experience in building scalable applications. You are proficient in multiple programming language and have extensive knowledge of frameworks. Your approach emphasizes best practices, such as test-driven development and continuous integration. You are also known for your clean, efficient, and maintainable code." +
			"In this scenario, you are advising a junior developer on designing a new feature that enhances user experience without compromising performance. Provide detailed guidance on the architectural choices, coding standards to follow, and any potential pitfalls to avoid."
	default:
		return "You are an AI assistant designed to provide accurate, clear, and helpful responses to user inquiries. Prioritize delivering information that is factual, well-researched, and up-to-date. Avoid assumptions and focus on providing reliable guidance. If you are unsure of an answer, be transparent about your limitations. Respond in a friendly and professional manner, using concise language that is easy to understand."
	}
}

func (bot *Bot) Start() error {
	if err := bot.Session.Open(); err != nil {
		return err
	}

	// Start processing messages from the Queue
	go bot.processMessages()

	return nil
}

func (bot *Bot) Stop() {
	bot.MessageQueue.Close() // Ensure the queue is properly closed
	bot.Session.Close()
	bot.MessageQueue.wg.Wait()
}

func (bot *Bot) messageCreate(session *discordgo.Session, message *discordgo.MessageCreate) {
	if message.Author.ID == session.State.User.ID {
		return // Ignore bot's own messages to prevent loops
	}

	// Construct the message reference
	msgRef := MessageReference{
		MessageID: message.ID,
		ChannelID: message.ChannelID,
		GuildID:   message.GuildID,
		Content:   message.Content,
	}

	// Enqueue the message reference for later processing
	bot.MessageQueue.Enqueue(msgRef)
}

func (bot *Bot) processMessages() {
	for {
		item := bot.MessageQueue.Dequeue()
		msgRef, ok := item.(MessageReference)
		if !ok {
			continue
		}

		persona := Personality()
		content := msgRef.Content
		mentionRegex := regexp.MustCompile(`<@!?` + regexp.QuoteMeta(bot.Session.State.User.ID) + `>`)
		content = strings.TrimSpace(mentionRegex.ReplaceAllString(content, ""))

		if content == "" {
			continue
		}

		bot.Session.ChannelTyping(msgRef.ChannelID)
		response, err := bot.sendRequestToCloudflareAI(persona + " " + content)
		if err != nil {
			bot.Session.ChannelMessageSend(msgRef.ChannelID, "Error processing your request.")
			continue
		}

		formattedResponse, err := HTMLToDiscordMarkdown(response)
		if err != nil {
			bot.Session.ChannelMessageSend(msgRef.ChannelID, "Error formatting response.")
			continue
		}

		bot.Session.ChannelMessageSend(msgRef.ChannelID, formattedResponse)
	}
}

func (bot *Bot) sendRequestToCloudflareAI(query string) (string, error) {
	accountID := os.Getenv("CLOUDFLARE_ACCOUNT_ID")
	apiKey := os.Getenv("CLOUDFLARE_API_TOKEN")
	if accountID == "" {
		return "", fmt.Errorf("CLOUDFLARE_ACCOUNT_ID is not set")
	}
	if apiKey == "" {
		return "", fmt.Errorf("CLOUDFLARE_API_TOKEN is not set")
	}
	response, err := ai.RunCloudflareAI(accountID, apiKey, query)
	if err != nil {
		return "", err
	}

	if !response.Success {
		return "", fmt.Errorf("Cloudflare API error: %v", response.Errors)
	}

	return response.Result.Response, nil
}

func HTMLToDiscordMarkdown(htmlContent string) (string, error) {
	doc, err := html.Parse(strings.NewReader(htmlContent))
	if err != nil {
		return "", err
	}
	var buf bytes.Buffer
	traverseNodes(doc, &buf, false, false)
	return buf.String(), nil
}

func traverseNodes(n *html.Node, buf *bytes.Buffer, insidePre bool, skip bool) {
	currentInsidePre, currentSkip := handleElementNodeStart(n, buf, insidePre, skip)
	if !currentSkip {
		handleTextNode(n, buf, currentInsidePre)
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		traverseNodes(c, buf, currentInsidePre, currentSkip)
	}
	handleElementNodeEnd(n, buf, currentInsidePre)
}

func handleElementNodeStart(n *html.Node, buf *bytes.Buffer, insidePre bool, skip bool) (bool, bool) {
	currentInsidePre := insidePre
	currentSkip := skip
	if n.Type == html.ElementNode {
		switch n.Data {
		case "pre":
			currentInsidePre = true  // Entering a <pre> tag
			currentSkip = false      // Ensure skip is reset when entering <pre>
			buf.WriteString("\n```") // Start code block for pre content
		case "div":
			if hasTargetDivClass(n) {
				currentSkip = true // Skip this div and its contents
			}
		case "code":
			if !currentInsidePre {
				buf.WriteString(" `") // Start inline code with a leading space for separation
			}
		}
	}
	return currentInsidePre, currentSkip
}

func handleTextNode(n *html.Node, buf *bytes.Buffer, insidePre bool) {
	if n.Type == html.TextNode {
		if insidePre {
			buf.WriteString(n.Data) // Preserve exact text inside <pre>
		} else {
			cleanText := strings.ReplaceAll(n.Data, "\n", " ") // Normalize space for text outside <pre>
			cleanText = strings.TrimSpace(cleanText)
			if len(cleanText) > 0 {
				if n.Parent != nil && n.Parent.Data == "code" && !insidePre {
					// If it's code text outside <pre>, don't add leading/trailing spaces
					buf.WriteString(cleanText)
				} else {
					buf.WriteString(" " + cleanText + " ") // Ensure spaces around normal text
				}
			}
		}
	}
}

func handleElementNodeEnd(n *html.Node, buf *bytes.Buffer, insidePre bool) {
	if n.Type == html.ElementNode {
		switch n.Data {
		case "code":
			if !insidePre {
				buf.WriteString("` ") // End inline code with a trailing space for separation
			}
		case "pre":
			buf.WriteString("```") // End code block for pre content
		}
	}
}

func hasTargetDivClass(n *html.Node) bool {
	for _, a := range n.Attr {
		if a.Key == "class" && strings.Contains(a.Val, "text-token-text-secondary") && strings.Contains(a.Val, "bg-token-main-surface-secondary") {
			return true
		}
	}
	return false
}

// reMultiLineCode := regexp.MustCompile("(?s)\\`\\`\\`(.*?)\\`\\`\\`")
func stripDiscordMarkdown(content string) string {
	// Remove multi-line code blocks
	reMultiLineCode := regexp.MustCompile("(?s)\\`\\`\\`(.*?)\\`\\`\\`")
	content = reMultiLineCode.ReplaceAllString(content, "$1")
	// Remove inline code blocks
	reInlineCode := regexp.MustCompile("`([^`]*)`")
	content = reInlineCode.ReplaceAllString(content, "$1")
	// Remove bold, italic, underline, and strikethrough by matching explicitly without back-references
	replacements := []struct {
		Regex       *regexp.Regexp
		Replacement string
	}{
		{Regex: regexp.MustCompile(`\*\*\*(.*?)\*\*\*`), Replacement: "$1"}, // Bold + Italic
		{Regex: regexp.MustCompile(`___(.*?)___`), Replacement: "$1"},       // Bold + Italic
		{Regex: regexp.MustCompile(`\*\*(.*?)\*\*`), Replacement: "$1"},     // Bold
		{Regex: regexp.MustCompile(`__(.*?)__`), Replacement: "$1"},         // Bold
		{Regex: regexp.MustCompile(`\*(.*?)\*`), Replacement: "$1"},         // Italic
		{Regex: regexp.MustCompile(`_(.*?)_`), Replacement: "$1"},           // Italic
		{Regex: regexp.MustCompile(`~~(.*?)~~`), Replacement: "$1"},         // Strikethrough
	}
	// Execute each replacement
	for _, repl := range replacements {
		content = repl.Regex.ReplaceAllString(content, repl.Replacement)
	}
	// Normalize whitespace and new lines
	content = strings.TrimSpace(content)
	content = strings.ReplaceAll(content, "\n\n", "\n")
	content = regexp.MustCompile(`\s+`).ReplaceAllString(content, " ") // Collapse multiple spaces into one
	return content
}
