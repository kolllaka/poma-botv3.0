package music

import "fmt"

var ErrorRequestToLong = fmt.Errorf("request слишком длинный")

func getErrorRequestToLong(time, needTime int) error {
	return fmt.Errorf("%w: %dc а должен быть меньше %dc", ErrorRequestToLong, time, needTime)
}
