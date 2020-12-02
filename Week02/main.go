package main

import (
	"database/sql"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
)

type svc struct {
	dbsts *Db
}

type Db struct {
	db_con *gorm.DB
}

func newDb(db *gorm.DB) *Db {
	return &Db{db}
}

func newSvc() svc {
	return svc{newDb(dbeng)}
}

var dbeng *gorm.DB

type production struct {
	pid     uint
	name    string
	storage uint
}

func (db *Db) queryProductionName(id uint) (*production, error) {
	return nil, errors.Wrapf(sql.ErrNoRows, "not found id = %d production name\n", id)
}

func (db *Db) queryProductionStorage(id uint) (*production, error) {
	return nil, errors.Wrapf(sql.ErrNoRows, "not found id = %d production storage\n", id)
}

//svc定义pass id是唯一确认产品，名字可以为空
func (s *svc) getProName(id uint) (*production, error) {
	p, err := s.dbsts.queryProductionName(id)
	if errors.Is(err, sql.ErrNoRows) {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}
	return p, err
}

//svc定义 error 存量不能为空
func (s *svc) getProStorage(id uint) (*production, error) {
	if p, err := s.dbsts.queryProductionStorage(id); err != nil {
		return nil, err
	} else {
		return p, nil
	}

}

func main() {
	s := newSvc()
	_, err := s.getProStorage(123)
	if errors.Is(err, sql.ErrNoRows) {
		fmt.Printf("404 not found storage %+v\n", err)
	} else {
		fmt.Printf("5xx server err %+v\n", err)
	}

	_, err = s.getProName(123)
	if err != nil {
		fmt.Printf("5xx %+v\n", err)
	}
	fmt.Println("ok")
}
