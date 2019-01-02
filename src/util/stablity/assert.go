package stablity

func Assert(exp bool,desc string)  {
	if false == exp {
		panic(desc)
	}
}

func AssertNil(v interface{},desc string)  {
	Assert(nil == v,desc)
}