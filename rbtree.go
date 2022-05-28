package rbtree

import (
	"golang.org/x/exp/constraints"
)

type color byte

const (
	red color = iota
	black
)

type (
	node[K constraints.Ordered, V any] struct {
		key    K
		value  V
		parent *node[K, V]
		left   *node[K, V]
		right  *node[K, V]
		color  color
	}

	RBTree[K constraints.Ordered, V any] struct {
		root *node[K, V]
		size int
		leaf *node[K, V]
	}
)

func (n *node[K, V]) changeColor() {
	if n.color == black {
		n.color = red
	} else {
		n.color = black
	}
}

func (n *node[K, V]) getFather() *node[K, V] {
	return n.parent
}

func (n *node[K, V]) getGrandfather() *node[K, V] {
	p := n.getFather()
	if p == nil {
		return nil
	}
	return p.getFather()
}

func (n *node[K, V]) getUncle() *node[K, V] {
	gp := n.getGrandfather()
	if gp != nil {
		if gp.left != n.getFather() {
			return gp.left
		} else {
			return gp.right
		}
	}
	return nil
}

func NewRBTree[K constraints.Ordered, V any]() *RBTree[K, V] {
	rbt := new(RBTree[K, V])
	rbt.size = 0
	rbt.root = nil
	var key K
	var value V
	rbt.leaf = &node[K, V]{
		key,
		value,
		nil,
		nil,
		nil,
		black,
	}
	return rbt
}

func (rbt *RBTree[K, V]) createNode(key K, value V) *node[K, V] {
	return &node[K, V]{
		key,
		value,
		nil,
		rbt.leaf,
		rbt.leaf,
		red,
	}
}

func (rbt *RBTree[K, V]) search(key K) *node[K, V] {
	curNode := rbt.root
	var prevNode *node[K, V] = nil
	for curNode != rbt.leaf {
		prevNode = curNode
		if key < curNode.key {
			curNode = curNode.left
		} else {
			curNode = curNode.right
		}
	}
	return prevNode
}

func (rbt *RBTree[K, V]) insert(key K, value V) {
	node := rbt.createNode(key, value)
	if rbt.root == nil {
		node.color = black
		rbt.root = node

	} else {
		parent := rbt.search(key)
		node.parent = parent

		if node.key < parent.key {
			parent.left = node
		} else {
			parent.right = node
		}
	}
	rbt.insertAdjust(node)
	rbt.size++
}

func (rbt *RBTree[K, V]) insertAdjust(n *node[K, V]) {
	// case 1: 节点是root 或者 父节点是黑色
	if n == rbt.root || n.parent.color == black {
		return
	}
	for n != rbt.root && n.parent.color != black {
		p := n.getFather()
		if p.color == red {
			gp := n.getGrandfather()
			u := n.getUncle()
			// case 2: 父节点和叔父节点都是红色
			// 		  GP(b)				 GP(r)
			// 		 /	 \				/  \
			// 	   P(r)   U(r)	=>    P(b)  U(b)
			//     /				  /
			//    n(r)               n(r)
			// 父节点和叔父节点变色， 祖父节点变色
			// 并以祖父节点为当前节点，继续向上调整红黑树
			if u != nil && u.color == red {
				p.changeColor()
				u.changeColor()
				gp.changeColor()
				n = gp
				continue
			}

			// case 3: 父节点红色，叔父节点不存在
			// case 3.1: 父节点是祖父节点的左节点
			if p == gp.left {
				// case 3.1.1 当前节点是父节点的左节点
				// 		  GP(b)				 P(b)
				// 		 /					/  \
				// 	   P(r) 		=>    n(r)  GP(r)
				//     /
				//    n(r)
				if n == p.left {
					// 祖父节点右旋
					rbt.rightRotate(gp)
					// 父节点 和 祖父节点都变色，父节点 r->b, 祖父节点 b->r
					p.changeColor()
					gp.changeColor()

				} else {
					// case 3.1.2 父节点是左节点，当前节点是右节点
					// 		  GP(b)				 GP(b)			 n(b)
					// 		 /	 				/				/  \
					// 	   P(r) 	    =>	   n(r) 	=>    P(r)  GP(r)
					//       \				  /
					//    	 n(r)			 P(r)
					// 父节点执行一次左旋后，变为case 3.1.1的情况
					rbt.leftRotate(p)
					// 进行下一次循环处理
					n = p
					continue
				}
			} else {
				// case 3.2: 父节点是祖父节点的右节点
				// case 3.2.1: 当前节点是父节点的右节点
				//	 GP(b)                  P(b)
				// 	   \				   /  \
				// 		P(r)		=>   GP(r) n(r)
				// 	     \
				//	      n(r)
				if n == p.right {
					// 祖父节点左旋
					rbt.leftRotate(gp)
					// 父节点 和 祖父节点都变色，父节点 r->b, 祖父节点 b->r
					p.changeColor()
					gp.changeColor()
				} else {
					// case 3.2.2: 当前节点是父节点的左节点
					// 		GP(b)                  GP(b)
					// 		   \					  \
					//			P(r)		=>         n(r)
					//		   /						\
					// 		 n(r)						P(r)
					// 父节点右旋，当前节点设置为父节点，进行下一轮循环处理
					rbt.rightRotate(p)
					n = p
					continue
				}
			}
		}
	}
	rbt.root.color = black
}

func (rbt *RBTree[K, V]) rightRotate(n *node[K, V]) {
	left := n.left
	n.left = left.right

	// 更新各个节点的父节点
	if left.right != rbt.leaf {
		left.right.parent = n
	}

	left.parent = n.parent
	if n.parent == nil {
		rbt.root = left
	} else {
		if n == n.parent.right {
			n.parent.right = left
		} else {
			n.parent.left = left
		}
	}

	n.parent = left
	left.right = n
}

func (rbt *RBTree[K, V]) leftRotate(n *node[K, V]) {
	right := n.right
	n.right = right.left

	if right != rbt.leaf {
		right.left.parent = n
	}
	right.parent = n.parent
	if n.parent == nil {
		rbt.root = right
	} else {
		if n == n.parent.left {
			n.parent.left = right
		} else {
			n.parent.right = right
		}
	}

	n.parent = right
	right.left = n
}
