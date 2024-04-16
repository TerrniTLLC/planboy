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
	Database *notionapi.DatabaseQueryResponse
	BlockId  notionapi.BlockID
}

func (np *NotionApi) QueryHabits(pageId string) (*notionapi.DatabaseQueryResponse, error) {
	var reqParams *notionapi.DatabaseQueryRequest
	// db, err := np.Client.Database.Get(context.Background(), notionapi.DatabaseID(pageId))
	res, err := np.Client.Database.Query(context.Background(), notionapi.DatabaseID(pageId), reqParams)
	if err != nil {
		return nil, err
	}

	return res, nil
}

func (np *NotionApi) QueryTodayHabits(pageId string, params *notionapi.DatabaseQueryRequest) (*notionapi.DatabaseQueryResponse, error) {
	res, err := np.Client.Database.Query(context.Background(), notionapi.DatabaseID(pageId), params)
	if err != nil {
		return nil, err
	}

	for _, item := range res.Results {
		fmt.Println(item)
	}

	return res, nil
}

func NewNotionApi(token, pageId string) (*NotionApi, error) {
	np := &NotionApi{
		Client: notionapi.NewClient(notionapi.Token(token)),
	}

	// np.Database = db
	// np.MainPage = page
	// np.BlockId = notionapi.BlockID(page.ID)

	// s, err := json.MarshalIndent(np.Database, "", " ")
	// if err != nil {
	// return nil, err
	// }

	// fmt.Println(string(s))

	return np, nil
}

func getTimestamp() string {
	tm := time.Now()
	return fmt.Sprintf("[%s]\n", tm.Format("02.01.06 15:04"))
}
