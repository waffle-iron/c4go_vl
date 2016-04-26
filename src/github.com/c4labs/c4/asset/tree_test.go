package asset_test

// import (
//   // "bytes"
//   "fmt"
//   "io"
//   "strconv"
//   "strings"
//   "testing"

//   "github.com/c4labs/c4/asset"
//   "github.com/cheekybits/is"
// )

// var _ io.Writer = (*asset.IDEncoder)(nil)
// var _ fmt.Stringer = (*asset.ID)(nil)

// func TestIDBranch(t *testing.T) {
//   is := is.New(t)
//   one := strings.NewReader(`1`)
//   two := strings.NewReader(`2`)

//   oneID := encode(one)
//   twoID := encode(two)

//   var idSlice asset.IDSlice

//   idSlice.Push(oneID)
//   idSlice.Push(twoID)
//   idSlice.Sort()

//   var idBranch asset.IDBranch

//   idBranch.Add(oneID)
//   idBranch.Add(twoID)
//   is.Equal(idSlice.ID().String(), idBranch.ID().String())
//   is.Equal(idSlice[0].String()+idSlice[1].String(), idBranch.String())
// }

// func TestIDBranchClear(t *testing.T) {
//   is := is.New(t)
//   var idBranch asset.IDBranch
//   one := strings.NewReader(`1`)
//   two := strings.NewReader(`2`)
//   oneID := encode(one)
//   twoID := encode(two)
//   idBranch.Add(oneID)
//   idBranch.Add(twoID)
//   idBranch.Clear()
//   is.Equal(idBranch.ID(), (*asset.ID)(nil))
// }

// func TestIDTree(t *testing.T) {
//   is := is.New(t)
//   var tree asset.IDTree

//   for i := 0; i < 6; i++ {
//     str := strconv.Itoa(i)
//     v := strings.NewReader(str)
//     tree.Add(encode(v))
//   }
//   is.Equal(tree.ID().String(), "c44KDwzF3aurRCWh5YM8ShV7yg4D9hRpbLJyvU7BfjwreZj1pjpVKvNFwjbi6CiLcHdGsAHuicGTiPyjyup1x678Fz")
// }

// '0': c41ZFa3qHm67fA4W1LwCrDRNyrEA5s7gK5UWukJUvSYC16oeU2JA3zhAnb95G28zTrKB9eTcJ2PEoaqg5XzZ85yS1R
// '1': c42yrSHMvUcscrQBssLhrRE28YpGUv9Gf95uH8KnwTiBv4odDbVqNnCYFs3xpsLrgVZfHebSaQQsvxgDGmw5CX1fVy
// '2': c42i2hTBA9Ej4nqEo9iUy3pJRRE53KAH9RwwMSWjmfaQN7LxCymVz1zL9hEjqeFYzxtxXz2wRK7CBtt71AFkRfHodu
// '3': c42cdkVJbn5kjthsdm8HXfkJzEVPFm8e89Zwg3R6qCDq5FowUAhJzuH2otny133sR2M4NYE8aNG44ZYLVK9TBFNTRV
// '4': c44gaNGLK7vU8JCsUgpYDsLMnrD9XK6CThiMj4SQKUNJYeWbbXMDYYksTreSZMbZvrWnqzPVKnuagVbyJsvii7boBj
// '5': c418Y8RA6GaM4oNgmvUwbj5QgFzG53ZHhCaR3KXqQkdfn6jszsQG89KfGwzcKV4gyfV2qmrB5p77hBXJidAUkqqK5i

// c44KDwzF3aurRCWh5YM8ShV7yg4D9hRpbLJyvU7BfjwreZj1pjpVKvNFwjbi6CiLcHdGsAHuicGTiPyjyup1x678Fz:
// 	c44iG1QDzNW9HeiXPvUfjGLKUbhGNxmm5NNd7qLXytYEMfFVszRwjECm3ByRuekiWQMmNZtPoU7PCZffg2bSYGBpwn:
// 		c43qq7t6eSBZvsCKss6fWfVccL8BrDz3urLTccvgE69jWq7QmXhGe4vNZhht4Va5iYg77brLGbj3ZTJGLLGeGGYRgJ:
// 			c42yrSHMvUcscrQBssLhrRE28YpGUv9Gf95uH8KnwTiBv4odDbVqNnCYFs3xpsLrgVZfHebSaQQsvxgDGmw5CX1fVy: '1'
// 			c44gaNGLK7vU8JCsUgpYDsLMnrD9XK6CThiMj4SQKUNJYeWbbXMDYYksTreSZMbZvrWnqzPVKnuagVbyJsvii7boBj: '4'
// 		c45psTWZU34q33EDVSqo5Py19Ya4wzkLWLM82oV3hsqRL7V2qcRafV453GB2yy5Y8Jf4AketsVGchYmicgQFQne5m2:
// 			c418Y8RA6GaM4oNgmvUwbj5QgFzG53ZHhCaR3KXqQkdfn6jszsQG89KfGwzcKV4gyfV2qmrB5p77hBXJidAUkqqK5i: '5'
// 			c41ZFa3qHm67fA4W1LwCrDRNyrEA5s7gK5UWukJUvSYC16oeU2JA3zhAnb95G28zTrKB9eTcJ2PEoaqg5XzZ85yS1R: '0'
// 	c45j8mMaf6T2Rexh3KYtVEHnKKkwBmJRjuDF9QXL8Suc1brPsnrvkGJAQXPKdgRq5FzjSoYJwC5vGnCMFc3VWhP4ar:
// 		c45ubKo6w9a153iCnv2BQhPepoxMvvUcJ2DnRa1KamdZVzCNZXe4LrgXBcDNEq2rzNwuwfaikXsiv1xoNvEgCJkkFn:
// 			c42cdkVJbn5kjthsdm8HXfkJzEVPFm8e89Zwg3R6qCDq5FowUAhJzuH2otny133sR2M4NYE8aNG44ZYLVK9TBFNTRV: '3'
// 			c42i2hTBA9Ej4nqEo9iUy3pJRRE53KAH9RwwMSWjmfaQN7LxCymVz1zL9hEjqeFYzxtxXz2wRK7CBtt71AFkRfHodu: '2'

