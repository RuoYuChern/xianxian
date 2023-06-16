package infra

import (
	"fmt"

	"google.golang.org/protobuf/proto"
	"xiyu.com/common"
	"xiyu.com/facade/grpc"
)

func UserSave(usr *grpc.UserDto) error {
	key := fmt.Sprintf("user.%s", usr.OpenId)
	out, err := proto.Marshal(usr)
	if err != nil {
		common.Logger.Warnf("Marshal user %s data failed:%s", usr.OpenId, err.Error())
		return err
	}

	err = KvBSet([]byte(key), out)
	if err != nil {
		common.Logger.Warnf("KvBSet user %s data failed:%s", usr.OpenId, err.Error())
		return err
	}
	return nil
}

func userGet(openId string, usr *grpc.UserDto) error {
	key := fmt.Sprintf("user.%s", openId)
	dat, err := KvGet(key)
	if err != nil {
		common.Logger.Infof("Get %s failed:%s", key, err.Error())
		return err
	}
	err = proto.Unmarshal(dat, usr)
	if err != nil {
		common.Logger.Warnf("Unmarshal user %s data failed:%s", openId, err.Error())
		return err
	}
	return nil
}

func UserSaveAvator(nickName string, openId string, dat []byte) error {
	fs, _ := GetFs(Aavator)
	ref, err := fs.Write(&grpc.FsBlockVo{OriginId: openId, Content: dat})
	if err != nil {
		common.Logger.Infof("write user:%s data failed:%s", openId, err.Error())
		return err
	}

	usr := &grpc.UserDto{}
	err = userGet(openId, usr)
	if err != nil {
		usr.OpenId = openId
		usr.NickName = nickName
		usr.AvatorRef = ref
	} else {
		usr.AvatorRef = ref
	}

	return UserSave(usr)

}

func UserGetAvator(openId string) ([]byte, error) {
	key := fmt.Sprintf("user.%s", openId)
	in, err := KvGet(key)
	if err != nil {
		common.Logger.Warnf("KvGet user %s data failed:%s", openId, err.Error())
		return nil, err
	}

	usr := &grpc.UserDto{}
	err = proto.Unmarshal(in, usr)
	if err != nil {
		common.Logger.Warnf("Unmarshal user %s data failed:%s", openId, err.Error())
		return nil, err
	}

	fs, _ := GetFs(Aavator)
	blk, err := fs.Read(usr.AvatorRef)
	if err != nil {
		common.Logger.Warnf("Read user %s data failed:%s", openId, err.Error())
		return nil, err
	}
	return blk.Content, nil
}
