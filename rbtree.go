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

func (rbt *RBTree[K, V]) search(key K) (prev, target *node[K, V]) {
	curNode := rbt.root
	for curNode != rbt.leaf {
		prev = curNode
		if key < curNode.key {
			curNode = curNode.left
		} else if key > curNode.key {
			curNode = curNode.right
		} else {
			prev = curNode.parent
			target = curNode
			return
		}
	}
	return prev, nil
}

func (rbt *RBTree[K, V]) insert(key K, value V) {
	if rbt.root == nil {
		node := rbt.createNode(key, value)
		node.color = black
		rbt.root = node

	} else {
		parent, target := rbt.search(key)
		if target != nil {
			target.value = value
			return
		}
		node := rbt.createNode(key, value)
		node.parent = parent
		if node.key < parent.key {
			parent.left = node
		} else {
			parent.right = node
		}
		rbt.insertAdjust(node)
	}
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

func (rbt *RBTree[K, V]) delete(key K) bool {
	parent, target := rbt.search(key)
	if target == nil {
		return false
	}

	var isLeft bool

	// case 1: 不存在子节点，直接删除
	//               50(b)
	//         /             \
	//       20(b)          80(b)
	//    /        \        /   \
	//  13(r)     25(r)   leaf leaf
	//  /   \     /   \
	// leaf leaf leaf leaf
	// 删除 80(b)
	if target.left == rbt.leaf && target.right == rbt.leaf {
		if parent == nil {
			rbt.root = nil
		} else if parent.left == target {
			parent.left = rbt.leaf
			isLeft = true
		} else {
			isLeft = false
			parent.right = rbt.leaf
		}
	} else if target.left == rbt.leaf {
		// case 2: 删除节点只有一个子节点，替换为子节点，如果子节点仍有子节点
		// 相当于转换为case2 或者 case3继续向下进行，总能替换为case1
		// 但是根据红黑树的约束，如果删除节点只有一个子节点，那么子节点一定不存在子节点
		// left = rbt.leaf则存在右子节点
		//            25(b)
		//        /         \
		//      20(b)      50(b)
		//    /      \     /   \
		//  13(r)   leaf leaf leaf
		//  /   \
		// leaf leaf
		// delete(20)
		rbt.exchange(target, target.right)
		isLeft = false
	} else if target.right == rbt.leaf {
		rbt.exchange(target, target.left)
		isLeft = true
	} else {
		// case 3: 左右子节点都存在，查找后继(前驱)节点
		// 后继节点（前驱节点）可能 右（左）子节点
		s := rbt.precursor(target)
		//                50(b)
		//         /              \
		//       20(b)           75(b)
		//    /        \        /    \
		//  13(b)     25(b)   60(b)  80(r)
		// 			   		   \
		// 			 		   65(r)
		// delete(20)
		// s = 25(r)
		target.key = s.key
		target.value = s.value

		target = s
	}

	// 如果被删除的替换节点颜色是黑色则需要调整
	if target.color == black {
		rbt.deleteAdjust(target, isLeft)
	}
	// 删除替换节点
	if target.parent == nil {
		rbt.root = nil
	} else if target.parent.left == target {
		target.parent.left = rbt.leaf
	} else {
		target.parent.right = rbt.leaf
	}

	// 替换节点只可能存在一个子节点
	if target.left != rbt.leaf {
		rbt.exchange(target, target.left)
	} else if target.right != rbt.leaf {
		rbt.exchange(target, target.right)
	}

	target.parent = nil
	target.left = nil
	target.right = nil

	rbt.size--
	return true
}

func (rbt *RBTree[K, V]) exchange(a, b *node[K, V]) {
	if a.parent == nil {
		rbt.root = b
	} else if a == a.parent.left {
		a.parent.left = b
	} else {
		a.parent.right = b
	}
	if b != nil {
		b.parent = a.parent
	}
}

func (rbt *RBTree[K, V]) deleteAdjust(n *node[K, V], isLeft bool) {
}

func (rbt *RBTree[K, V]) precursor(n *node[K, V]) *node[K, V] {
	// 比n小的最大节点
	// 有左子树，前驱节点就是左子树的最右节点
	if n.left != rbt.leaf {
		cur := n.left
		for cur.right != rbt.leaf {
			cur = cur.right
		}
		return cur
	}
	// 没有左子树，前驱节点的右节点为该节点的父节点或祖父节点
	p := n.parent
	for p != rbt.leaf && n == p.left {
		n = p
		p = p.parent
	}
	return p
}

func (rbt *RBTree[K, V]) successor(n *node[K, V]) *node[K, V] {
	// 比n大的最小节点
	if n.right != rbt.leaf {
		cur := n.right
		for cur.left != rbt.leaf {
			cur = cur.left
		}
		return cur
	}
	p := n.parent
	for p != rbt.leaf && n == p.right {
		n = p
		p = p.parent
	}

	return p
}