// c44KDwzF3aurRCWh5YM8ShV7yg4D9hRpbLJyvU7BfjwreZj1pjpVKvNFwjbi6CiLcHdGsAHuicGTiPyjyup1x678Fz: 0
// c44iG1QDzNW9HeiXPvUfjGLKUbhGNxmm5NNd7qLXytYEMfFVszRwjECm3ByRuekiWQMmNZtPoU7PCZffg2bSYGBpwn: 1
// c43qq7t6eSBZvsCKss6fWfVccL8BrDz3urLTccvgE69jWq7QmXhGe4vNZhht4Va5iYg77brLGbj3ZTJGLLGeGGYRgJ: 2
// c42yrSHMvUcscrQBssLhrRE28YpGUv9Gf95uH8KnwTiBv4odDbVqNnCYFs3xpsLrgVZfHebSaQQsvxgDGmw5CX1fVy: 3
// c44gaNGLK7vU8JCsUgpYDsLMnrD9XK6CThiMj4SQKUNJYeWbbXMDYYksTreSZMbZvrWnqzPVKnuagVbyJsvii7boBj: 4
// c45psTWZU34q33EDVSqo5Py19Ya4wzkLWLM82oV3hsqRL7V2qcRafV453GB2yy5Y8Jf4AketsVGchYmicgQFQne5m2: 5
// c418Y8RA6GaM4oNgmvUwbj5QgFzG53ZHhCaR3KXqQkdfn6jszsQG89KfGwzcKV4gyfV2qmrB5p77hBXJidAUkqqK5i: 6
// c41ZFa3qHm67fA4W1LwCrDRNyrEA5s7gK5UWukJUvSYC16oeU2JA3zhAnb95G28zTrKB9eTcJ2PEoaqg5XzZ85yS1R: 7
// c45j8mMaf6T2Rexh3KYtVEHnKKkwBmJRjuDF9QXL8Suc1brPsnrvkGJAQXPKdgRq5FzjSoYJwC5vGnCMFc3VWhP4ar: 8
// c45ubKo6w9a153iCnv2BQhPepoxMvvUcJ2DnRa1KamdZVzCNZXe4LrgXBcDNEq2rzNwuwfaikXsiv1xoNvEgCJkkFn: 9
// c42cdkVJbn5kjthsdm8HXfkJzEVPFm8e89Zwg3R6qCDq5FowUAhJzuH2otny133sR2M4NYE8aNG44ZYLVK9TBFNTRV: 10
// c42i2hTBA9Ej4nqEo9iUy3pJRRE53KAH9RwwMSWjmfaQN7LxCymVz1zL9hEjqeFYzxtxXz2wRK7CBtt71AFkRfHodu: 11

// 3 A1(A,B)
// 	1 A(1,2)
// 	 1: c418Y8RA6GaM4oNgmvUwbj5QgFzG53ZHhCaR3KXqQkdfn6jszsQG89KfGwzcKV4gyfV2qmrB5p77hBXJidAUkqqK5i
// 	 2: c41ZFa3qHm67fA4W1LwCrDRNyrEA5s7gK5UWukJUvSYC16oeU2JA3zhAnb95G28zTrKB9eTcJ2PEoaqg5XzZ85yS1R
// 	2 B(3,4)
// 	 3: c42cdkVJbn5kjthsdm8HXfkJzEVPFm8e89Zwg3R6qCDq5FowUAhJzuH2otny133sR2M4NYE8aNG44ZYLVK9TBFNTRV
// 	 4: c42i2hTBA9Ej4nqEo9iUy3pJRRE53KAH9RwwMSWjmfaQN7LxCymVz1zL9hEjqeFYzxtxXz2wRK7CBtt71AFkRfHodu
// 5 B1(C, nil)
//   4 C(5,6)
//     5: c42yrSHMvUcscrQBssLhrRE28YpGUv9Gf95uH8KnwTiBv4odDbVqNnCYFs3xpsLrgVZfHebSaQQsvxgDGmw5CX1fVy
//     6: c44gaNGLK7vU8JCsUgpYDsLMnrD9XK6CThiMj4SQKUNJYeWbbXMDYYksTreSZMbZvrWnqzPVKnuagVbyJsvii7boBj

