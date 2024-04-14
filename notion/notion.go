package notion

import (
	"context"
	"fmt"
	"time"

	"github.com/jomei/notionapi"
	"github.com/terrnitllc/planboy/utils/jsonify"
)

type NotionApi struct {
	Client   *notionapi.Client
	MainPage *notionapi.Page
	Database *notionapi.Database
	BlockId  notionapi.BlockID
}

func NewNotionApi(token, pageId string) (*NotionApi, error) {
	np := &NotionApi{
		Client: notionapi.NewClient(notionapi.Token(token)),
	}
	db, err := np.Client.Database.Get(context.Background(), notionapi.DatabaseID(pageId))
	if err != nil {
		return nil, err
	}
	np.Database = db
	// np.MainPage = page
	// np.BlockId = notionapi.BlockID(page.ID)

	s, _ := jsonify.Marshal(np.Database)
	fmt.Println(string(s))

	return np, nil
}

func getTimestamp() string {
	tm := time.Now()
	return fmt.Sprintf("[%s]\n", tm.Format("02.01.06 15:04"))
}
