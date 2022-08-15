package util

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/shaineminkyaw/microservice/user/ds"
	"github.com/shaineminkyaw/microservice/user/model"
	"golang.org/x/crypto/scrypt"
	"gorm.io/gorm"
)

//@@@hash password
func HashPassword(pass string) (string, error) {
	//
	salt := make([]byte, 64)
	_, err := rand.Read(salt)
	if err != nil {
		return "", nil
	}

	unHash, err := scrypt.Key([]byte(pass), salt, 32768, 8, 1, 32)
	if err != nil {
		log.Fatalf("error on crypting password %v", err.Error())
	}
	hashed := fmt.Sprintf("%v.%v", hex.EncodeToString(unHash), hex.EncodeToString(salt))
	return hashed, nil

}

//validate hash
func ValidateHashedPassword(hash, plainPassword string) (bool, error) {
	//
	splitHash := strings.Split(hash, ".")
	salt, err := hex.DecodeString(splitHash[1])
	if err != nil {
		log.Fatalf("error on : %v", err.Error())
	}
	plainHashed, err := scrypt.Key([]byte(plainPassword), salt, 32768, 8, 1, 32)
	if err != nil {
		log.Fatalf("error on %v", err.Error())
	}

	return (hex.EncodeToString(plainHashed)) == (splitHash[0]), nil
}

const (
	Yangon    = "0001"
	Mandalay  = "0002"
	Naypyitaw = "0003"
	Taunggyi  = "0004"
)

func GetBankCardNumber(city string) (string, error) {
	//
	var allCity int64
	var c, cStr, bankStr string
	info := &model.UserCityTotalBankCard{}

	tx := ds.Auth_DB.Begin()
	db := tx.Model(&model.UserCityTotalBankCard{}).Where("id = ?", 1)
	err := db.First(&info).Error
	if err == gorm.ErrRecordNotFound {
		totalCard := &model.UserCityTotalBankCard{}
		switch city {
		case "yangon":
			allCity += 1
			c = Yangon
			totalCard = &model.UserCityTotalBankCard{
				YangonCard: allCity,
			}
		case "mandalay":
			allCity += 1
			c = Mandalay
			totalCard = &model.UserCityTotalBankCard{
				MandalayCard: allCity,
			}
		case "naypyitaw":
			allCity += 1
			c = Naypyitaw
			totalCard = &model.UserCityTotalBankCard{
				NaypyitawCard: allCity,
			}
		case "taunggyi":
			allCity += 1
			c = Taunggyi
			totalCard = &model.UserCityTotalBankCard{
				TaunggyiCard: allCity,
			}
		default:
			allCity = 0
			c = "Unknown"
		}

		err = tx.Model(&model.UserCityTotalBankCard{}).Create(&totalCard).Error
		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
		}
		if len(city) > 0 {
			cStr = strconv.Itoa(int(allCity))
			length := len(cStr)
			for i := length; i < 8; i++ {
				bankStr += "0"
			}
		}

	} else {
		var aStr int64
		updateBankCard := &model.UserCityTotalBankCard{}
		if city == "yangon" {
			aStr = info.YangonCard + 1
			c = Yangon
			updateBankCard = &model.UserCityTotalBankCard{
				YangonCard: aStr,
			}
		} else if city == "mandalay" {
			aStr = info.MandalayCard + 1
			c = Mandalay
			updateBankCard = &model.UserCityTotalBankCard{
				MandalayCard: aStr,
			}
		} else if city == "naypyitaw" {
			aStr = info.NaypyitawCard + 1
			c = Naypyitaw
			updateBankCard = &model.UserCityTotalBankCard{
				NaypyitawCard: aStr,
			}
		} else if city == "taunggyi" {
			aStr = info.TaunggyiCard + 1
			c = Taunggyi
			updateBankCard = &model.UserCityTotalBankCard{
				TaunggyiCard: aStr,
			}
		}

		db := tx.Model(&model.UserCityTotalBankCard{}).Where("id = ?", 1)
		err = db.Updates(&updateBankCard).Error
		if err != nil {
			log.Println(err.Error())
			tx.Rollback()
		}

		if len(city) > 0 {
			cStr = strconv.Itoa(int(aStr))
			length := len(cStr)
			for i := length; i < 8; i++ {
				bankStr += "0"
			}
		}

	}
	if err != nil && err != gorm.ErrRecordNotFound {
		log.Println(err.Error())
	}
	err = tx.Commit().Error
	if err != nil {
		log.Println(err.Error())
		tx.Rollback()
	}

	return c + bankStr + cStr, nil
}

func SaveUserBankCard(city string, userId uint64, bankcard string) error {
	//

	card := &model.UserBankCard{
		Uid:           userId,
		BanCardNumber: bankcard,
	}
	err := ds.Auth_DB.Model(&model.UserBankCard{}).Create(&card).Error
	if err != nil {
		log.Println(err.Error())
	}

	return nil
}
