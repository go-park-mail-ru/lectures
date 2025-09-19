package main

import (
	"encoding/json"
	"fmt"
)

// 1. Базовая дженерик-структура с методами
type Box[T any] struct {
	Value T
}

// Метод получает значение. T уже объявлен в структуре.
func (b Box[T]) GetValue() T {
	return b.Value
}

// Метод устанавливает значение
func (b *Box[T]) SetValue(v T) {
	b.Value = v
}

// 2. Структура с ограничением типов
type Number interface {
	~int | ~float64
}

type Calculator[T Number] struct {
	Value T
}

// Методы используют то же ограничение T Number
func (c Calculator[T]) Add(a T) T {
	return c.Value + a
}

func (c Calculator[T]) Multiply(a T) T {
	return c.Value * a
}

// 3. Методы с преобразованием данных
type Container[T any] struct {
	Data T
}

func (c Container[T]) ToJSON() ([]byte, error) {
	return json.Marshal(c.Data)
}

func (c *Container[T]) FromJSON(data []byte) error {
	return json.Unmarshal(data, &c.Data)
}

// 4. Функция с дополнительным generic-параметром (методы не могут иметь свои параметры типов)
func CompareWith[T any, U any](b Box[T], other U, compareFunc func(T, U) bool) bool {
	return compareFunc(b.Value, other)
}

// 5. Практический пример: кэш с разными типами значений
type Cache[K comparable, V any] struct {
	data map[K]V
}

func NewCache[K comparable, V any]() *Cache[K, V] {
	return &Cache[K, V]{
		data: make(map[K]V),
	}
}

func (c *Cache[K, V]) Set(key K, value V) {
	c.data[key] = value
}

func (c *Cache[K, V]) Get(key K) (V, bool) {
	val, exists := c.data[key]
	return val, exists
}

func (c *Cache[K, V]) Delete(key K) {
	delete(c.data, key)
}

// 6. Пример с интерфейсным ограничением
type Stringer interface {
	String() string
}

type Printer[T Stringer] struct {
	Item T
}

func (p Printer[T]) Print() string {
	return p.Item.String()
}

// Реализуем интерфейс Stringer для кастомного типа
type MyInt int

func (m MyInt) String() string {
	return fmt.Sprintf("MyInt: %d", m)
}

func main() {
	// 1. Работа с Box
	intBox := Box[int]{Value: 42}
	strBox := Box[string]{Value: "hello"}

	fmt.Printf("intBox value: %d\n", intBox.GetValue())
	fmt.Printf("strBox value: %s\n", strBox.GetValue())

	intBox.SetValue(100)
	fmt.Printf("New intBox value: %d\n", intBox.GetValue())

	// 2. Работа с Calculator
	intCalc := Calculator[int]{Value: 10}
	floatCalc := Calculator[float64]{Value: 3.14}

	fmt.Printf("10 * 2 = %d\n", intCalc.Multiply(2))
	fmt.Printf("3.14 + 1.5 = %.2f\n", floatCalc.Add(1.5))

	// 3. Работа с Cache
	stringCache := NewCache[string, string]()
	stringCache.Set("name", "Alice")
	stringCache.Set("city", "Moscow")

	if name, exists := stringCache.Get("name"); exists {
		fmt.Printf("Name from cache: %s\n", name)
	}

	intCache := NewCache[string, int]()
	intCache.Set("age", 25)
	intCache.Set("score", 100)

	if age, exists := intCache.Get("age"); exists {
		fmt.Printf("Age from cache: %d\n", age)
	}

	// 4. Работа с методом, имеющим дополнительный generic-параметр
	box := Box[int]{Value: 42}
	result := CompareWith(box, "42", func(a int, b string) bool {
		return fmt.Sprint(a) == b
	})
	fmt.Printf("Compare 42 with '42': %t\n", result)

	// 5. Работа с интерфейсным ограничением
	myInt := MyInt(42)
	printer := Printer[MyInt]{Item: myInt}
	fmt.Println(printer.Print())
}
