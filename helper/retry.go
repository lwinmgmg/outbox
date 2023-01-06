package helper

func Retry[T any](times int, callBack func(...T) error, arg ...T) error {
	var err error
	for i := 0; i < times; i++ {
		if err = callBack(arg...); err == nil {
			return nil
		}
	}
	return err
}
