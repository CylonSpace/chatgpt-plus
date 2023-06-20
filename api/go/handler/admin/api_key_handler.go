package admin

import (
	"chatplus/core"
	"chatplus/core/types"
	"chatplus/handler"
	"chatplus/store/model"
	"chatplus/store/vo"
	"chatplus/utils"
	"chatplus/utils/resp"
	"time"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ApiKeyHandler struct {
	handler.BaseHandler
	db *gorm.DB
}

func NewApiKeyHandler(app *core.AppServer, db *gorm.DB) *ApiKeyHandler {
	h := ApiKeyHandler{db: db}
	h.App = app
	return &h
}

func (h *ApiKeyHandler) Save(c *gin.Context) {
	var data struct {
		Id         uint   `json:"id"`
		UserId     uint   `json:"user_id"`
		Value      string `json:"value"`
		LastUsedAt string `json:"last_used_at"`
		CreatedAt  int64  `json:"created_at"`
	}
	if err := c.ShouldBindJSON(&data); err != nil {
		resp.ERROR(c, types.InvalidArgs)
		return
	}

	apiKey := model.ApiKey{Value: data.Value, UserId: data.UserId, LastUsedAt: utils.Str2stamp(data.LastUsedAt)}
	apiKey.Id = data.Id
	if apiKey.Id > 0 {
		apiKey.CreatedAt = time.Unix(data.CreatedAt, 0)
	}
	res := h.db.Save(&apiKey)
	if res.Error != nil {
		resp.ERROR(c, "更新数据库失败！")
		return
	}

	var keyVo vo.ApiKey
	err := utils.CopyObject(apiKey, &keyVo)
	if err != nil {
		resp.ERROR(c, "数据拷贝失败！")
		return
	}
	keyVo.Id = apiKey.Id
	keyVo.CreatedAt = apiKey.CreatedAt.Unix()
	resp.SUCCESS(c, keyVo)
}

func (h *ApiKeyHandler) List(c *gin.Context) {
	userId := h.GetInt(c, "user_id", -1)
	query := h.db.Session(&gorm.Session{})
	if userId >= 0 {
		query = query.Where("user_id", userId)
	}
	var items []model.ApiKey
	var keys = make([]vo.ApiKey, 0)
	res := query.Find(&items)
	if res.Error == nil {
		for _, item := range items {
			var key vo.ApiKey
			err := utils.CopyObject(item, &key)
			if err == nil {
				key.Id = item.Id
				key.CreatedAt = item.CreatedAt.Unix()
				key.UpdatedAt = item.UpdatedAt.Unix()
				keys = append(keys, key)
			} else {
				logger.Error(err)
			}
		}
	}
	resp.SUCCESS(c, keys)
}

func (h *ApiKeyHandler) Remove(c *gin.Context) {
	id := h.GetInt(c, "id", 0)

	if id > 0 {
		res := h.db.Where("id = ?", id).Delete(&model.ApiKey{})
		if res.Error != nil {
			resp.ERROR(c, "更新数据库失败！")
			return
		}
	}
	resp.SUCCESS(c)
}
