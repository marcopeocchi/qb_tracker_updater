package main

func reduce[T, R any](s []T, f func(R, T) R, initValue R) R {
	acc := initValue
	for _, v := range s {
		acc = f(acc, v)
	}
	return acc
}
