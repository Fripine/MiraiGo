package client

import (
	"fmt"

	"github.com/Fripine/MiraiGo/client/pb/richmedia"
	"github.com/Fripine/MiraiGo/internal/proto"
)

// OidbSvcTrpcTcp.0x126d_200
func (c *QQClient) buildRecordDownloadReqPacket(Uid string, FileId string, groupUin int64, isGroup bool) (uint16, []byte) {
	scene := &richmedia.SceneInfo{
		RequestType:  2,
		BusinessType: 3,
		SceneType:    1,
		C2C: &richmedia.C2CUserInfo{
			AccountType: 2,
			TargetUid:   Uid,
		},
	}
	if isGroup {
		scene.RequestType = 1
		scene.SceneType = 2
		scene.Group = &richmedia.NTGroupInfo{
			GroupUin: uint32(groupUin),
		}
	}
	body := &richmedia.NTV2RichMediaReq{
		ReqHead: &richmedia.MultiMediaReqHead{
			Common: &richmedia.CommonHead{
				RequestId: 3,
				Command:   200,
			},
			Scene: scene,
			Client: &richmedia.ClientMeta{
				AgentType: 2,
			},
		},
		Download: &richmedia.DownloadReq{
			Node: &richmedia.IndexNode{
				Info: &richmedia.FileInfo{
					Type: &richmedia.FileType{
						Type:        3,
						VoiceFormat: 1,
					},
					Time: 1,
				},
				FileUuid: FileId,
				StoreId:  1,
			},
		},
	}
	b, err := proto.Marshal(body)
	if err != nil {
		fmt.Println(err)
	}
	payload := c.packOIDBPackage(4717, 200, b)
	return c.uniPacket("OidbSvcTrpcTcp.0x126d_200", payload)
}

func (c *QQClient) ParseRecordDownloadRspPacket(body []byte) (string, error) {
	rp := &richmedia.NTV2RichMediaRsp{}
	if err := proto.Unmarshal(body, rp); err != nil && rp.MediaResp.DownloadResp.Info != nil {
		c.error("parse RecordDownloadRspPacket error: %v", err)
		return "", err
	}
	return fmt.Sprintf("https://%s%s%s", rp.MediaResp.DownloadResp.Info.Domain, rp.MediaResp.DownloadResp.Info.UrlPath, rp.MediaResp.DownloadResp.Rkey), nil
}

// GetPrivateRecordDownloadUrl 获取私聊语音文件下载地址
func (c *QQClient) GetPrivateRecordDownloadUrl(selfUid string, FileId string) (string, error) {
	body, err := c.sendAndWaitDynamic(c.buildRecordDownloadReqPacket(selfUid, FileId, 0, false))
	if err != nil {
		return "", err
	}
	return c.ParseRecordDownloadRspPacket(body)
}

// GetGroupRecordDownloadUrl 获取群聊语音文件下载地址
func (c *QQClient) GetGroupRecordDownloadUrl(selfUid string, FileId string, groupUin int64) (string, error) {
	body, err := c.sendAndWaitDynamic(c.buildRecordDownloadReqPacket(selfUid, FileId, groupUin, true))
	if err != nil {
		return "", err
	}
	return c.ParseRecordDownloadRspPacket(body)
}
