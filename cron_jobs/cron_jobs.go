package cron_jobs

import (
	"time"

	"gorm.io/gorm"
)

type CronJob struct {
	gorm.Model
	JobName        string    `gorm:"not null;type:varchar(255);comment:任务名"`
	Protocal       int8      `gorm:"not null;type:tinyint;comment:协议类型"`
	Address        string    `gorm:"not null;type:varchar(20);comment:远程调用地址"`
	Port           string    `gorm:"not null;type:varchar(10);comment:远程调用端口"`
	Path           string    `gorm:"not null;type:text;comment:远程调用URL"`
	CronExpression string    `gorm:"not null;type:varchar(50);comment:cron 表达式"`
	NextToggleTime time.Time `gorm:"not null;comment:下次触发时间"`
}

type CronJobDao interface {
	List() []CronJob
	Insert(cronJob CronJob)
	UpdateNextToggleTime(cronJob CronJob)
}

func GetCronJobDao(db *gorm.DB) CronJobDao {
	cronJobDao := &cronJobMysqlDao{
		db: db,
	}
	cronJobDao.db.AutoMigrate(&CronJob{})
	return cronJobDao
}

type cronJobMysqlDao struct {
	db *gorm.DB
}

// Insert implements CronJobDao.
func (dao *cronJobMysqlDao) Insert(cronJob CronJob) {
	dao.db.Create(&cronJob)
}

// List implements CronJobDao.
func (dao *cronJobMysqlDao) List() []CronJob {
	var cronJobs []CronJob
	dao.db.Find(&cronJobs)
	return cronJobs
}

// UpdateNextToggleTime implements CronJobDao.
func (dao *cronJobMysqlDao) UpdateNextToggleTime(cronJob CronJob) {
	dao.db.Updates(cronJob)
}

var _ CronJobDao = &cronJobMysqlDao{}
