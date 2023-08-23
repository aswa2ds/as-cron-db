package as_cron_db

import (
	"fmt"
	"sync"

	"github.com/aswa2ds/as-cron-db/cron_jobs"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type MysqlCfg struct {
	Username     string
	Password     string
	Address      string
	Port         string
	DatebaseName string
}

type Config struct {
	mysql MysqlCfg
}

var client clientSet

type Interface interface {
	CronJobs() cron_jobs.CronJobDao
}

type clientSet struct {
	db   *gorm.DB
	lock sync.Mutex
}

// CronJobs implements Interface.
func (c *clientSet) CronJobs() cron_jobs.CronJobDao {
	return cron_jobs.GetCronJobDao(c.db)
}

func BuildFromConfig(cfg Config) (Interface, error) {
	if client.db != nil {
		return &client, nil
	}

	client.lock.Lock()
	defer client.lock.Unlock()

	if client.db != nil {
		return &client, nil
	}

	dsn := fmt.Sprintf(
		"%s:%s@tcp(%s:%s)/%s?charset=utf8mb4&parseTime=True&loc=Local",
		cfg.mysql.Username,
		cfg.mysql.Password,
		cfg.mysql.Address,
		cfg.mysql.Port,
		cfg.mysql.DatebaseName,
	)
	var err error
	client.db, err = gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}
	return &client, nil
}
