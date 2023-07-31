package controllers

import (
	"math"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
)

type pagingResult struct {
	Page      int `json:"page"`
	Limit     int `json:"limit"`
	PrevPage  int `json:"prevPage"`
	NextPage  int `json:"nextPage"`
	Count     int `json:"count"`
	TotalPage int `json:"totalPage"`
}

type pagination struct {
	ctx     *gin.Context
	query   *gorm.DB
	records interface{}
}

func (p *pagination) paginate() *pagingResult {
	// 1. get limit, page from query ?limit=10&page=2
	page, _ := strconv.Atoi(p.ctx.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(p.ctx.DefaultQuery("limit", "12"))

	// 2. count records
	ch := make(chan int)
	go p.countRecords(ch)

	// 3. find records
	// page => 1, 1-10, offset = 0
	// page => 2, 11-20, offset = 10
	offset := (page - 1) * limit
	p.query.Limit(limit).Offset(offset).Find(p.records)

	// 4. total pages
	count := <-ch
	totalPage := int(math.Ceil(float64(count) / float64(limit)))

	// 5. find next page
	var nextPage int
	if page == totalPage {
		nextPage = totalPage
	} else {
		nextPage = page - 1
	}

	// 6. create pagingResult
	return &pagingResult{
		Page:      page,
		Limit:     limit,
		Count:     count,
		PrevPage:  page - 1,
		NextPage:  nextPage,
		TotalPage: totalPage,
	}
}

func (p *pagination) countRecords(ch chan int) {
	var count int
	p.query.Model(p.records).Count(&count)

	ch <- count
}
