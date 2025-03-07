package dao

import (
	"context"
	"user_growth/comm"
	"user_growth/dbhelper"
	"user_growth/models"
	"xorm.io/xorm"
)

type CoinDetailDao struct {
	db  *xorm.Engine
	ctx context.Context
}

func NewCoinDetailDao(ctx context.Context) *CoinDetailDao {
	return &CoinDetailDao{
		db:  dbhelper.GetDb(),
		ctx: ctx,
	}
}

func (dao *CoinDetailDao) Get(id int) (*models.TbCoinDetail, error) {
	data := &models.TbCoinDetail{}
	_, err := dao.db.ID(id).Get(data)

	return data, err
}

func (dao *CoinDetailDao) FindByUid(uid, page, size int) ([]models.TbCoinDetail, int64, error) {
	datalist := make([]models.TbCoinDetail, 0)
	sess := dao.db.Where("uid=?", uid)
	start := (page - 1) * size
	count, err := sess.Desc("id").Limit(size, start).FindAndCount(&datalist)
	return datalist, count, err
}

func (dao *CoinDetailDao) FindAllPage(page, size int) ([]models.TbCoinDetail, int64, error) {
	datalist := make([]models.TbCoinDetail, 0)
	if page < 1 {
		page = 1
	}
	if size < 1 {
		size = 100
	}
	start := (page - 1) * size
	count, err := dao.db.Desc("id").Limit(size, start).FindAndCount(&datalist)
	return datalist, count, err
}
func (dao *CoinDetailDao) Insert(data *models.TbCoinDetail) error {
	data.SysCreated = comm.Now()
	data.SysUpdated = comm.Now()

	_, err := dao.db.Insert(data)
	return err
}

func (dao *CoinDetailDao) update(data *models.TbCoinDetail, musColumns ...string) error {
	sess := dao.db.ID(data.Id)

	if len(musColumns) > 0 {
		sess.MustCols(musColumns...)
	}
	_, err := sess.Update(data)
	return err
}

func (dao *CoinDetailDao) Save(data *models.TbCoinDetail, musColumns ...string) error {
	if data.Id > 0 {
		return dao.update(data, musColumns...)
	} else {
		return dao.Insert(data)
	}
}
