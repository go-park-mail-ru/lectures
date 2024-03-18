package main

func main() {
	// простое условие
	// boolVal := true
	// if boolVal {
	// 	fmt.Println("boolVal is true")
	// }

	// // условие с блоком инициализации
	// if keyValue, keyExist := mapVal["name"]; keyExist {
	// 	fmt.Println("name =", keyValue)
	// }
	// // получаем только признак сущестования ключа
	// if _, keyExist := mapVal["name"]; keyExist {
	// 	fmt.Println("key 'name' exist")
	// }

	// cond := 1
	// // множественные if else
	// if cond == 1 {
	// 	fmt.Println("cond is 1")
	// } else if cond == 2 {
	// 	fmt.Println("cond is 2")
	// }

	// // switch по 1 переменной
	// strVal := "name"
	// switch strVal {
	// case "name":
	// 	// print 123
	// 	fallthrough
	// case "test", "lastName":
	// 	// some work
	// default:
	// 	// some work
	// }

	// // switch как замена многим ifelse
	// var val1, val2 = 2, 2
	// switch {
	// case val1 > 1 || val2 < 11:
	// 	fmt.Println("first block")
	// case val2 > 10:
	// 	fmt.Println("second block")
	// }

	mapVal := map[string]string{"1": "rvasily", "2": "sulaev", "3": "pokatilov"}

	// выход из цикла, находясь внутри switch
	for key, val := range mapVal {
		println(key, val)
	} // конец for
}
