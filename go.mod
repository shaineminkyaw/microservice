module github.com/shaineminkyaw/microservice

go 1.17

require (
	// github.com/shaineminkyaw/microservice/pb v0.0.0-00010101000000-000000000000
	github.com/dgrijalva/jwt-go v3.2.0+incompatible
	github.com/golang/protobuf v1.5.2
	github.com/google/uuid v1.3.0
	github.com/mazen160/go-random v0.0.0-20210308102632-d2b501c85c03
	golang.org/x/crypto v0.0.0-20220722155217-630584e8d5aa
	google.golang.org/grpc v1.48.0
	google.golang.org/protobuf v1.28.1
	gopkg.in/ini.v1 v1.67.0
	gorm.io/driver/mysql v1.3.5
	gorm.io/gorm v1.23.8
)

require (
	github.com/go-sql-driver/mysql v1.6.0 // indirect
	github.com/google/go-cmp v0.5.8 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	github.com/stretchr/testify v1.8.0 // indirect
	golang.org/x/net v0.0.0-20220624214902-1bab6f366d9e // indirect
	golang.org/x/sys v0.0.0-20220610221304-9f5ed59c137d // indirect
	golang.org/x/text v0.3.7 // indirect
	google.golang.org/genproto v0.0.0-20220805133916-01dd62135a58 // indirect
)


replace "github.com/shaineminkyaw/microservice/pb" => ./pb