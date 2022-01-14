package sample

import "gogrpc/pb"

func NewPerson() *pb.Person {
	shirt := &pb.Shirt{
		Brand: RandomBrand("Tommy Hilfiger", "LACOSTE"),
		Name: RandomName("aaa", "bbb"),
		Color: pb.Color_BLACK,
	}

	chino := &pb.Chino{
		//Brand: RandomBrand("brand1", "brand2"),
		//Name: RandomName("name1", "name2"),
		//Color: pb.Color_BLUE,
	}

	person := &pb.Person{
		Shirt: shirt,
		Chino: chino,
	}

	return person
}