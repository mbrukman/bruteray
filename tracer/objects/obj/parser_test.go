package obj

import "fmt"

func ExampleParse() {
	obj, err := ParseFile("testdata/goph.obj")
	if err != nil {
		fmt.Println(err)
		return
	}
	fmt.Println(obj)

	//Output:
	//{[[0 0.10000000149011612 0.20000000298023224] [1 1.100000023841858 1.2000000476837158] [2 2.0999999046325684 2.200000047683716] [3 3.0999999046325684 3.200000047683716] [4 4.099999904632568 4.199999809265137] [5 5.099999904632568 5.199999809265137] [5 5.099999904632568 5.199999809265137]] map[Body:[[0 1 2] [1 2 3 4]] SkinColor:[[1 2 3 4] [2 3 4 5] [3 4 5]]]}
}