package notion

import (
	"context"
	"fmt"
	"time"

	"github.com/jomei/notionapi"
)

type NotionApi struct {
	Client   *notionapi.Client
	MainPage *notionapi.Page
	BlockId  notionapi.BlockID
}

func NewNotionApi(token, pageId string) (*NotionApi, error) {
	np := &NotionApi{
		Client: notionapi.NewClient(notionapi.Token(token)),
	}
	page, err := np.Client.Page.Get(context.Background(), notionapi.PageID(pageId))
	if err != nil {
		return nil, err
	}
	np.MainPage = page
	np.BlockId = notionapi.BlockID(page.ID)

	fmt.Println(np.BlockId)
	return np, nil
}

func getTimestamp() string {
	tm := time.Now()
	return fmt.Sprintf("[%s]\n", tm.Format("02.01.06 15:04"))
}
