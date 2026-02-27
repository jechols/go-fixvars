package a

func f() {
	x := 1 // want "use explicit var declaration instead of :="
	_ = x

	var y = 2
	_ = y

	z, w := 3, 4 // want "use explicit var declaration instead of :="
	_, _ = z, w
}
