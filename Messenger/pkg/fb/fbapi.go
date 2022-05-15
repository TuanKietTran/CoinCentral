package fb

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/valyala/fasthttp"
)

const (
	// uriSendMessage = "https://graph.facebook.com/v13.0/me/messages"
	uriSendMessage = "https://graph.facebook.com/v2.6/me/messages"
	// uriSendMessage = "https://graph.facebook.com/v13.0/me/messages" //?access_token=<PAGE_ACCESS_TOKEN>

	defaultRequestTimeout = 10 * time.Second
)

const (
	messageTypeResponse = "RESPONSE"
)

var (
	client = fasthttp.Client{}
)

// Respond responds to a user in FB messenger. This includes promotional and non-promotional messages sent inside the 24-hour standard messaging window.
// For example, use this tag to respond if a person asks for a reservation confirmation or an status update.
func Respond(ctx context.Context, recipientID, msgText string) error {
	return callAPI(ctx, uriSendMessage, SendMessageRequest{
		MessagingType: messageTypeResponse,
		RecipientID: MessageRecipient{
			ID: recipientID,
		},
		Message: Message{
			Text: msgText,
		},
	})
}

func popUpAllCoinButtons(ctx context.Context, textPostback string, recipientID string, buttons AttachmentButtons) error {
	att := Attachment{
		Type: "template",
		Payload: AttachmentPayload{
			TemplateType: "button",
			Text:         textPostback,
			Buttons:      buttons,
		},
	}
	return callAPI(ctx, uriSendMessage, SendMessageRequest{
		MessagingType: messageTypeResponse,
		RecipientID: MessageRecipient{
			ID: recipientID,
		},
		Message: Message{
			Attachment: &att,
		},
	})
}

func callAPI(ctx context.Context, reqURI string, reqBody SendMessageRequest) error {
	req := fasthttp.AcquireRequest()
	defer fasthttp.ReleaseRequest(req)

	req.SetRequestURI(fmt.Sprintf("%s?access_token=%s", reqURI, accessToken))
	req.Header.SetMethod(fasthttp.MethodPost)
	req.Header.Add("Content-Type", "application/json")

	body, err := json.Marshal(&reqBody)
	if err != nil {
		return fmt.Errorf("marshal: %w", err)
	}
	req.SetBody(body)

	res := fasthttp.AcquireResponse()
	defer fasthttp.ReleaseResponse(res)

	dl, ok := ctx.Deadline()
	if !ok {
		dl = time.Now().Add(defaultRequestTimeout)
	}

	err = client.DoDeadline(req, res, dl)
	if err != nil {
		return fmt.Errorf("do deadline: %w", err)
	}

	resp := APIResponse{}

	err = json.Unmarshal(res.Body(), &resp)
	if err != nil {
		return fmt.Errorf("unmarshal response: %w", err)
	}
	if resp.Error != nil {
		return fmt.Errorf("the response API is %s \n Response error: %s", resp.RecipientID, resp.Error.Error())
	}
	if res.StatusCode() != fasthttp.StatusOK {
		return fmt.Errorf("unexpected rsponse status %d", res.StatusCode())
	}

	return nil
}
