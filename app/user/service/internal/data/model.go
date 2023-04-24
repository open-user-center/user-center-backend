package data

import "time"

//easyjson:json
type User struct {
	Id           int32
	UserName     string `gorm:"column:username"`
	UserAccount  string `gorm:"column:userAccount"`
	AvatarUrl    string `gorm:"column:avatarUrl"`
	Gender       int32
	UserPassword string `gorm:"column:userPassword"`
	Phone        string
	Email        string
	UserStatus   int32     `gorm:"column:userStatus"`
	CreateTime   time.Time `gorm:"column:createTime"`
	UpdateTime   time.Time `gorm:"column:updateTime"`
	IsDelete     int32     `gorm:"column:isDelete"`
	Role         int32
}

////easyjson:json
//type Profile struct {
//	CreatedAt time.Time
//	Updated   int64
//	Uuid      string `gorm:"primaryKey;size:20"`
//	Username  string `gorm:"index;not null;size:100"`
//	Avatar    string `gorm:"size:200"`
//	School    string `gorm:"size:50"`
//	Company   string `gorm:"size:50"`
//	Job       string `gorm:"size:50"`
//	Homepage  string `gorm:"size:100"`
//	Github    string `gorm:"size:100"`
//	Gitee     string `gorm:"size:100"`
//	Introduce string `gorm:"size:100"`
//}
//
//type ProfileUpdate struct {
//	Profile
//	Status int32 `gorm:"default:1"`
//}
//
//type AvatarReview struct {
//	gorm.Model
//	Uuid     string `gorm:"index;size:20"`
//	JobId    string `gorm:"size:100"`
//	Url      string `gorm:"size:1000"`
//	Label    string `gorm:"size:100"`
//	Result   int32
//	Category string `gorm:"size:100"`
//	SubLabel string `gorm:"size:100"`
//	Score    int32
//}
//
//type CoverReview struct {
//	gorm.Model
//	Uuid     string `gorm:"index;size:20"`
//	JobId    string `gorm:"size:100"`
//	Url      string `gorm:"size:1000"`
//	Label    string `gorm:"size:100"`
//	Result   int32
//	Category string `gorm:"size:100"`
//	SubLabel string `gorm:"size:100"`
//	Score    int32
//}
//
//type Follow struct {
//	gorm.Model
//	Follow   string `gorm:"uniqueIndex:idx_follow;size:20"`
//	Followed string `gorm:"uniqueIndex:idx_follow;index;size:20"`
//	Status   int32
//}
