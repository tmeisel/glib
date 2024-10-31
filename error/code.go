package error

type Code int

const (
	CodeUser         Code = 40000
	CodeNotFound     Code = 40400
	CodeConflict     Code = 40900
	CodeDuplicateKey Code = 40901
	CodeInternal     Code = 50000
)

func (c Code) String() string {
	return statusText(c)
}
