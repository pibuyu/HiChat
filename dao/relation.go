package dao

import (
	"HiChat/global"
	"HiChat/models"

	"errors"

	"go.uber.org/zap"
)

// FriendList 获取好友列表
func FriendList(userId uint) (*[]models.UserBasic, error) {
	relation := make([]models.Relation, 0)
	if tx := global.DB.Where("owner_id = ? and type=1", userId).Find(&relation); tx.RowsAffected == 0 {
		zap.S().Info("未查询到Relation数据")
		return nil, errors.New("未查到好友关系")
	}

	userID := make([]uint, 0)
	for _, v := range relation {
		userID = append(userID, v.TargetID)
	}

	user := make([]models.UserBasic, 0)
	if tx := global.DB.Where("id in ?", userID).Find(&user); tx.RowsAffected == 0 {
		zap.S().Info("未查询到Relation好友关系")
		return nil, errors.New("未查到好友")
	}
	return &user, nil
}

// AddFriendByName 昵称加好友
func AddFriendByName(userId uint, targetName string) (int, error) {
	user, err := FindUserByName(targetName)
	if err != nil {
		return -1, errors.New("该用户不存在")
	}
	if user.ID == 0 {
		zap.S().Info("未查询到用户")
		return -1, errors.New("该用户不存在")
	}
	return AddFriend(userId, user.ID)
}

// DeleteFriend 删除好友
func DeleteFriend(userId uint, targetName string) (int, error) {
	//首先通过targetName查到对方的userId
	zap.S().Info("待删除的好友name：", targetName)
	targetUser, err := FindUserByName(targetName)
	if err != nil {
		return -1, errors.New("查询要删除的好友出错：" + err.Error())
	}
	tarId := targetUser.ID
	zap.S().Info("待删除的好友id：", tarId)
	//然后去relation表中双向删除好友
	tx := global.DB.Debug().Delete(&models.Relation{}, "owner_id = ? and target_id = ? and type = 1", userId, tarId)
	if tx.RowsAffected == 0 {
		return -1, errors.New("删除好友失败")
	}
	tx = global.DB.Debug().Delete(&models.Relation{}, "owner_id = ? and target_id = ? and type = 1", tarId, userId)
	if tx.RowsAffected == 0 {
		return -1, errors.New("删除好友失败")
	}
	return 0, nil
}

// AddFriend 加好友
func AddFriend(userID, TargetId uint) (int, error) {

	//不能加自己为好友
	if userID == TargetId {
		return -2, errors.New("userID和TargetId相等")
	}
	//通过id查询用户
	targetUser, err := FindUserID(TargetId)
	if err != nil {
		return -1, errors.New("未查询到用户")
	}
	if targetUser.ID == 0 {
		zap.S().Info("未查询到用户")
		return -1, errors.New("未查询到用户")
	}

	relation := models.Relation{}

	//双向查询是否已经有好友
	if tx := global.DB.Debug().Where("owner_id = ? and target_id = ? and type = 1", userID, TargetId).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("该好友存在")
		return 0, errors.New("好友已经存在")
	}

	if tx := global.DB.Debug().Where("owner_id = ? and target_id = ?  and type = 1", TargetId, userID).First(&relation); tx.RowsAffected == 1 {
		zap.S().Info("该好友存在")
		return 0, errors.New("好友已经存在")
	}

	//开启事务
	tx := global.DB.Begin()

	relation.OwnerId = userID
	relation.TargetID = targetUser.ID
	relation.Type = 1

	//同时添加两条好友记录，双向加好友
	if t := tx.Debug().Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("创建失败")

		//事务回滚
		tx.Rollback()
		return -1, errors.New("创建好友记录失败")
	}

	relation = models.Relation{}
	relation.OwnerId = TargetId
	relation.TargetID = userID
	relation.Type = 1

	if t := tx.Debug().Create(&relation); t.RowsAffected == 0 {
		zap.S().Info("创建失败")

		//事务回滚
		tx.Rollback()
		return -1, errors.New("创建好友记录失败")
	}

	//提交事务
	tx.Commit()
	return 1, nil
}
