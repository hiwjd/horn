package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"log"
	"net/http"
	"os"

	"strconv"

	"github.com/BurntSushi/toml"
	"github.com/hiwjd/horn/mysql"
	"github.com/hiwjd/horn/redis"
	"github.com/hiwjd/horn/state"
)

var (
	configPath string
	config     state.Config
)

func init() {
	flag.StringVar(&configPath, "c", "./state.toml", "配置文件的路径")
}

func main() {
	flag.Parse()

	if _, err := os.Stat(configPath); os.IsNotExist(err) {
		log.Printf("配置文件 %s 不存在, 使用默认配置 \r\n", configPath)
	} else {
		_, err := toml.DecodeFile(configPath, &config)
		if err != nil {
			log.Fatal(err)
		}
	}

	log.Println(config)

	redisManager := redis.New(config.RedisConfigs)
	mysqlManager := mysql.New(config.MysqlConfigs)
	stateService := state.New(mysqlManager, redisManager)

	http.HandleFunc("/ping", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
	})

	// POSTS 客服上线
	http.HandleFunc("/api/state/staff/online", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			oid, err := strconv.Atoi(r.FormValue("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}
			mid := r.FormValue("mid")
			sid := r.FormValue("sid")

			err = stateService.StaffOnline(oid, mid, sid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case http.MethodGet:
			// 获取在线客服列表
			q := r.URL.Query()
			oid, err := strconv.Atoi(q.Get("oid"))
			if err != nil {
				fmt.Fprintf(w, `{"code":%d,"msg":"%s"}`, http.StatusBadRequest, err.Error())
				return
			}

			staffs, err := stateService.OnlineStaffList(oid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				bs, err := json.Marshal(staffs)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
					return
				}

				w.Write(bs)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	// POST 客服下线
	http.HandleFunc("/api/state/staff/offline", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			oid, err := strconv.Atoi(r.FormValue("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}
			mid := r.FormValue("mid")
			sid := r.FormValue("sid")

			err = stateService.StaffOffline(oid, mid, sid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case http.MethodGet:
		// 获取下线客服列表
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	// POST 访客上线
	// GET 在线访客列表
	http.HandleFunc("/api/state/visitor/online", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			oid, err := strconv.Atoi(r.FormValue("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}
			mid := r.FormValue("mid")
			vid := r.FormValue("vid")

			err = stateService.VisitorOnline(oid, mid, vid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case http.MethodGet:
		// 获取在线访客列表
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	// POST 访客下线
	http.HandleFunc("/api/state/visitor/offline", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			oid, err := strconv.Atoi(r.FormValue("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}
			mid := r.FormValue("mid")
			vid := r.FormValue("vid")

			err = stateService.VisitorOffline(oid, mid, vid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				w.WriteHeader(http.StatusOK)
			}
		case http.MethodGet:
		// 获取下线访客列表
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	// POST 创建对话
	http.HandleFunc("/api/state/chat/create", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			oid, err := strconv.Atoi(r.FormValue("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}
			mid := r.FormValue("mid")
			cid := r.FormValue("cid")
			uid := r.FormValue("uid")

			err = stateService.CreateChat(oid, mid, cid, uid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				w.WriteHeader(http.StatusOK)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	// POST 客服/访客加入对话
	http.HandleFunc("/api/state/chat/join", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			oid, err := strconv.Atoi(r.FormValue("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}
			mid := r.FormValue("mid")
			cid := r.FormValue("cid")
			uid := r.FormValue("uid")

			err = stateService.JoinChat(oid, mid, cid, uid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				w.WriteHeader(http.StatusOK)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	// POST 客服/访客离开对话
	http.HandleFunc("/api/state/chat/leave", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodPost:
			oid, err := strconv.Atoi(r.FormValue("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}
			mid := r.FormValue("mid")
			cid := r.FormValue("cid")
			uid := r.FormValue("uid")

			err = stateService.LeaveChat(oid, mid, cid, uid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				w.WriteHeader(http.StatusOK)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	// GET 获取对话用户列表
	http.HandleFunc("/api/state/chat/uids", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			q := r.URL.Query()
			oid, err := strconv.Atoi(q.Get("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}
			cid := q.Get("cid")

			uids, err := stateService.GetUidsInChat(oid, cid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				bs, err := json.Marshal(uids)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
					return
				}

				w.Write(bs)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	// GET 根据用户获取对话ID
	http.HandleFunc("/api/state/user/cids", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			q := r.URL.Query()
			oid, err := strconv.Atoi(q.Get("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}
			uid := q.Get("uid")

			cids, err := stateService.GetChatIdsByUid(oid, uid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				bs, err := json.Marshal(cids)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
					return
				}

				w.Write(bs)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	// GET 根据用户获取推送地址
	http.HandleFunc("/api/state/user/pusher_addr", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			q := r.URL.Query()
			oid, err := strconv.Atoi(q.Get("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}
			uid := q.Get("uid")

			addr, err := stateService.GetPushAddrByUid(oid, uid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				bs, err := json.Marshal(addr)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
					return
				}

				w.Write(bs)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	// GET 根据组织ID获取所有的客服ID
	http.HandleFunc("/api/state/org/sids", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			q := r.URL.Query()
			oid, err := strconv.Atoi(q.Get("oid"))
			if err != nil {
				w.WriteHeader(http.StatusBadRequest)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
				return
			}

			sids, err := stateService.GetSidsInOrg(oid)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
			} else {
				bs, err := json.Marshal(sids)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					fmt.Fprintf(w, `{"error":"%s"}`, err.Error())
					return
				}

				w.Write(bs)
			}
		default:
			w.WriteHeader(http.StatusMethodNotAllowed)
			fmt.Fprintf(w, `{"error":"%s"}`, http.StatusText(http.StatusMethodNotAllowed))
		}
	})

	log.Printf("STATE[%s]启动，地址: %s \r\n", "X", ":9094")
	log.Fatal(http.ListenAndServe(":9094", nil))
}
