package lib

import (
	"database/sql"
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

func (tp *Tips) NewTips(connstr string) error {
	var err error
	tp.db, err = sql.Open("mysql", connstr)
	if err != nil {
		return err
	}
	tp.db.SetMaxIdleConns(2)
	tp.db.SetMaxOpenConns(5)

	err = tp.db.Ping()
	log.Printf("%#v", tp)
	if err != nil {
		return err
	}
	return err
}

func (tp *Tips) GetMessage() string {
	log.Printf("%#v", tp)
	rows, err := tp.db.Query("select a.id,a.addtime,b.username from rd_account_cash as a left join rd_user as b on a.user_id=b.user_id order by a.id desc limit 10")
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

	switch {
	case newcom > 1:
		s := strconv.Itoa(newcom)
		returnstr := "有X单提现申请等待处理！"
		return strings.Replace(returnstr, "X", s, 1)
	case newcom == 1:
		returnstr := "客户：N有新的提现申请！"
		returnstr = strings.Replace(returnstr, "N", lastusername, 1)
		return returnstr
	default:
		return ""
	}

}
