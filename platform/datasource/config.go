package datasource

import "time"

type Config struct {
	ID        uint      `json:"id,omitempty" gorm:"primary_key;autoIncrement"`
	Key       string    `json:"key,omitempty" gorm:"column:c_key"`
	Field     string    `json:"field,omitempty" gorm:"column:c_field"`
	Value     string    `json:"value,omitempty" gorm:"column:c_value"`
	CreatedAt time.Time `json:"created_at,omitempty" gorm:"column:created_at;autoCreateTime:nano"`
	UpdatedAt time.Time `json:"updated_at,omitempty" gorm:"column:updated_at;autoUpdateTime:nano"`
}

func (c *Config) TableName() string {
	return "t_config"
}
