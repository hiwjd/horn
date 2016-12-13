package remote

import (
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"

	"github.com/hiwjd/horn/state"

	"strconv"
)

type remoteState struct {
	apihost string
}

func New(apihost string) state.State {
	return &remoteState{
		apihost: apihost,
	}
}

// 客服上线
func (c *remoteState) StaffOnline(oid int, mid string, sid string) error {
	path := fmt.Sprintf("%s/api/state/staff/online", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"mid": {mid},
		"sid": {sid},
	}

	return post1(path, values)
}

// 客服下线
func (c *remoteState) StaffOffline(oid int, mid string, sid string) error {
	path := fmt.Sprintf("%s/api/state/staff/offline", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"mid": {mid},
		"sid": {sid},
	}

	return post1(path, values)
}

// 访客上线
func (c *remoteState) VisitorOnline(oid int, mid string, vid string) error {
	path := fmt.Sprintf("%s/api/state/visitor/online", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"mid": {mid},
		"vid": {vid},
	}

	return post1(path, values)
}

// 访客上线
func (c *remoteState) VisitorOffline(oid int, mid string, vid string) error {
	path := fmt.Sprintf("%s/api/state/visitor/offline", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"mid": {mid},
		"vid": {vid},
	}

	return post1(path, values)
}

// 创建对话
func (c *remoteState) CreateChat(oid int, mid string, cid, creator, sid, vid, tid string) error {
	path := fmt.Sprintf("%s/api/state/chat/create", c.apihost)
	values := url.Values{
		"oid":     {strconv.Itoa(oid)},
		"mid":     {mid},
		"cid":     {cid},
		"creator": {creator},
		"sid":     {sid},
		"vid":     {vid},
		"tid":     {tid},
	}

	return post1(path, values)
}

// 加入对话
func (c *remoteState) JoinChat(oid int, mid string, cid string, uid string) error {
	path := fmt.Sprintf("%s/api/state/chat/join", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"mid": {mid},
		"cid": {cid},
		"uid": {uid},
	}

	return post1(path, values)
}

// 离开对话
func (c *remoteState) LeaveChat(oid int, mid string, cid string, uid string) error {
	path := fmt.Sprintf("%s/api/state/chat/leave", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"mid": {mid},
		"cid": {cid},
		"uid": {uid},
	}

	return post1(path, values)
}

// 获取对话中的用户ID
func (c *remoteState) GetUidsInChat(oid int, cid string) ([]string, error) {
	path := fmt.Sprintf("%s/api/state/chat/uids", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"cid": {cid},
	}

	path = fmt.Sprintf("%s?%s", path, values.Encode())
	resp, err := http.Get(path)

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var r []string
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return r, nil
	} else {
		var r struct {
			Error string
		}
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(r.Error)
	}
}

func (c *remoteState) OnlineStaffList(oid int) ([]*state.Staff, error) {
	path := fmt.Sprintf("%s/api/state/staff/online", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
	}

	path = fmt.Sprintf("%s?%s", path, values.Encode())
	resp, err := http.Get(path)

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var r []*state.Staff
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return r, nil
	} else {
		var r struct {
			Error string
		}
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(r.Error)
	}
}

func (c *remoteState) GetChatIdsByUid(oid int, uid string) ([]string, error) {
	path := fmt.Sprintf("%s/api/state/user/cids", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"uid": {uid},
	}

	path = fmt.Sprintf("%s?%s", path, values.Encode())
	resp, err := http.Get(path)

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var r []string
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return r, nil
	} else {
		var r struct {
			Error string
		}
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(r.Error)
	}
}

func (c *remoteState) GetPushAddrByUid(oid int, uid string) (string, error) {
	path := fmt.Sprintf("%s/api/state/user/pusher_addr", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"uid": {uid},
	}

	path = fmt.Sprintf("%s?%s", path, values.Encode())
	resp, err := http.Get(path)

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	if resp.StatusCode == http.StatusOK {
		var r string
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return "", err
		}

		return r, nil
	} else {
		var r struct {
			Error string
		}
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return "", err
		}

		return "", errors.New(r.Error)
	}
}

func (c *remoteState) GetSidsInOrg(oid int) ([]string, error) {
	path := fmt.Sprintf("%s/api/state/org/sids", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
	}

	path = fmt.Sprintf("%s?%s", path, values.Encode())
	resp, err := http.Get(path)

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var r []string
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return r, nil
	} else {
		var r struct {
			Error string
		}
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(r.Error)
	}
}

func (c *remoteState) GetVisitor(oid int, vid string) (*state.Visitor, error) {
	path := fmt.Sprintf("%s/api/state/visitor/info", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"vid": {vid},
	}

	path = fmt.Sprintf("%s?%s", path, values.Encode())
	resp, err := http.Get(path)

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var r state.Visitor
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return &r, nil
	} else {
		var r struct {
			Error string
		}
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(r.Error)
	}
}

func (c *remoteState) GetStaff(oid int, sid string) (*state.Staff, error) {
	path := fmt.Sprintf("%s/api/state/staff/info", c.apihost)
	values := url.Values{
		"oid": {strconv.Itoa(oid)},
		"sid": {sid},
	}

	path = fmt.Sprintf("%s?%s", path, values.Encode())
	resp, err := http.Get(path)

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var r state.Staff
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return &r, nil
	} else {
		var r struct {
			Error string
		}
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(r.Error)
	}
}

func (c *remoteState) GetVisitorLastTracks(oid int, vid string, limit int) ([]*state.Track, error) {
	path := fmt.Sprintf("%s/api/state/visitor/tracks/last", c.apihost)
	values := url.Values{
		"oid":   {strconv.Itoa(oid)},
		"vid":   {vid},
		"limit": {strconv.Itoa(limit)},
	}

	path = fmt.Sprintf("%s?%s", path, values.Encode())
	resp, err := http.Get(path)

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		var r []*state.Track
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return r, nil
	} else {
		var r struct {
			Error string
		}
		err = json.Unmarshal(bs, &r)
		if err != nil {
			return nil, err
		}

		return nil, errors.New(r.Error)
	}
}

func post1(path string, values url.Values) error {
	resp, err := http.PostForm(path, values)
	if err != nil {
		return err
	}

	if resp.StatusCode == http.StatusOK {
		return nil
	}

	defer resp.Body.Close()
	bs, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return err
	}

	var r struct {
		Error string
	}
	err = json.Unmarshal(bs, &r)
	if err != nil {
		return err
	}

	return errors.New(r.Error)
}
