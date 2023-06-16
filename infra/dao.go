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

func UserGetAvator(opendId string) ([]byte, error) {
	key := fmt.Sprintf("user.%s", opendId)
	in, err := KvGet(key)
	if err != nil {
		common.Logger.Warnf("KvGet user %s data failed:%s", opendId, err.Error())
		return nil, err
	}

	usr := &grpc.UserDto{}
	err = proto.Unmarshal(in, usr)
	if err != nil {
		common.Logger.Warnf("Unmarshal user %s data failed:%s", opendId, err.Error())
		return nil, err
	}

	fs, _ := GetFs(Aavator)
	blk, err := fs.Read(usr.AvatorRef)
	if err != nil {
		common.Logger.Warnf("Read user %s data failed:%s", opendId, err.Error())
		return nil, err
	}
	return blk.Content, nil
}
