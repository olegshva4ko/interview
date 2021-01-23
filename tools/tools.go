package tools

var (
	pattern      = []byte(`/api/block/`)
	totalPattern = []byte(`/total`)
	nums         = []byte{'0', '1', '2', '3', '4', '5', '6', '7', '8', '9'}
)

//MatchPath checks if url is matching the path /api/block/<block_number>/total
func MatchPath(url []byte) bool {
	pointer := 0
	if len(url) <= len(pattern) { //check if first part is /api/block by len
		return false
	}
	for i := range pattern {
		if url[i] == pattern [i] {
			continue
		}
		return false
	}
M:
	for i, v := range url[len(pattern):] { //check for number
		for j := range nums { //accept only number
			if v == nums[j] {
				continue M
			}
		}
		if v == '/' { // should be / before 'total' for the next processing
			pointer = len(pattern) + i
			break
		}
		return false
	}

	if len(url[pointer:]) != len(totalPattern) {//check for length, if doesnt match there can not be "total"
		return false
	}
	
	for j := range totalPattern {
		if url[pointer+j] == totalPattern[j] {
			continue
		}
		return false
	}

	return true
}
