package test

import (
	"interview/tools"
	"testing"
)

func TestCheckPath(t *testing.T) {
	URLs := []struct {
		url     string
		isValid bool
	}{
		{
			"/api/block/11508993/total", true,
		},
		{
			"/apu/block/11508993/total", false,
		},
		{
			"/api/block/aaa/total", false,
		},
		{
			"/api/block/11508993total", false,
		},
		{
			"/api/block/13x23/total", false,
		},
		{
			"/api/block/12313/totaly", false,
		},
	}

	for _, url := range URLs {
		if url.isValid != tools.MatchPath([]byte(url.url)) {
			t.Errorf(url.url)
		}
	}
}
