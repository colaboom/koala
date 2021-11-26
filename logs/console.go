package logs

type ConsoleOutputer struct {
}

func (c *ConsoleOutputer) Write(data *LogData) {

}

func (c *ConsoleOutputer) Close() {

}

func NewConsoleOutputer() *ConsoleOutputer {
	return &ConsoleOutputer{}
}
