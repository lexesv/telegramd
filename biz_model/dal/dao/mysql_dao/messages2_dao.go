/*
 *  Copyright (c) 2017, https://github.com/nebulaim
 *  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package mysql_dao

import (
	"fmt"
	"github.com/golang/glog"
	"github.com/jmoiron/sqlx"
	"github.com/nebulaim/telegramd/biz_model/dal/dataobject"
	"github.com/nebulaim/telegramd/mtproto"
)

type Messages2DAO struct {
	db *sqlx.DB
}

func NewMessages2DAO(db *sqlx.DB) *Messages2DAO {
	return &Messages2DAO{db}
}

// insert into messages2(user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, random_id, message_type, message_data, date2) values (:user_id, :user_message_box_id, :dialog_message_id, :sender_user_id, :message_box_type, :peer_type, :peer_id, :random_id, :message_type, :message_data, :date2)
// TODO(@benqi): sqlmap
func (dao *Messages2DAO) Insert(do *dataobject.Messages2DO) int64 {
	var query = "insert into messages2(user_id, user_message_box_id, dialog_message_id, sender_user_id, message_box_type, peer_type, peer_id, random_id, message_type, message_data, date2) values (:user_id, :user_message_box_id, :dialog_message_id, :sender_user_id, :message_box_type, :peer_type, :peer_id, :random_id, :message_type, :message_data, :date2)"
	r, err := dao.db.NamedExec(query, do)
	if err != nil {
		errDesc := fmt.Sprintf("NamedExec in Insert(%v), error: %v", do, err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	id, err := r.LastInsertId()
	if err != nil {
		errDesc := fmt.Sprintf("LastInsertId in Insert(%v)_error: %v", do, err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}
	return id
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = :user_id and user_message_box_id in (:idList) order by user_message_box_id desc
// TODO(@benqi): sqlmap
func (dao *Messages2DAO) SelectByMessageIdList(user_id int32, idList []int32) []dataobject.Messages2DO {
	var q = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = ? and user_message_box_id in (?) order by user_message_box_id desc"
	query, a, err := sqlx.In(q, user_id, idList)
	rows, err := dao.db.Queryx(query, a...)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByMessageIdList(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.Messages2DO
	for rows.Next() {
		v := dataobject.Messages2DO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByMessageIdList(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	return values
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = :user_id and user_message_box_id = :user_message_box_id limit 1
// TODO(@benqi): sqlmap
func (dao *Messages2DAO) SelectByMessageId(user_id int32, user_message_box_id int32) *dataobject.Messages2DO {
	var query = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = ? and user_message_box_id = ? limit 1"
	rows, err := dao.db.Queryx(query, user_id, user_message_box_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectByMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.Messages2DO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectByMessageId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	return do
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = :user_id and peer_type = :peer_type and peer_id = :peer_id and user_message_box_id < :user_message_box_id order by user_message_box_id desc limit :limit
// TODO(@benqi): sqlmap
func (dao *Messages2DAO) SelectBackwardByPeerOffsetLimit(user_id int32, peer_type int8, peer_id int32, user_message_box_id int32, limit int32) []dataobject.Messages2DO {
	var query = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = ? and peer_type = ? and peer_id = ? and user_message_box_id < ? order by user_message_box_id desc limit ?"
	rows, err := dao.db.Queryx(query, user_id, peer_type, peer_id, user_message_box_id, limit)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectBackwardByPeerOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.Messages2DO
	for rows.Next() {
		v := dataobject.Messages2DO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectBackwardByPeerOffsetLimit(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	return values
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = :user_id and ((sender_user_id = :user_id and peer_id = :peer_id) or (sender_user_id = :peer_id and peer_id = :user_id)) and peer_type = :peer_type and user_message_box_id < :user_message_box_id order by user_message_box_id desc limit :limit
// TODO(@benqi): sqlmap
func (dao *Messages2DAO) SelectBackwardByPeerUserOffsetLimit(user_id int32, peer_id int32, peer_type int8, user_message_box_id int32, limit int32) []dataobject.Messages2DO {
	var query = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = ? and ((sender_user_id = ? and peer_id = ?) or (sender_user_id = ? and peer_id = ?)) and peer_type = ? and user_message_box_id < ? order by user_message_box_id desc limit ?"
	rows, err := dao.db.Queryx(query, user_id, user_id, peer_id, peer_id, user_id, peer_type, user_message_box_id, limit)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectBackwardByPeerUserOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.Messages2DO
	for rows.Next() {
		v := dataobject.Messages2DO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectBackwardByPeerUserOffsetLimit(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	return values
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = :user_id and peer_type = :peer_type and peer_id = :peer_id and user_message_box_id >= :user_message_box_id order by user_message_box_id asc limit :limit
// TODO(@benqi): sqlmap
func (dao *Messages2DAO) SelectForwardByPeerOffsetLimit(user_id int32, peer_type int8, peer_id int32, user_message_box_id int32, limit int32) []dataobject.Messages2DO {
	var query = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = ? and peer_type = ? and peer_id = ? and user_message_box_id >= ? order by user_message_box_id asc limit ?"
	rows, err := dao.db.Queryx(query, user_id, peer_type, peer_id, user_message_box_id, limit)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectForwardByPeerOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.Messages2DO
	for rows.Next() {
		v := dataobject.Messages2DO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectForwardByPeerOffsetLimit(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	return values
}

// select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = :user_id and ((sender_user_id = :user_id and peer_id = :peer_id) or (sender_user_id = :peer_id and peer_id = :user_id)) and peer_type = :peer_type and user_message_box_id >= :user_message_box_id order by user_message_box_id asc limit :limit
// TODO(@benqi): sqlmap
func (dao *Messages2DAO) SelectForwardByPeerUserOffsetLimit(user_id int32, peer_id int32, peer_type int8, user_message_box_id int32, limit int32) []dataobject.Messages2DO {
	var query = "select user_id, user_message_box_id, sender_user_id, message_box_type, peer_type, peer_id, message_type, message_data, date2 from messages2 where user_id = ? and ((sender_user_id = ? and peer_id = ?) or (sender_user_id = ? and peer_id = ?)) and peer_type = ? and user_message_box_id >= ? order by user_message_box_id asc limit ?"
	rows, err := dao.db.Queryx(query, user_id, user_id, peer_id, peer_id, user_id, peer_type, user_message_box_id, limit)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectForwardByPeerUserOffsetLimit(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	var values []dataobject.Messages2DO
	for rows.Next() {
		v := dataobject.Messages2DO{}

		// TODO(@benqi): 不使用反射
		err := rows.StructScan(&v)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectForwardByPeerUserOffsetLimit(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
		values = append(values, v)
	}

	return values
}

// select user_message_box_id, message_box_type from messages2 where user_id = :peerId and dialog_message_id = (select dialog_message_id from messages2 where user_id = :user_id and user_message_box_id = :user_message_box_id limit 1)
// TODO(@benqi): sqlmap
func (dao *Messages2DAO) SelectPeerMessageId(peerId int32, user_id int32, user_message_box_id int32) *dataobject.Messages2DO {
	var query = "select user_message_box_id, message_box_type from messages2 where user_id = ? and dialog_message_id = (select dialog_message_id from messages2 where user_id = ? and user_message_box_id = ? limit 1)"
	rows, err := dao.db.Queryx(query, peerId, user_id, user_message_box_id)

	if err != nil {
		errDesc := fmt.Sprintf("Queryx in SelectPeerMessageId(_), error: %v", err)
		glog.Error(errDesc)
		panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
	}

	defer rows.Close()

	do := &dataobject.Messages2DO{}
	if rows.Next() {
		err = rows.StructScan(do)
		if err != nil {
			errDesc := fmt.Sprintf("StructScan in SelectPeerMessageId(_), error: %v", err)
			glog.Error(errDesc)
			panic(mtproto.NewRpcError(int32(mtproto.TLRpcErrorCodes_DBERR), errDesc))
		}
	} else {
		return nil
	}

	return do
}