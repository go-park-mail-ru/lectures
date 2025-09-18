package main

import "fmt"

func main() {
	buf := []int{1, 2, 3, 4, 5}
	fmt.Println(buf)

	/*
		 buf:    [1, 2, 3, 4, 5]   ← базовый массив
		         ↑   ↑   ↑   ↑   ↑
		         │   │   │   │   │
		sl2[:2]: │───│   │   │   │ → len=2, cap=5
		sl1[1:4]:    │───┼───┤   │ → len=3, cap=4
		sl3[2:]:         │───┼───┼───┤ → len=3, cap=3

	*/

	// получение среза, указывающего на ту же память
	sl1 := buf[1:4] // [2, 3, 4]    (с 1-го до 4-го, не включая 4-й)
	sl2 := buf[:2]  // [1, 2]       (с начала до 2-го, не включая 2-й)
	sl3 := buf[2:]  // [3, 4, 5]    (со 2-го до конца)
	fmt.Println(sl1, sl2, sl3)

	fmt.Println()

	newBuf := buf[:] // [1, 2, 3, 4, 5]
	// buf = [9, 2, 3, 4, 5], т.к. та же память
	newBuf[0] = 9

	// newBuf теперь указывает на другие данные
	newBuf = append(newBuf, 6)

	// buf    = [9, 2, 3, 4, 5], не изменился
	// newBuf = [1, 2, 3, 4, 5, 6], изменился
	newBuf[0] = 1
	fmt.Println("buf", buf)
	fmt.Println("newBuf", newBuf)

	fmt.Println()

	// копирование одного слайса в другой
	var emptyBuf []int // len=0, cap=0
	// неправильно - скопирует меньшее (по len) из 2-х слайсов
	copied := copy(emptyBuf, buf) // copied = 0
	fmt.Println("Invalid copy ", copied, emptyBuf)

	// правильно
	newBuf = make([]int, len(buf), len(buf))
	copy(newBuf, buf)
	fmt.Println("Valid copy ", newBuf)

	// можно копировать в часть существующего слайса
	ints := []int{1, 2, 3, 4}
	copy(ints[1:3], []int{5, 6}) // ints = [1, 5, 6, 4]
	fmt.Println("Part copy ", ints)
}