// 0,1,A,3,4,B,A1,5,6,C,nil,nil,B1,A2

// A2,A1,B1,A,B,C,nil,1,2,3,4,5,6,

// 1,2,3,4,5,6

// 8: 1,2,3,4,5,6,nil,nil
// 4: 1-2,3-4,5-6,nil
// 2: 1-4,5-nil
// 1: 1-nil

// 8+4+2+1

// 1-n,1-4,5-n,1-2,3-4,5-6,nil, 1,  2,  3,  4,  5,  6, nil,nil
// [ ],[ ],[ ],[ ],[ ],[ ],[ ],[ ],[ ],[ ],[ ],[ ],[ ],[ ],[ ]

// 1,2,3,4,3-4,5-6,nil,1,2,3,4,5,6,nil,nil

// 1,2,3,4,5,6,7

// 1,2,3,4,5,6,7,8

// 16: 1,2,3,4,5,6,7,8,9,nil
// 8:  1-2,3-4,5-6,7-8,9-nil,
// 4:    1-4,    5-8,     9-nil',
// 2:        1-8              9-nil''
// 1:               1-9

// 0     1    2       3     4    5      6   7   8   9   10   11 12 13 14 15 16 17 18 19
// 1-9, 1-8, 9-nil'', 1-4, 5-8, 9-nil',1-2,3-9,5-6,7-8,9-nil, 1, 2, 3, 4, 5, 6, 7, 8, 9

// (cnt / 2)+1

// 1-2

// buf: c418Y8RA6GaM4oNgmvUwbj5QgFzG53ZHhCaR3KXqQkdfn6jszsQG89KfGwzcKV4gyfV2qmrB5p77hBXJidAUkqqK5i
// buf: ID(buf+c41ZFa3qHm67fA4W1LwCrDRNyrEA5s7gK5UWukJUvSYC16oeU2JA3zhAnb95G28zTrKB9eTcJ2PEoaqg5XzZ85yS1R)
// buf2: c42cdkVJbn5kjthsdm8HXfkJzEVPFm8e89Zwg3R6qCDq5FowUAhJzuH2otny133sR2M4NYE8aNG44ZYLVK9TBFNTRV
// buf2: ID(buf2+c42i2hTBA9Ej4nqEo9iUy3pJRRE53KAH9RwwMSWjmfaQN7LxCymVz1zL9hEjqeFYzxtxXz2wRK7CBtt71AFkRfHodu)
// buf: ID(buf, buf2)
// buf2: c42yrSHMvUcscrQBssLhrRE28YpGUv9Gf95uH8KnwTiBv4odDbVqNnCYFs3xpsLrgVZfHebSaQQsvxgDGmw5CX1fVy
// buf2: c4
// c44gaNGLK7vU8JCsUgpYDsLMnrD9XK6CThiMj4SQKUNJYeWbbXMDYYksTreSZMbZvrWnqzPVKnuagVbyJsvii7boBj

// 0,1,2,3,4,5,6,7,8|  9, 10,  11,  12, 13,  14,
// A,B,C,D,E,F,G,H,I  AB, CD, ABCD, EF, GH  EFGH,

// lc> 6
// rc> 7
// q> [11, 12, 13]
// q_left> 0
// q_right> 1
// list_size> 9

// var q []int
// var list IDSlice

// list_size := len(list)
// left_cursor  = 0
// right_cursor = 1
// q_left = 2
// q_right = 2

// for {
// 	while right_cursor < list_size && list[left_cursor] == list[right_cursor] {
// 		right_cursor++
// 	}

// 	if(right_cursor >= list_size) {
// 		break
// 	}

// 	new = id(concat(list[left_cursor], list[right_cursor])
// 	q = append(q, len(list))
// 	list = append(list, new)
// 	q_right++
// 	if(q_right-q_left == 2) {
// 		q_right--
// 		new = concat(list[q[q_left]], list[q[q_right]])
// 		list = append(list, new)
// 		q[q_left] = len(list)
// 		q_left++
// 	}

// 	left_cursor = right_cursor+1

// 	while left_cursor < list_size && list[left_cursor] == list[right_cursor] {
// 		left_cursor++
// 	}

// 	right_cursor = left_cursor+1
// }
