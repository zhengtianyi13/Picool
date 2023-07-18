package serializer

import "mall/model"

type Boss struct {
	ID       uint   `json:"id"`
	BossName string `json:"boss_name"`
	NickName string `json:"nickname"`
	Type     int    `json:"type"`
	Email    string `json:"email"`
	Status   string `json:"status"`
	Avatar   string `json:"avatar"`
	CreateAt int64  `json:"create_at"`
}

//BuildUser 序列化用户
func BuildBoss(boss *model.Boss) Boss {
	return Boss{
		ID:       boss.ID,
		BossName: boss.BossName,
		NickName: boss.NickName,
		Email:    boss.Email,
		Status:   boss.Status,
		Avatar:   boss.AvatarURL(),
		CreateAt: boss.CreatedAt.Unix(),
	}
}

func BuildBosses(items []*model.Boss) (bosses []Boss) {
	for _, item := range items {
		boss := BuildBoss(item)
		bosses = append(bosses, boss)
	}
	return bosses
}
