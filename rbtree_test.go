package rbtree

import (
	"fmt"
	"strconv"
	"strings"
	"testing"
	"unsafe"
)

func (c color) String() string {
	if c == 0 {
		return "red"
	}
	return "black"
}

func (n *node[K, V]) String() string {
	return fmt.Sprintf("key=%v, value=%v, parent=%v, left=%v, right=%v, color=%v",
		n.key, n.value, unsafe.Pointer(n.parent), unsafe.Pointer(n.left), unsafe.Pointer(n.right), n.color)
}

func String(rbt *RBTree[int, int]) string {

	sb := strings.Builder{}
	root := rbt.root
	que := make([]*node[int, int], 0, rbt.size)
	que = append(que, root)
	for len(que) != 0 {
		cnt := len(que)
		for i := 0; i < cnt; i++ {
			n := que[0]
			que = que[1:]
			if n == rbt.leaf {
				sb.WriteString("leaf    ")
				continue
			}
			sb.WriteString(strconv.Itoa(n.key))
			sb.WriteString("(" + n.color.String() + ")")
			sb.WriteString(strings.Repeat(" ", 5))
			if n.left != nil {
				que = append(que, n.left)
			}
			if n.right != nil {
				que = append(que, n.right)
			}
		}
		sb.WriteString("\n")
	}
	return sb.String()
}

func TestNewRBTree(t *testing.T) {
	rbt := NewRBTree[string, int]()
	if rbt.root != nil {
		t.Fatal("error: rbtree root should nil")
	}
	if rbt.leaf.color != black {
		t.Fatal("error: rbtree leaf should black")
	}
}

func TestRBTreeInsert(t *testing.T) {
	rbt := NewRBTree[string, int]()

	rbt.insert("key05", 0)
	if rbt.size != 1 {
		t.Fatal("error: rbtree size should 1")
	}
	if rbt.root.key != "key05" {
		t.Fatal("error: rbtree root key should equal 'key'")
	}
	if rbt.root.value != 0 {
		t.Fatal("error: rbtree root value should equal 0")
	}
	if rbt.root.color != black {
		t.Fatal("error: rbtree root color should equal 0")
	}
	rbt.insert("key01", 1)
	rbt.insert("key10", 1)
	rbt.insert("key11", 2)
	if rbt.root.left.key != "key01" || rbt.root.left == rbt.leaf {
		t.Fatal("error: rbtree root.left node should equal 'key01'")
	}
	if rbt.root.right.key != "key10" || rbt.root.right == rbt.leaf {
		t.Fatal("error: rbtree root.right node should equal 'key10'")
	}
	if rbt.root.right.right.key != "key11" || rbt.root.right.right == rbt.leaf {
		t.Fatal("error: rbtree root.right.right node should equal 'key11'")
	}
	if rbt.size != 4 {
		t.Fatal("error: rbtree size should equal 4")
	}
}

func TestRBTreeInsertCase1And2(t *testing.T) {
	rbt := NewRBTree[int, int]()
	rbt.insert(10, 0)
	//  case 2: 父节点和叔父节点都是红色
	// 		  GP
	// 		 /	\
	// 	   P(r)  BP(r)
	// 父节点和叔父节点变色，祖父节点变色，如果祖父节点为root，仍是黑色
	rbt.insert(5, 0)
	rbt.insert(20, 0)

	rbt.insert(2, 0)

	fmt.Println(String(rbt))
}

func TestRBTreeInsertCase31(t *testing.T) {
	// 		10(b)
	//      /
	//     5(r)
	//    /
	//   4(r)
	// 可以通过usfca的数据结构可视化来进行验证
	// https://www.cs.usfca.edu/~galles/visualization/RedBlack.html
	rbt := NewRBTree[int, int]()
	rbt.insert(10, 0)
	rbt.insert(5, 0)
	rbt.insert(4, 0)
	// case 3.1.1
	fmt.Println(String(rbt))

	rbt.insert(2, 0)
	fmt.Println(String(rbt))
	// case 3.1.2
	rbt.insert(3, 0)
	fmt.Println(String(rbt))
}

func TestRBTreeInsertCase32(t *testing.T) {
	rbt := NewRBTree[int, int]()
	rbt.insert(10, 0)
	rbt.insert(15, 0)
	rbt.insert(20, 0)
	fmt.Println(String(rbt))

	rbt.insert(25, 0)
	fmt.Println(String(rbt))
	rbt.insert(22, 0)
	fmt.Println(String(rbt))
}
