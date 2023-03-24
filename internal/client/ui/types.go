package ui

type DataType int

const (
	TypeCard DataType = iota
	TypeCred
	TypeText
	TypeFile
)

func (t DataType) String() string {
	return [...]string{"Карта", "Логин/пароль", "Текстовые данные", "Файл"}[t]
}

type TabName int

const (
	TabCard TabName = iota
	TabCred
	TabText
	TabFile
)

func (t TabName) String() string {
	return [...]string{"Карты", "Логин/пароль", "Текстовые данные", "Файлы"}[t]
}
