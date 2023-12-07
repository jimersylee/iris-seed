package main

import (
	"fmt"
	"time"
)

// snowflake算法
// 0-00000000000000000000000000000000000000000-00000-00000-000000000000
// 1bit不用     41bit毫秒时间戳                  5位数据中心id 5位机器id 12位序列号
// 41bit时间戳最大值为2^41-1,即2199023255551,可用到2039-09-07  5位数据中心id 2^5,0-31最大支持32个,机器id也是最大32个,序列号12位,2^12,最大数为4095,理论上计算,就算是同一个机器,同一毫秒内,序列号可以生成4096个,够用

const workIdShift int32 = 12
const dataCenterIdShift int32 = 17
const timestampShift int32 = 22

// 序号占的位数,12位,2^12,最大数为4095
const sequenceBit = 12

type Snowflake struct {
	timestamp     int64
	workerId      int64
	dataCenterId  int64
	sequence      int64
	lastTimestamp int64
	sequenceMask  int64
}

func NewSnowflake(workerId int64, dataCenterId int64) *Snowflake {
	return &Snowflake{
		timestamp:    1414213562373, //这个值可以任选一个
		workerId:     workerId,
		dataCenterId: dataCenterId,
		sequence:     0,
		sequenceMask: 4095, //2^12-1
	}
}

func (sf *Snowflake) NextId() (int64, error) {
	timeNow := time.Now().UnixMilli()
	if timeNow < sf.lastTimestamp {
		return 0, fmt.Errorf("clock moved backwards")
	}
	if timeNow == sf.lastTimestamp {
		sf.sequence = (sf.sequence + 1) & sf.sequenceMask
		if sf.sequence == 0 {
			//同意毫秒内已经达到最大序号4095
			timeNow = timeNow + 1
		}
	} else {
		sf.sequence = 0
	}
	sf.lastTimestamp = timeNow
	return (timeNow-sf.timestamp)<<timestampShift | sf.dataCenterId<<dataCenterIdShift | sf.workerId<<workIdShift | sf.sequence, nil

}

func main() {
	sf := NewSnowflake(1, 1)
	for i := 0; i < 100; i++ {
		id, err := sf.NextId()
		if err != nil {
			fmt.Println(err)
		} else {
			fmt.Println(id)
		}
	}
}
