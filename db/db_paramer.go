package db

type DBParamer struct {
	format string
	args   []interface{}
	index  int
}

func DBParamerNew(format string, args ...interface{}) *DBParamer {
	return &DBParamer{format: format, args: args, index: 1}
}

func (this *DBParamer) GetFormat() string {
	return this.format
}

func (this *DBParamer) GetArgs() []interface{} {
	return this.args
}

func (this *DBParamer) Append(format string, args ...interface{}) *DBParamer {
	this.format = this.format + " " + format
	this.args = append(this.args, args...)
	this.index = this.index + 1
	return this
}

func (this *DBParamer) AppendwithSplit(split, format string, args ...interface{}) *DBParamer {
	if this.index == 1 {
		this.format = " (" + this.format + ") "
	}
	this.format = this.format + " " + " (" + format + ") "
	this.args = append(this.args, args...)
	this.index = this.index + 1
	return this
}
