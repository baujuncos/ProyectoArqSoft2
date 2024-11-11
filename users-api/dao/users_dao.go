package dao

type Users struct {
	User_id  int64  `gorm:"primaryKey;AUTO_INCREMENT;column:user_id"`
	Email    string `gorm:"type:varchar(100);not null;unique" binding:"required"`
	Password string `gorm:"type:varchar(100);not null" binding:"required"`
	Nombre   string `gorm:"type:varchar(100);not null" binding:"required"`
	Apellido string `gorm:"type:varchar(100);not null" binding:"required"`
	Admin    bool   `gorm:"type:boolean;not null" binding:"required"`
}
type Userss []Users
