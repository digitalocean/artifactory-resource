package resource

// Put performs the Put operation for the resource
func Put(req PutRequest, dir string) (GetResponse, error) {
	var v Version

	get := GetRequest{
		Source:  req.Source,
		Version: v,
		Params:  req.Params.Get,
	}

	return Get(get, dir)
}
