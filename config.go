package q2p

var connection_num = 32
var print = func(level int, v ...any){ return }

func SetConnectionNum(num int) {
	if num < 1 || num > 128 {
		return
	}

	connection_num = num
}

func SetPrint(f func(int, ...any)) {
	print = f
}
