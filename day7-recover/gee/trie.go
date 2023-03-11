package gee

import "strings"

/* jason-comment

1、实现树的数据结构，node
2、在树上，实现2个方法，
	matchChild（）：给定一个字符串，判断该字符串是否和某node的所有子节点记录，匹配，如果匹配返回该该node，否则返回nil；
  matchChildren（）：给定一个字符串，判断该字符串，匹配到的某node的，所有子节点，返回一个slice，比如/wang 就匹配到了/:name /:who /:test等多个子节点
3、树的构建和查询方法
	insert（）
	search（）

*/

// 这里其实就是Trie树中，每个节点的数据组成部分，也代表树，因为：一个树就可以由一个“根”节点标识，然后靠着不断的向下索引，就可以遍历整个树
type Trie struct {
	pattern  string  // 比如：注册的/p/:lang/doc ，插入树中后，会分别处于3层，但只有最后一层，该pattern字段，有值，上面所有层没值，可以作为判断该条分支是否到底，即该node还有没有子节点的一个标志！
	part     string  // 比如：注册的/p/:lang/doc ，那么3层，每一层分别是p :lang doc
	children []*Trie // 所有子节点的内存地址，组成的slice
	isWild   bool    // 该节点的part部分，是否通配，比如part是，:lang *filepath时，该值为true

}

func (t *Trie) matchChild(subPath string) *Trie {
	// 用于在树上插入节点时使用，查找对应层匹配到的child，如果能查找，那么递归查到的child节点，继续匹配查询后续部分，如果不能说明某一部分，在某一层就不存在匹配，就需要在这层新建一个child，挂载到children中，

	/* jason-comment
	给定一个字符串，判断该字符串是否和某node的所有子节点记录，匹配，如果匹配返回该该node，否则返回nil；
	在某个层级的、某个node上，判断给定的字符串，是否在它的children列表中，（遍历children列表，得到每个子节点的地址，然后拿到子节点的part字段，判断其是否和给定的subPath匹配，如果完全匹配，或者子节点允许通配即，isWild为true，即匹配
		匹配就返回，匹配到子节点的地址，（外层调用方，再对该子节点进一步调用matchChild方法，匹配下一个subPath）

	ps：根节点是一个虚节点，它的children字段，存储了所有一级url的内容，比如/p/:lang/doc /intro 就存储了intro 和p这2个节点的地址
	*/

	for _, child := range t.children {
		if subPath == child.part || child.isWild == true {
			return child
		}
	}

	// 遍历完成都没有，返回nil
	return nil
}

func (t *Trie) insert(pattern string, parts []string, index int) {
	/* jason-comment
	在树上插入节点
	以在一个空树上，插入/p/:lang/doc为例
		pattern是/p/:lang/doc ，最后一个节点doc，它的pattern会存储这个字符串，上2层，该字段都为空
		parts是p :lang doc 组成的slice ，分别存储在3层节点的part中，
		index是本次插入时，要插入的第几个部分，是parts的索引

	*/
	// 递归退出条件，假设插入3层，那么就需要有4层函数调用，第4层函数调用，是作用在叶子节点的的node上，此时需要填充一下第4层的pattern字段，就行，无需填充children字段，n个part，对应n+1层函数调用，n+1层是递归返回的时候

	if len(parts) == index {
		t.pattern = pattern // 填充该字段，查询时依靠节点的该字段，判断是否匹配到树的叶子节点
		return
	}
	part := parts[index] // 根据index确定，本次要插入第几个部分，index应从0开始，到？结束

	node := t.matchChild(part)
	if node == nil {
		// 说明没有匹配的子节点，和该part匹配，那么就构建一个node
		node = &Trie{part: part, children: make([]*Trie, 0)} // children引用类型，需要make初始化，

		// 根据part是否有*或者/，决定isWild是否为true
		if part[0] == ':' || part[0] == '*' {
			node.isWild = true
		}
		t.children = append(t.children, node) // 加入子节点列表
	}
	// 或者node能找到，直接用，总之：这里node一定有值,指向一个节点
	node.insert(pattern, parts, index+1) // 该node节点，要向自己的子节点列表中，插入parts中的后一个（递归）

}

func (t *Trie) matchChildren(subPath string) []*Trie {
	/* jason-comment
	在某个节点的children中搜索，找出所有 子node节点的part和subPath匹配的，并放在一个列表中返回
	*/

	nodes := make([]*Trie, 0)

	for _, child := range t.children {
		if child.part == subPath || child.isWild == true {
			nodes = append(nodes, child)
		}
	}

	return nodes // 没有一个匹配，就是nil
}

func (t *Trie) search(parts []string, index int) *Trie {
	/* jason-comment
	在给定的树上，搜索某个path的parts，是否匹配树上的一条分支，匹配的条件：path的每个part都在树中能找到对应的分支节点，且最后一部分匹配到的一定是个叶子节点，即最后一个节点的pattern要非“”
	匹配成功：
		返回最后一个叶子节点的地址，这里可以找到path所匹配到的完成的pattern；
	匹配失败：
		node节点地址为nil，
	如何匹配？
		每次从parts列表中，取index位置的part，然后利用t的matchChildren（）方法，找到所有匹配的child节点
			如果返回是nil，说明一个匹配都无，那么return nil， nil
			如果返回的slice有值，那么进行遍历
				针对每个child，都继续调用其上的search方法，path，parts原样传入，index + 1 （每个child再继续搜索，后面一个part是否匹配，递归）
		退出条件：
			当len（parts）和传入的index相等时，说明已经匹配到了最后一个节点，




	*/
	// 递归退出条件：1、匹配了parts的所有部分，2、或者该node的part是*开头，（*开头后续无需匹配）但：此时还需要判断一下，node的pattern，因为：如果还是“”，说明它下面还有子节点，不是完全匹配
	// 比如/p/:lang/doc 匹配 /p/golang/doc ，但是如果树的一条分支是/p/:lang/doc/chapter，那么/p/golang/doc就不匹配
	if len(parts) == index || strings.HasPrefix(t.part, "*") {
		if t.pattern == "" {
			return nil
		}
		return t
	}
	part := parts[index]
	nodes := t.matchChildren(part)

	for _, node := range nodes {
		result := node.search(parts, index+1)
		if result != nil {
			return result
		}
	}

	return nil

}
