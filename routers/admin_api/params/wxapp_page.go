package params

import (
	"github.com/xifengzhu/eshop/models"
)

type WxAppPageParams struct {
	Name                 string      `json:"name,omitempty"`
	Key                  string      `json:"key,omitempty"`
	PageData             models.JSON `json:"page_data,omitempty"`
	ShareSentence        string      `json:"share_sentence,omitempty"`
	ShareCover           string      `json:"share_cover,omitempty"`
	ShareBackgroundCover string      `json:"share_background_cover,omitempty"`
}
