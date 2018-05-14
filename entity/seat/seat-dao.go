/***************************************************************************
Copyright(C) 2018
Author: huziang
Description: seat的dao层，使用数据库层接口，为上层逻辑层提供接口
	PS：该文件所有错误全都panic抛出，不进行错误返回
Date: 2018年5月4日 星期五 上午10:17
****************************************************************************/

package seat

import (
	"fmt"

	"labix.org/v2/mgo"
	"labix.org/v2/mgo/bson"

	"github.com/book-library-seat-system/go-server/mgdb"
	. "github.com/book-library-seat-system/go-server/util"
)

// seat 使用同一个数据库，但是使用不同的表
var database *mgo.Database

func init() {
	database = mgdb.Mydb.DB("seat")
	fmt.Println("Seat database init!")
}

// TItemsAtomicService 一个空类型
type TItemsAtomicService struct{}

// service 空类型的指针，使用函数
var service = TItemsAtomicService{}

/*************************************************
Function: Insert
Description: 保存TItems信息到数据库
InputParameter:
	titems: 要保存的TItems的信息
Return: none
*************************************************/
func (*TItemsAtomicService) Insert(stitem *STItem) {
	c := database.C(stitem.School)
	for i := 0; i < len(stitem.Titems); i++ {
		err := c.Insert(stitem.Titems[i])
		CheckNewErr(err, "101|数据库座位信息插入出现错误")
	}
}

/*************************************************
Function: FindBySchool
Description: 通过主键School查询数据
InputParameter:
	school: 主键
Return: 找到的对应TItems，如果没有为nil
*************************************************/
func (*TItemsAtomicService) FindBySchool(school string) []TItem {
	c := database.C(school)
	titems := []TItem{}
	err := c.Find(nil).All(&titems)
	CheckNewErr(err, "103|数据库座位信息查找出现错误")
	return titems
}

/*************************************************
Function: FindBySchoolAndTimeInterval
Description: 通过两个主键查询数据
InputParameter:
	school: 主键1
	timeinterval: 主键2
Return: 查找到的座位信息，如果不存在返回nil
*************************************************/
func (this *TItemsAtomicService) FindBySchoolAndTimeInterval(school string, timeinterval TimeInterval) []Item {
	c := database.C(school)
	titem := TItem{}
	err := c.Find(bson.M{"timeinterval": timeinterval}).One(&titem)
	CheckNewErr(err, "103|数据库座位信息查找出现错误")
	return titem.Items
}

/*************************************************
Function: UpdateAllSeat
Description: 通过两个主键更新多个座位信息
InputParameter:
	school: 主键1
	timeinterval: 主键2
	seats: 要更改的座位信息
Return: none
*************************************************/
func (*TItemsAtomicService) UpdateAllSeat(
	school string,
	timeinterval TimeInterval,
	seats []Item) {
	c := database.C(school)
	err := c.Update(
		bson.M{"timeinterval": timeinterval},
		bson.M{"$set": bson.M{"items": seats}},
	)
	CheckNewErr(err, "102|数据库座位信息更新出现错误")
}

/*************************************************
Function: UpdateOneSeat
Description: 通过两个主键更新单个座位信息
InputParameter:
	school: 主键1
	timeinterval: 主键2
	seat: 要更改的座位信息
Return: none
*************************************************/
func (*TItemsAtomicService) UpdateOneSeat(
	school string,
	timeinterval TimeInterval,
	seat Item) {
	c := database.C(school)
	err := c.Update(
		bson.M{
			"timeinterval": timeinterval,
			"items":        bson.M{"$elemMatch": bson.M{"seatid": seat.SeatID}},
		},
		bson.M{
			"$set": bson.M{
				"items.$.seatinfo":  seat.Seatinfo,
				"items.$.studentid": seat.StudentID,
			},
		},
	)
	CheckNewErr(err, "102|数据库座位信息更新出现错误")
}

/*************************************************
Function: DeleteBySchool
Description: 通过主键School删除数据
InputParameter:
	school: 主键1
Return: none
*************************************************/
func (*TItemsAtomicService) DeleteBySchool(school string) {
	err := database.C(school).DropCollection()
	CheckNewErr(err, "104|数据库座位信息删除出现错误")
}

/*************************************************
Function: DeleteBySchoolAndTimeInterval
Description: 通过两个主键删除数据（删除过时信息）
InputParameter:
	school: 主键1
	timeinterval: 主键2
Return: none
*************************************************/
func (*TItemsAtomicService) DeleteBySchoolAndTimeInterval(school string, timeinterval TimeInterval) {
	c := database.C(school)
	err := c.Remove(bson.M{"timeinterval": timeinterval})
	CheckNewErr(err, "104|数据库座位信息删除出现错误")
}
