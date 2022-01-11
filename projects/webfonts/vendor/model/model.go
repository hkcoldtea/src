package model

import (
	"log"

	"inventory"
)

const (
	fontsdirectory = "./fonts"
)

var (
	fontInventory *inventory.Inventory = inventory.New()
)
/*
func getEnv(key, fallback string) string {
	if value, ok := os.LookupEnv(key); ok {
		return value
	}
	return fallback
}
*/
// Note: You can customize the value by desire.
// export DB_HOST=localhost
// export DB_PORT=5432
// export DB_USER=hugo
// export DB_PASSWORD=postgres
// export DB_DATABASE=project.db
/*
// InitialMigration for project with db.AutoMigrate
func InitialMigration() {
	defer func() {
		err := recover()
		if err != nil {
			panic(err)
		}
	}()

//	host := getEnv("DB_HOST", "localhost")
//	port := getEnv("DB_PORT", "5432")
//	user := getEnv("DB_USER", "hugo")
//	password := getEnv("DB_PASSWORD", "postgres")
	database := getEnv("DB_DATABASE", "/tmp/project.db")

	var err error
	db, err = gorm.Open(sqlite.Open("file:"+database+"?mode=memory&cache=shared&_fk=1"), &gorm.Config{
		PrepareStmt: true,
		NamingStrategy: schema.NamingStrategy{
			TablePrefix: "web.",
		},
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				IgnoreRecordNotFoundError: true,
				Colorful:      false,
			},
		),
	})
	if err != nil {
		panic("Failed to connect to database")
	}
}
*/
/*
//GetDB ...
func GetDB() *gorm.DB {
	return db
}
*/
// Build font inventory
func BuildFontInventory() {
	if err := fontInventory.Build(fontsdirectory); err != nil {
		log.Println(err.Error())
	}
	if fontInventory.Len() == 0 {
		log.Println("Empty font library")
	}
}

func GetFontInventory() *inventory.Inventory {
	return fontInventory
}
