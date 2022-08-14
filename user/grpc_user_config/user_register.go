package grpc_user_config

import (
	"context"
	"log"
	"time"

	"github.com/mazen160/go-random"
	"github.com/shaineminkyaw/microservice/pb"
	"github.com/shaineminkyaw/microservice/user/model"
	"github.com/shaineminkyaw/microservice/user/util"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"gorm.io/gorm"
)

func (server *Server) UserRegister(ctx context.Context, req *pb.UserRequest) (*pb.UserResponse, error) {
	//

	vUser := &model.VerifyCode{}
	//@@@validate verifycode
	err := server.Database.Sql.Model(&model.VerifyCode{}).Where("email = ?", req.GetEmail()).First(&vUser).Error
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error on validate verify code ")
	}
	if vUser == nil {
		return nil, status.Errorf(codes.NotFound, "email not registered for verify code !")
	}
	if req.GetVerifyCode() != vUser.Code || time.Now().Unix() > vUser.ExpireTime.UnixNano() {
		return nil, status.Errorf(codes.Internal, "verifycode invalid")
	}

	//hash password
	userPassword, err := util.HashPassword(req.GetPassword())
	log.Println("User Password ....", userPassword)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "error on hash password")
	}
	charset := "12345678"
	length := 8
	usr, _ := random.Random(length, charset, true)
	userName := "U" + usr
	currency := "USD"
	bankCard, err := util.GetBankCardNumber(req.City)
	if err != nil {
		return nil, status.Errorf(codes.Internal, "erorr on generate bank card")
	}
	user := &model.User{
		Username:       userName,
		Password:       userPassword,
		Email:          req.GetEmail(),
		RegesiterIP:    "",
		LastLoginIP:    "",
		NationID:       req.GetNationId(),
		BankCardNumber: bankCard,
		City:           req.GetCity(),
		Balance:        100,
		Currency:       currency,
		Type:           int8(req.GetGenderType()),
	}

	// if ctx.Err() == context.Canceled {
	// 	log.Printf("request is canceled")
	// 	return nil, status.Errorf(codes.Aborted, "request is Cancel")
	// }
	// // if ctx.Err() == context.Background().Err() {
	// // 	log.Printf("request is canceled")
	// // 	return nil, status.Errorf(codes.Aborted, "request is Cancel!!!!")
	// // }
	fUser := &model.User{}
	data := server.Database.Sql.Model(&model.User{}).Where("email =?", req.GetEmail())
	err = data.First(&fUser).Error
	if err == gorm.ErrRecordNotFound {
		err := server.Database.Sql.Model(&model.User{}).Create(&user).Error
		if err != nil {
			return nil, status.Errorf(codes.Internal, "erorr on user create !")
		}
		err = util.SaveUserBankCard(req.GetCity(), user.ID, bankCard)
		if err != nil {
			return nil, status.Errorf(codes.Internal, "erorr on store bank card ")
		}
	} else {
		return nil, status.Errorf(codes.Internal, "user already exists")
	}

	resp := &pb.UserResponse{
		User: convert(user),
	}

	return resp, nil
}
