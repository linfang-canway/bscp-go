/*
 * Tencent is pleased to support the open source community by making Blueking Container Service available.
 * Copyright (C) 2019 THL A29 Limited, a Tencent company. All rights reserved.
 * Licensed under the MIT License (the "License"); you may not use this file except
 * in compliance with the License. You may obtain a copy of the License at
 * http://opensource.org/licenses/MIT
 * Unless required by applicable law or agreed to in writing, software distributed under
 * the License is distributed on an "AS IS" BASIS, WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND,
 * either express or implied. See the License for the specific language governing permissions and
 * limitations under the License.
 *
 */

package upstream

import (
	"fmt"

	"bscp.io/pkg/kit"
	pbfs "bscp.io/pkg/protocol/feed-server"
	sfs "bscp.io/pkg/sf-share"
)

// Handshake to the upstream server
func (uc *upstreamClient) Handshake(vas *kit.Vas, msg *pbfs.HandshakeMessage) (*pbfs.HandshakeResp, error) {
	if err := uc.wait.WaitWithContext(vas.Ctx); err != nil {
		return nil, err
	}

	resp, err := uc.client.Handshake(vas.Ctx, msg)
	if err != nil {
		return nil, err
	}

	return resp, nil
}

// Watch release related messages from upstream feed server.
func (uc *upstreamClient) Watch(vas *kit.Vas, payload []byte) (pbfs.Upstream_WatchClient, error) {
	if err := uc.wait.WaitWithContext(vas.Ctx); err != nil {
		return nil, err
	}

	meta := &pbfs.SideWatchMeta{
		ApiVersion: sfs.CurrentAPIVersion,
		Payload:    payload,
	}

	return uc.client.Watch(vas.Ctx, meta)
}

// Messaging is a message pipeline to send message to the upstream feed server.
func (uc *upstreamClient) Messaging(vas *kit.Vas, typ sfs.MessagingType, payload []byte) (*pbfs.MessagingResp,
	error) {

	if err := uc.wait.WaitWithContext(vas.Ctx); err != nil {
		return nil, err
	}

	if err := typ.Validate(); err != nil {
		return nil, fmt.Errorf("invalid message type, %s", err.Error())
	}

	msg := &pbfs.MessagingMeta{
		ApiVersion: sfs.CurrentAPIVersion,
		Rid:        vas.Rid,
		Type:       uint32(typ),
		Payload:    payload,
	}

	return uc.client.Messaging(vas.Ctx, msg)
}

// PullAppFileMeta pulls the app file meta from upstream feed server.
func (uc *upstreamClient) PullAppFileMeta(vas *kit.Vas, req *pbfs.PullAppFileMetaReq) (
	*pbfs.PullAppFileMetaResp, error) {

	if err := uc.wait.WaitWithContext(vas.Ctx); err != nil {
		return nil, err
	}

	return uc.client.PullAppFileMeta(vas.Ctx, req)
}

// GetDownloadURL gets the file temp download url from upstream feed server.
func (uc *upstreamClient) GetDownloadURL(vas *kit.Vas, req *pbfs.GetDownloadURLReq) (*pbfs.GetDownloadURLResp, error) {

	if err := uc.wait.WaitWithContext(vas.Ctx); err != nil {
		return nil, err
	}

	return uc.client.GetDownloadURL(vas.Ctx, req)
}
