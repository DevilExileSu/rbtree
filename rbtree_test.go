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
		return "r"
	}
	return "b"
}

/**
                 15(b)
           /                  \
        22(b)                22(b)
	  /      \             /	   \
   20(r)    20(r)       20(r)     20(r)
   /  \    	/   \
leaf leaf leaf leaf

一个子节点最长长度 - 自身长度的1/2，即为该节点的空格长度
*/
func (n *node[K, V]) String() string {
	return fmt.Sprintf("key=%v, value=%v, parent=%v, left=%v, right=%v, color=%v",
		n.key, n.value, unsafe.Pointer(n.parent), unsafe.Pointer(n.left), unsafe.Pointer(n.right), n.color)
}

type printTree struct {
	str   string
	left  *printTree
	right *printTree
	x     int
	level int
}

type printRoot struct {
	root *printTree
}

func max(a, b int) int {
	if a > b {
		return a
	}
	return b
}

func min(a, b int) int {
	if a < b {
		return a
	}
	return b
}

func String(rbt *RBTree[int, int]) string {
	levelX := make([]int, 1000)
	l := 0
	var initPrintTree func(*node[int, int], int) *printTree
	initPrintTree = func(n *node[int, int], level int) *printTree {
		if n == nil {
			l = max(level-1, l)
			return nil
		}

		left := initPrintTree(n.left, level+1)
		right := initPrintTree(n.right, level+1)

		p := new(printTree)

		if n == rbt.leaf {
			p.str = "leaf"
		} else {
			p.str = strconv.Itoa(n.key) + "(" + n.color.String() + ")"
		}
		p.left = left
		p.right = right
		p.level = level
		return p
	}

	var generateCoor func(*printTree)

	pt := initPrintTree(rbt.root, 0)
	generateCoor = func(p *printTree) {
		if p == nil {
			return
		}
		if levelX[p.level] < 0 {
			for i := 0; i < p.level; i++ {
				levelX[p.level] = -levelX[p.level]
			}
			generateCoor(pt)
			return
		}
		generateCoor(p.left)
		generateCoor(p.right)

		var px int
		if p.left != nil && p.right != nil {
			px = max((p.right.x+p.left.x)/2, 1)
		} else if p.left != nil {
			px = p.left.x + len(p.left.str)/2 + 1
		} else if p.right != nil {
			px = max(p.right.x+len(p.right.str)/2+1, levelX[p.level]+len(p.str)/2)
		} else {
			px = levelX[p.level]
		}

		if p.level > 0 && px+len(p.str)/2 <= levelX[p.level] {

			if p.left != nil {
				levelX[p.level+1] = levelX[p.level] - len(p.left.str)/2 + 1
			} else {
				levelX[p.level+1] = levelX[p.level] - len(p.str)/2 + 1
			}

			generateCoor(p)
			if p.left != nil && p.right != nil {
				px = max((p.right.x+p.left.x)/2, 1)
			} else if p.left != nil {
				px = p.left.x + len(p.left.str)/2
			} else if p.right != nil {
				px = p.right.x + len(p.right.str)/2
			} else {
				px = levelX[p.level] + 1
			}
		}
		p.x = px
		if p.right != nil {
			levelX[p.level] = max(p.right.x+len(p.right.str)+1, px+len(p.str)+1)
		} else {
			if p.left != nil {
				levelX[p.level] = px + max(len(p.left.str)/2, len(p.str)) + 1
			} else {
				levelX[p.level] = px + len(p.str) + 1
			}
		}
	}

	generateCoor(pt)

	que := make([]*printTree, 0, rbt.size)
	sb := strings.Builder{}
	que = append(que, pt)
	for len(que) != 0 {
		cnt := len(que)
		curX := 0
		curS := 0
		tmpB := strings.Builder{}
		for i := 0; i < cnt; i++ {
			n := que[0]
			que = que[1:]
			if n.x <= curX {
				n.x += (curX - n.x) + 1
			}
			sb.WriteString(strings.Repeat(" ", n.x-curX))
			sb.WriteString(n.str)
			curX = n.x + len(n.str)

			if n.left != nil {
				tmpB.WriteString(strings.Repeat(" ", n.left.x-curS+len(n.left.str)/2))
				if n.left.x != n.x {
					tmpB.WriteString("/")
				} else {
					tmpB.WriteString("|")
				}
				curS += n.left.x - curS + len(n.left.str)/2 + 1
				que = append(que, n.left)
			}
			if n.right != nil {
				tmpB.WriteString(strings.Repeat(" ", n.right.x-curS+len(n.right.str)/2-1))
				if n.right.x != n.x {
					tmpB.WriteString("\\")
				} else {
					tmpB.WriteString("|")
				}
				curS += n.right.x - curS + len(n.right.str)/2

				que = append(que, n.right)
			}
		}
		sb.WriteString("\n")
		sb.WriteString(tmpB.String() + "\n")
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
	rbt.insert(1354444, 0)
	rbt.insert(23223453334, 0)
	rbt.insert(52553453, 0)
	fmt.Println(String(rbt))

	rbt.insert(25, 0)
	fmt.Println(String(rbt))
	rbt.insert(123100823, 0)
	rbt.insert(13670065221, 0)
	rbt.insert(1436340053, 0)
	//rbt.insert(13, 0)

	fmt.Println(String(rbt))
}

func TestRBTreeDeleteCase1And2(t *testing.T) {
	rbt := NewRBTree[int, int]()
	rbt.insert(50, 0)
	rbt.insert(25, 0)
	rbt.insert(75, 0)
	rbt.insert(13, 0)
	rbt.insert(20, 0)
	fmt.Println(String(rbt))

	rbt.delete(25)
	fmt.Println(String(rbt))

	rbt.delete(20)
	fmt.Println(String(rbt))

}

func TestRBTreeDeleteCase3(t *testing.T) {
	rbt := NewRBTree[int, int]()
	rbt.insert(20, 0)
	rbt.insert(50, 0)
	rbt.insert(25, 0)
	rbt.insert(75, 0)
	rbt.insert(80, 0)

	rbt.insert(13, 0)
	rbt.insert(22, 0)
	fmt.Println(String(rbt))
	rbt.delete(20)
	fmt.Println(String(rbt))
	rbt.delete(50)
	fmt.Println(String(rbt))
}

func TestRBTreeDeleteAdjustcase1_1(t *testing.T) {
	rbt := NewRBTree[int, int]()
	rbt.insert(20, 0)
	rbt.insert(16, 0)
	rbt.insert(24, 0)
	rbt.insert(23, 0)
	rbt.insert(25, 0)
	rbt.insert(26, 0)
	fmt.Println(String(rbt))
	rbt.delete(26)
	fmt.Println(String(rbt))

	rbt.delete(24)
	fmt.Println(String(rbt))
}

func TestRBTreeDeleteAdjustcase1_2(t *testing.T) {
	rbt := NewRBTree[int, int]()
	rbt.insert(20, 0)
	rbt.insert(25, 0)
	rbt.insert(10, 0)
	rbt.insert(30, 0)
	fmt.Println(String(rbt))

	rbt.delete(10)
	fmt.Println(String(rbt))
}

func TestRBTreeDeleteAdjustcase1_3(t *testing.T) {
	rbt := NewRBTree[int, int]()
	rbt.insert(20, 0)
	rbt.insert(10, 0)

	rbt.insert(30, 0)
	rbt.insert(25, 0)
	fmt.Println(String(rbt))

	rbt.delete(20)
	fmt.Println(String(rbt))
}
