package models

import "time"

type User struct {
	Id          string     `json:"id" gorm:"column:id;primary_key" bson:"_id" dynamodbav:"id" firestore:"-" validate:"required,max=40"`
	Username    string     `json:"username" gorm:"column:username" bson:"username" dynamodbav:"username" firestore:"username" validate:"required,username,max=100"`
	Email       string     `json:"email" gorm:"column:email" bson:"email" dynamodbav:"email" firestore:"email" validate:"email,max=100"`
	Phone       string     `json:"phone" gorm:"column:phone" bson:"phone" dynamodbav:"phone" firestore:"phone" validate:"required,phone,max=18"`
	DateOfBirth *time.Time `json:"dateOfBirth" gorm:"column:date_of_birth" bson:"dateOfBirth" dynamodbav:"dateOfBirth" firestore:"dateOfBirth"`
	CreateTime  *time.Time `json:"createTime" gorm:"column:create_time" bson:"createTime" dynamodbav:"createTime" firestore:"-"`
	UpdateTime  *time.Time `json:"updateTime" gorm:"column:update_time" bson:"updateTime" dynamodbav:"updateTime" firestore:"-"`
}
