package model

import (
	"time"

	"github.com/shopspring/decimal"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

type Asset struct {
	ID           uint            `gorm:"primaryKey;column:id;AUTO_INCREMENT" json:"id"`
	Name         string          `gorm:"column:name;not null" json:"name"`
	ParentID     uint            `gorm:"default:null;column:parent_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"parent_id"`
	Parent       *Asset          `gorm:"foreignKey:ParentID;references:ID;default:null" json:"parent"`
	UserID       uint            `gorm:"default:null;column:user_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"user_id"`
	User         User            `gorm:"foreignKey:UserID;references:ID;default:null" json:"user"`
	DepartmentID uint            `gorm:"default:null;column:department_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"department_id"`
	Department   Department      `gorm:"foreignKey:DepartmentID;references:ID;default:null" json:"department"`
	Price        decimal.Decimal `gorm:"type:decimal(10,2);column:price" json:"price"`
	Description  string          `gorm:"column:description" json:"description"`
	Position     string          `gorm:"column:position" json:"position"`
	Expire       uint            `gorm:"column:expire;default:0" json:"expire"`
	ClassID      uint            `gorm:"default:null;column:class_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"class_id"`
	Class        AssetClass      `gorm:"foreignKey:ClassID;references:ID;default:null" json:"class"`
	Number       int             `gorm:"column:number" json:"count"`
	Type         int             `gorm:"column:type" json:"type"`   // 1-条目型资产 2-数量型资产
	State        uint            `gorm:"column:state" json:"state"` // 0idle;1in_use;2in_maintain;3retired;4deleted
	MaintainerID uint            `gorm:"default:null;column:maintainer_id;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"maintainer_id"`
	Maintainer   User            `gorm:"foreignKey:MaintainerID;references:ID;default:null" json:"maintainer"`
	Property     datatypes.JSON  `gorm:"column:property;" json:"property"`
	TaskList     []*Task         `gorm:"many2many:task_assets;" json:"task_list"`
	CreatedAt    *ModelTime      `gorm:"column:created_at" json:"created_at"`
	NetWorth     decimal.Decimal `gorm:"type:decimal(10,2);column:net_worth" json:"net_worth"`
}

func (asset *Asset) BeforeSave(tx *gorm.DB) error {
	if asset.ID == 0 {
		return nil
	}

	var err error
	err = nil
	if asset.State < 3 && asset.Expire != 0 {
		interval := getDiffDays(time.Time(*asset.CreatedAt), time.Now())
		if interval >= int(asset.Expire) {
			asset.NetWorth = decimal.Zero
			asset.State = 3

			skipHookDB := tx.Session(&gorm.Session{
				SkipHooks: true,
			})

			var subAssets []*Asset
			result := skipHookDB.Model(&Asset{}).Where("parent_id = ?", asset.ID).Find(subAssets)

			if result.Error == gorm.ErrRecordNotFound {
				err = nil
			} else if result.Error != nil {
				err = result.Error
			} else {
				err = nil
			}
			if err == nil {
				for _, subAsset := range subAssets {
					subAsset.ParentID = 0
					err = skipHookDB.Save(subAsset).Error
					if err != nil {
						break
					}
				}
			}
		} else {
			rate := float64(interval) / float64(asset.Expire)
			asset.NetWorth = asset.Price.Mul(decimal.NewFromFloat(rate))
		}
	}

	return err
}

func getDiffDays(t1, t2 time.Time) int {
	timezone, _ := time.LoadLocation("Asia/Shanghai")
	timeDay1 := time.Date(t1.Year(), t1.Month(), t1.Day(), 0, 0, 0, 0, timezone)
	timeDay2 := time.Date(t2.Year(), t2.Month(), t2.Day(), 0, 0, 0, 0, timezone)

	return int(timeDay2.Sub(timeDay1).Hours() / 24)
}
