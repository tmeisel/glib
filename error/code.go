package error

type Code int

const (
	CodeUser               Code = 40000
	CodeAuthRequired       Code = 40100
	CodeForbidden          Code = 40300
	CodeNotFound           Code = 40400
	CodeConflict           Code = 40900
	CodeDuplicateKey       Code = 40901
	CodePreconditionFailed Code = 41200
	CodeGone               Code = 41000
	CodeTooManyRequests    Code = 42900
	CodeInternal           Code = 50000
)

func (c Code) String() string {
	return statusText(c)
}
