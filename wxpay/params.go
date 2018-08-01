package wxpay

import "strconv"

type Params map[string]string

// map本来已经是引用类型了，所以不需要 *Params
func (p Params) SetString(k, s string) {
	p[k] = s
}

func (p Params) GetString(k string) string {
	s, _ := p[k]
	return s
}

func (p Params) SetInt64(k string, i int64) {
	p[k] = strconv.FormatInt(i, 10)
}

func (p Params) GetInt64(k string) int64 {
	i, _ := strconv.ParseInt(p.GetString(k), 10, 64)
	return i
}
