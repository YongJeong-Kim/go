package main

import (
	"log"
	"strings"
	"testing"
)

func TestDD(t *testing.T) {
	log.Println(fibo(1))
	ss := "aaa"
	log.Println(strings.Split(ss, ","))
}

func fibo(a int) int {
	if a <= 1 {
		return a
	}
	return fibo(a-1) + fibo(a-2)
}

// 0 1 1 2 3 5 8 13 21 34
//fibo(3) + fibo(2) -> 4
//	fibo(2) 									+ fibo(1) -> 2
//		fibo(1) + fibo(0) -> 1
//
//f(2) = f(1) + f(0) = 1
//f(3) = f(2) + f(1) = 2
//f(4) = f(3) + f(2) = 3 x -> 2
//f(5) = f(4) + f(3) = 5
//
//// product
//id	name	description	created_at		updated_at
//											default now
//
//// order
//id	product_id	ordered_at	status
//
//
//// delivery
//
//// cart
//
//// status
//id	name
//		결제완료
//		취소완료
//		//배송중
//		//배송준비중
//		//배송완료
//
//		//반품중
//		//반품완료
//		//환불중
//		//환불완료
