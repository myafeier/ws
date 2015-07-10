package lib

import (
	"database/sql"
	"encoding/json"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"strconv"
	"strings"
	"time"
)

type Tips struct {
	db     *sql.DB
	sended map[int]int64
}

type mess struct {
	Id   int64  `json:"id"`
	Mess string `json:"mess"`
}

func (tp *Tips) NewTips(connstr string) error {
	var err error
	tp.db, err = sql.Open("mysql", connstr)
	if err != nil {
		return err
	}
	tp.db.SetMaxIdleConns(2)
	tp.db.SetMaxOpenConns(5)

	err = tp.db.Ping()
	if err != nil {
		return err
	}
	return err
}

func (tp *Tips) GetHotMessage() string {
	lstime := time.Now().Unix() - 24*3600
	var num int
	err := tp.db.QueryRow("select count(1) as num from rd_account_cash where status=0 and addtime>?", lstime).Scan(&num)
	// defer db.Close()
	if err != nil {
		log.Println("sql query error:", err)
		return ""
	}
	if num > 0 {
		returnstr := "欢迎回来，当前有 X 单提现申请等待处理！"
		returnstr = strings.Replace(returnstr, "X", strconv.Itoa(num), 1)
		messnew := &mess{Id: lstime, Mess: returnstr}
		return_json, err := json.Marshal(messnew)
		if err != nil {
			log.Println("err:", err)
		}
		// log.Println(return_json)
		if err == nil {
			return string(return_json)
		} else {
			return ""
		}

	} else {
		return ""
	}
}
func (tp *Tips) GetMessage() string {
	rows, err := tp.db.Query("select a.id,a.addtime,b.username from rd_account_cash as a left join rd_user as b on a.user_id=b.user_id where a.status=0 order by a.id desc limit 10")
	defer rows.Close()
	if err != nil {
		log.Println("sql_query error:", err)
		return ""
	}
	newcom := 0
	var id int
	var username string
	var addtime int64
	var issended bool
	var lastusername string
	if tp.sended == nil {
		tp.sended = make(map[int]int64)
	}
	oldnums := len(tp.sended)
	now := time.Now().Unix()
	if oldnums > 0 { //把过期的信息删除掉
		for k, v := range tp.sended {
			if (now - v) > 3600 {
				delete(tp.sended, k)
			}
		}
	}

	for rows.Next() {
		err := rows.Scan(&id, &addtime, &username)
		if err != nil {
			log.Println("query result error:", err)
			return ""
		}
		if (now - addtime) > 3600 {
			continue
		}
		if oldnums == 0 {
			tp.sended[id] = now
			oldnums++
			newcom++
		} else {

			for oid, _ := range tp.sended {
				if oid == id {
					issended = true
					break

				}
			}
			if !issended {
				lastusername = username
				tp.sended[id] = now
				oldnums++
				newcom++
			}

		}
	}
	var returnstr string
	switch {
	case newcom > 1:
		s := strconv.Itoa(newcom)
		returnstr = "有 X 单提现申请等待处理！"
		returnstr = strings.Replace(returnstr, "X", s, 1)
	case newcom == 1:
		returnstr = "N 有新的提现申请！"
		returnstr = strings.Replace(returnstr, "N", lastusername, 1)
	default:
		return ""
	}
	rmsg := &mess{Id: now, Mess: returnstr}
	return_json, err := json.Marshal(rmsg)
	if err != nil {
		log.Println("err:", err)
	}
	// log.Println(return_json)
	if err == nil {
		return string(return_json)
	} else {
		return ""
	}

}
