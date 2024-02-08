module ivpn.net/email-service

go 1.19

require github.com/mhale/smtpd v0.8.2

require github.com/joho/godotenv v1.5.1

require (
	github.com/satori/go.uuid v1.2.0
	gopkg.in/gomail.v2 v2.0.0-20160411212932-81ebce5c23df
	gorm.io/driver/mysql v1.5.4
	gorm.io/gorm v1.25.7
)

require (
	github.com/go-sql-driver/mysql v1.7.1 // indirect
	github.com/jinzhu/inflection v1.0.0 // indirect
	github.com/jinzhu/now v1.1.5 // indirect
	gopkg.in/alexcesaro/quotedprintable.v3 v3.0.0-20150716171945-2caba252f4dc // indirect
	gopkg.in/check.v1 v1.0.0-20201130134442-10cb98267c6c // indirect
)
