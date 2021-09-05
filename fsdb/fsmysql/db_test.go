package fsmysql

import (
	"fmt"
	"testing"

	"fsky.pro/fstest"
)

func TestGetObjTagMembers(t *testing.T) {
	fstest.PrintTestBegin("dbkeyMapValues")
	defer fstest.PrintTestEnd()

	type A struct {
		V1 string  `db:"v1"`
		v2 int     `db:"v2"`
		v3 float32 `mysql:"v3"`
		v4 int64   `mysql:"-"`
	}

	a := new(A)
	kv, err := dbkeyMapValues(a, "V1,v3,v4", true)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	*(kv["v1"].(*string)) = "xxxx"
	*(kv["v3"].(*float32)) = 1.1
	fmt.Println(a)

	kv, err = dbkeyMapValues(a, "", true)
	if err != nil {
		fmt.Println("error:", err.Error())
		return
	}
	*(kv["v2"].(*int)) = 100
	fmt.Println(a)
}

type Player struct {
	Uid                  int
	HeadPath             string //头像路径
	ChannelId            string //渠道号
	UserName             string //玩家名称
	UserLevel            int    //玩家等级
	UserExp              int    //玩家经验
	Gold                 int    //玩家金币
	Physical             int    //玩家体力
	MaxPhysical          int    //最大体力
	Diamond              int    //玩家钻石
	VIPLevel             int    //VIP等级
	Sex                  uint8  //性别
	LastUpdateTime       int64  //上次更新时间
	PhysicalTimer        int    //体力倒计时
	PhysicalTimeInitTime int64  //增加体力的时间
	ReNewTime            int    //刷新体力时间
	StartGameTime        int    //进入游戏时间点
	GameID               int    //游戏ID
	HangUpID             int    //挂机ID
	HardProgress         int    //困难模式进度

	HangUpTime  int64 //挂机时间
	ReceiveTime int64 //领取挂机奖励的时间
	Recharge    int   //充值金额
	Energy      int   //能量水晶
	FriendShip  int   //友情点
	//Atk         int   //玩家战力

	PysicalCDTime    int64 //免费获取体力时间
	GetPysicalNum    int   //获取体力次数
	CurrentChatperID int   //当前关卡进度 用于增加玩家经验
	RewardEndless    int   // 已经领取的无尽奖励
	Endless          int   //无尽关卡
	EndlessNumDay    int   //今天进入无尽模式的次数
	LassEndlessTime  int64 //最后一次进入无尽模式的时间
	LngType          int8  //语言类型

	Active           int   //活动副本
	ActiveNumDay     int   //今天进入活动副本的次数
	CreatePlayerTime int64 //创建角色时间
	LassActiveTime   int64 //最后一次进入活动副本的时间

	ZombiesNumDay   int   //今天玩僵尸局的次数
	LastZombiesTime int64 //最后一次进入僵尸局的时间

	BossFightNumDay   int   //今天玩boss战的次数
	LastBossFightTime int64 //最后一次进入boss战的时间

	FirstLogin     int
	LassAccessTime int64  //最后一次登录时间
	OnlineTime     int    //在线时长
	Age            int    //年龄
	State          uint8  //账号状态
	HasGuideId     string //新手引导
	SiteID         int    //当前正在进行的场次
	WearIndex      string //已穿戴的武器
	Remark         string
	SignDay        int //签到时间
}

func TestQuery(t *testing.T) {
	fstest.PrintTestBegin("Query")
	defer fstest.PrintTestEnd()

	player := new(Player)
	info := &S_DBInfo{
		Host:     "localhost",
		Port:     3306,
		User:     "root",
		Password: "sz123456",
		DBName:   "tzbb_1",
	}

	db, err := Open(info)
	if err != nil {
		fmt.Printf("open mysql database fail. error: %v\n", err)
		return
	}
	fmt.Println("database has been opened.")

	result := db.QueryObjects(player, "UserName,UserLevel,Gold", "`player`", "LIMIT ?", 10)
	if result.Err() != nil {
		fmt.Printf("load player from db fail, db error: %v. sqltx:\n%s.\n", result.Err(), result.SQLText())
		return
	}

	result.ForEach(func(obj interface{}, err error) bool {
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(obj)
		}
		return true
	})
}
