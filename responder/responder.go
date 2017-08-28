package responder

import (
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/jrrickard/slack-ooo-responder/common"
	"github.com/jrrickard/slack-ooo-responder/utils"
	"github.com/nlopes/slack"
)

type Responder struct {
}

type channelLastSentTime struct {
	lastSentTime map[string]time.Time
	sync.Mutex
}

var lastSentCache channelLastSentTime

func getContactSuggestions() (slack.Attachment, bool) {
	var attachment slack.Attachment
	return attachment, false
}

func generateContactSuggestion(config *common.Config) (string, bool) {
	var contactSuggestion string
	var ok = false
	now := time.Now()
	nowHour := now.Hour()

	var selectedContact common.ContactSuggestion
	if config.Contacts.Len() > 0 {
		for _, contact := range config.Contacts {
			startHour := contact.BeginTime.Hour()
			endHour := contact.EndTime.Hour()

			match := nowHour >= startHour && nowHour <= endHour
			if match {
				selectedContact = contact
				break
			}
		}

		if &selectedContact != nil {
			contactSuggestion = fmt.Sprintf("If you need something right now, try contacting : %v", selectedContact.Users)
			ok = true

		}
	}
	return contactSuggestion, ok
}

func generateSuggestions(config *common.Config) []slack.Attachment {
	suggestions := make([]slack.Attachment, len(config.Suggestions))
	for _, suggestion := range config.Suggestions {
		attachment := slack.Attachment{
			Title:     suggestion.Text,
			TitleLink: suggestion.URL,
		}
		suggestions = append(suggestions, attachment)
	}
	return suggestions
}

func sendMessages(client *slack.Client, channel chan slack.Msg, config *common.Config) {

	for msg := range channel {
		params := slack.PostMessageParameters{}
		params.AsUser = true
		params.LinkNames = 1
		params.Markdown = true
		suggestions := generateSuggestions(config)
		if len(suggestions) > 0 {
			params.Attachments = suggestions
		}
		//fmt.Printf("%v", msg)
		contactSuggestion, ok := generateContactSuggestion(config)
		if ok {
			msg.Text = msg.Text + ". " + contactSuggestion
		}
		client.PostMessage(msg.Channel, msg.Text, params)
	}
}

func sendMessage(channel chan slack.Msg, destination string, config *common.Config) {
	msg := slack.Msg{Channel: destination,
		Text: config.Message}
	lastSentCache.Lock()
	var shouldSend bool
	lastSentTime, ok := lastSentCache.lastSentTime[destination]
	if !ok {
		shouldSend = true
	} else {

		shouldSend = time.Since(lastSentTime).Minutes() > float64(config.SurpressMessages)

	}
	if shouldSend {

		channel <- msg
		lastSentCache.lastSentTime[destination] = time.Now()
	}
	lastSentCache.Unlock()
}

func handleMessage(event *slack.MessageEvent, channel chan slack.Msg, config *common.Config) {

	timestamp, err := utils.ConvertTimestamp(event.Timestamp[0:strings.Index(event.Timestamp, ".")])

	if err != nil {
		log.Printf("Couldn't determine time, skipping : %v", err.Error())
		return
	}
	if timestamp.After(config.InitalizedTime) {
		//Ignore messages from myself
		if event.User != config.User {
			//This is a direct message, respond.
			if strings.HasPrefix(event.Channel, "D") {
				sendMessage(channel, event.Channel, config)
				//For now, only handle direct messages
				// } else {
				// 	if strings.Contains(event.Text, config.User.Name) {
				// 		//sendMessage(channel, event.Channel, config)
				// 	}
				// }
			}
		}
	}
}

func (responder *Responder) Connect() {
	config := *utils.GetConfig()
	lastSentCache = channelLastSentTime{lastSentTime: make(map[string]time.Time)}
	api := slack.New(config.Token)
	rtm := api.NewRTM()
	go rtm.ManageConnection()
	outgoing := make(chan slack.Msg)
	go sendMessages(api, outgoing, &config)
	for msg := range rtm.IncomingEvents {
		switch ev := msg.Data.(type) {
		case *slack.MessageEvent:
			go handleMessage(ev, outgoing, &config)
		default:
		}
	}
}
