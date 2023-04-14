//go:build !js
// +build !js

package db

import (
	"context"
	"encoding/json"
	"errors"
	"open_im_sdk/pkg/db/model_struct"
	"open_im_sdk/pkg/log"
	"open_im_sdk/pkg/server_api_params"
	"open_im_sdk/pkg/utils"

	"github.com/OpenIMSDK/Open-IM-Server/pkg/proto/sdkws"
)

func (d *DataBase) GetMessageReactionExtension(ctx context.Context, msgID string) (result *model_struct.LocalChatLogReactionExtensions, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var l model_struct.LocalChatLogReactionExtensions
	return &l, utils.Wrap(d.conn.WithContext(ctx).Where("client_msg_id = ?",
		msgID).Take(&l).Error, "GetMessageReactionExtension failed")
}

func (d *DataBase) InsertMessageReactionExtension(ctx context.Context, messageReactionExtension *model_struct.LocalChatLogReactionExtensions) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	return utils.Wrap(d.conn.WithContext(ctx).Create(messageReactionExtension).Error, "InsertMessageReactionExtension failed")
}
func (d *DataBase) UpdateMessageReactionExtension(ctx context.Context, c *model_struct.LocalChatLogReactionExtensions) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	t := d.conn.WithContext(ctx).Updates(c)
	if t.RowsAffected == 0 {
		return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
	}
	return utils.Wrap(t.Error, "UpdateConversation failed")
}
func (d *DataBase) GetAndUpdateMessageReactionExtension(ctx context.Context, msgID string, m map[string]*sdkws.KeyValue) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var temp model_struct.LocalChatLogReactionExtensions
	err := d.conn.WithContext(ctx).Where("client_msg_id = ?",
		msgID).Take(&temp).Error
	if err != nil {
		temp.ClientMsgID = msgID
		temp.LocalReactionExtensions = []byte(utils.StructToJsonString(m))
		return d.conn.WithContext(ctx).Create(&temp).Error
	} else {
		oldKeyValue := make(map[string]*sdkws.KeyValue)
		err = json.Unmarshal(temp.LocalReactionExtensions, &oldKeyValue)
		if err != nil {
			log.Error("special handle", err.Error())
		}
		log.Warn("special handle", oldKeyValue)
		for k, newValue := range m {
			oldKeyValue[k] = newValue
		}
		temp.LocalReactionExtensions = []byte(utils.StructToJsonString(oldKeyValue))
		t := d.conn.WithContext(ctx).Updates(temp)
		if t.RowsAffected == 0 {
			return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
		}
	}
	return nil
}
func (d *DataBase) DeleteMessageReactionExtension(ctx context.Context, msgID string) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	temp := model_struct.LocalChatLogReactionExtensions{ClientMsgID: msgID}
	return d.conn.WithContext(ctx).Delete(&temp).Error

}
func (d *DataBase) DeleteAndUpdateMessageReactionExtension(ctx context.Context, msgID string, m map[string]*sdkws.KeyValue) error {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var temp model_struct.LocalChatLogReactionExtensions
	err := d.conn.WithContext(ctx).Where("client_msg_id = ?",
		msgID).Take(&temp).Error
	if err != nil {
		return err
	} else {
		oldKeyValue := make(map[string]*server_api_params.KeyValue)
		_ = json.Unmarshal(temp.LocalReactionExtensions, &oldKeyValue)
		for k, _ := range m {
			if _, ok := oldKeyValue[k]; ok {
				delete(oldKeyValue, k)
			}
		}
		temp.LocalReactionExtensions = []byte(utils.StructToJsonString(oldKeyValue))
		t := d.conn.WithContext(ctx).Updates(temp)
		if t.RowsAffected == 0 {
			return utils.Wrap(errors.New("RowsAffected == 0"), "no update")
		}
	}
	return nil
}
func (d *DataBase) GetMultipleMessageReactionExtension(ctx context.Context, msgIDList []string) (result []*model_struct.LocalChatLogReactionExtensions, err error) {
	d.mRWMutex.Lock()
	defer d.mRWMutex.Unlock()
	var messageList []model_struct.LocalChatLogReactionExtensions
	err = utils.Wrap(d.conn.WithContext(ctx).Where("client_msg_id IN ?", msgIDList).Find(&messageList).Error, "GetMultipleMessageReactionExtension failed")
	for _, v := range messageList {
		v1 := v
		result = append(result, &v1)
	}
	return result, err
}
