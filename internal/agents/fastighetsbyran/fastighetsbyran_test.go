package fastighetsbyran

import "testing"

func Test(t *testing.T) {
	f := Crawler{}
	res, err := f.GetForSale()
	if err != nil {
		t.Fatal(err)
	}
	t.Log(res)
}
