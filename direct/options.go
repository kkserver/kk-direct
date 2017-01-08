package direct

type Options map[interface{}]interface{}

func (O Options) Name() string {
	v, ok := O["name"]
	if ok {
		vv, ok := v.(string)
		if ok {
			return vv
		}
	}
	return ""
}
