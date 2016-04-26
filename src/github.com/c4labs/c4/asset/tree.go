package asset

// import (
// 	"bytes"
// 	"io"
// 	"strconv"
// )

// // type IDBranch [2]*ID

// type IDBranch struct {
// 	id       *ID
// 	Children []*IDBranch
// }

// type IDTree struct {
// 	idList     IDSlice
// 	branchList []*IDBranch
// 	q          []*IDBranch
// 	Root       *IDBranch
// }

// type errIDBranchOverflow int

// func (e errIDBranchOverflow) Error() string {
// 	return "IDBranch.Add called more than twice." + strconv.Itoa(int(e))
// }

// func (t *IDTree) Add(id *ID) {
// 	t.idList.Push(id)
// }

// func (t *IDTree) ID() *ID {
// 	t.idList.Sort()
// 	Root = t.buildTree(t.idList)
// }

// func (t *IDTree) buildTree() {
// 	left_cursor, right_cursor := 0, 1
// 	q_left, q_right := 0, 1
// 	list_size := len(t.idList)
// 	var br *IDBranch
// 	for left_cursor < list_size {
// 		// advance past any duplicates
// 		for right_cursor < list_size && t.idList[left_cursor] == t.idList[right_cursor] {
// 			right_cursor++
// 		}
// 		br = &IDBranch{}
// 		br.Add(t.idList[left_cursor])
// 		if right_cursor >= list_size {
// 			t.branchList = append(branchList, br)
// 			break
// 		}

// 		br.Add(t.idList[right_cursor])
// 		t.branchList = append(branchList, br)
// 		left_cursor = right_currson + 1
// 		right_currson += 2
// 		left_cursor = right_cursor - 1
// 	}
// 	left_cursor = 0
// 	right_cursor = 1
// 	list_size = len(t.branchList)
// 	for {
// 		for right_cursor < list_size && t.branchList[left_cursor] == t.branchList[right_cursor] {
// 			right_cursor++
// 		}
// 		br = &IDBranch{}
// 		br.Add(t.idList[left_cursor])
// 		if right_cursor >= list_size {
// 			t.branchList = append(branchList, br)
// 			left_cursor = list_size
// 			right_cursor = left_cursor + 1
// 			list_size = len(t.branchList)
// 			continue
// 		}
// 		br.Add(t.idList[right_cursor])
// 		t.branchList = append(branchList, br)
// 		left_cursor = right_currson + 1
// 		right_currson += 2
// 		left_cursor = right_cursor - 1
// 	}

// 	// build rest of tree from the initial set of branches

// 	// new = id(concat(idList[left_cursor], idList[right_cursor])
// 	// q = append(q, len(idList))
// 	// idList = append(idList, new)
// 	// q_right++
// 	// if(q_right-q_left == 2) {
// 	//   q_right--
// 	//   new = concat(idList[q[q_left]], idList[q[q_right]])
// 	//   idList = append(idList, new)
// 	//   q[q_left] = len(idList)
// 	//   q_left++
// 	// }

// 	// advance from right cursor past any duplicates

// 	// left_cursor = right_cursor+1
// 	// while left_cursor < list_size && idList[left_cursor] == idList[right_cursor] {
// 	//   left_cursor++
// 	// }
// 	// right_cursor = left_cursor+1
// }

// func (b *IDBranch) Add(branch IDable) error {
// 	if len(b.Children) > 1 {
// 		return errIDBranchOverflow(1)
// 	}
// 	var child *IDBranch
// 	var br interface{}
// 	br = branch
// 	switch i := br.(type) {
// 	case *ID:
// 		child = &IDBranch{
// 			id: i,
// 		}
// 	case *IDBranch:
// 		child = i
// 	}
// 	b.Children = append(b.Children, child)
// 	b.id = nil
// 	return nil
// }

// func (b *IDBranch) Clear() {
// 	b.id = nil
// 	b.Children = nil
// }

// func (b *IDBranch) ID() *ID {
// 	if b.id != nil {
// 		return b.id
// 	}
// 	encoder := NewIDEncoder()
// 	if len(b.Children) == 1 {
// 		_, err := io.Copy(encoder, bytes.NewReader(b.Children[0].ID().Bytes()))
// 		if err != nil {
// 			panic(err)
// 		}
// 		b.id = encoder.ID()
// 	} else if len(b.Children) == 2 {
// 		b.id = b.Children[0].ID().Sum(b.Children[1].ID())
// 	}
// 	return b.id
// }

// func (b *IDBranch) String() string {
// 	if len(b.Children) == 1 {
// 		return b.Children[0].ID().String()
// 	} else if len(b.Children) == 2 {
// 		var ids [2]*ID
// 		ids[0] = b.Children[0].ID()
// 		ids[1] = b.Children[1].ID()
// 		l := 0
// 		r := 1
// 		if ids[r].Cmp(ids[l]) < 0 {
// 			r = 0
// 			l = 1
// 		}
// 		return ids[l].String() + ids[r].String()
// 	}
// 	return ""
// }
