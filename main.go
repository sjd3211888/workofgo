package main

import (
	"fmt"
	_ "golearn/commontask"
	"io/ioutil"
	"math"
	"net/http"
	"reflect"
)

var w [5]int = [5]int{5, 5, 3, 4, 3}      //物品占用容量数组
var v [5]int = [5]int{40, 50, 20, 30, 30} //物品价值数组
var cap int = 10                          //背包容量10

func computer3(nIndex int, cap int) int {
	size := cap + 1             //+1是因为把第一个索引0表示为0个容量，虽然没有实际意义，但是让从索引1开始的位置代表1容量，便于理解。
	preRes := make([]int, size) //上一轮最大价值存储
	res := make([]int, size)    //这一轮最大价值存储

	//填充边界格子,把第一个物品放入能容纳下的第一行格子中
	for i := 0; i <= cap; i++ {
		if i < w[0] {
			preRes[i] = 0
		} else {
			preRes[i] = v[0]
		}
	}
	fmt.Println(1, preRes)

	//填充其他格子，外层循环是物品数量，内层循环是容量
	for i := 1; i < nIndex; i++ {
		for j := 0; j <= cap; j++ {
			vCurrent := v[i] //当前物品价值
			wCurrent := w[i] //当前物品容量
			f1 := preRes[j]  //上一个不装的最大值
			//判断是否装的下
			if j < wCurrent {
				res[j] = f1
				fmt.Println("----", j, "装不下", wCurrent, "取值=", f1)
			} else {
				capCurrent := j - wCurrent //装下后的剩余容量
				vPre := preRes[capCurrent] //获取上一轮剩余容量能存的最大价值
				f2 := vPre + vCurrent
				var biger int
				if f1 >= f2 {
					biger = f1
				} else {
					biger = f2
				}
				//biger := maxForInt(f1, f2)
				fmt.Println("----", j, ">=", wCurrent, "装的下", f1, "vs", f2, "(", vPre, "+", vCurrent, ")", "=", biger)
				res[j] = biger
			}
		}
		//用深拷贝，把res赋值给上一个数组preRes，如果用preRes=res，则是操作一个数组
		copy(preRes, res)
		fmt.Println(i+1, res)
	}

	return res[cap]
}
func Decimal(value float64) string {
	return (fmt.Sprintf("%.2f", value))
}

// 定义可排序接口
type Sortable interface {
	Len() int
	Less(int, int) bool
	Swap(int, int)
}

func BubbleSortable(arr Sortable) {
	length := arr.Len()
	for i := 0; i < length; i++ {
		for j := i; j < length-i-1; j++ {
			if arr.Less(j, j+1) {
				arr.Swap(j, j+1)
			}
		}
	}
}

type IntArr []int
type sjd []map[string]string

func (ss sjd) Len() int {
	return len(ss)
}

// 给IntArr提供Len方法
func (arr IntArr) Len() int {
	return len(arr)
}

// 给IntArr提供Less方法
func (arr IntArr) Less(i int, j int) bool {
	return arr[i] < arr[j]
}

// 给IntArr提供Swap方法
func (arr IntArr) Swap(i int, j int) {
	arr[i], arr[j] = arr[j], arr[i]
}

// ToMap 结构体转为Map[string]interface{}
func ToMap(in interface{}, tagName string) (map[string]interface{}, error) {
	out := make(map[string]interface{})

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct { // 非结构体返回错误提示
		return nil, fmt.Errorf("ToMap only accepts struct or struct pointer; got %T", v)
	}

	t := v.Type()
	// 遍历结构体字段
	// 指定tagName值为map中key;字段值为map中value
	for i := 0; i < v.NumField(); i++ {
		fi := t.Field(i)
		if tagValue := fi.Tag.Get(tagName); tagValue != "" {
			out[tagValue] = v.Field(i).Interface()
		}
	}
	return out, nil
}

type Node struct {
	value int
	left  *Node
	right *Node
}

type Tree struct {
	root   *Node
	length int
}

func creatTree() {
	arrList := []int{77, 2, 5, 7, 23, 35, 12, 17, 31}
	myTree := Tree{}
	for i := 0; i < len(arrList); i++ {
		myTree = insertNode(myTree, arrList[i])
		myTree.length++
	}
	fmt.Println(myTree)
	LDR(myTree.root)
	TreeHeight(myTree)
	fmt.Println(myTree.root)
}
func TreeHeight(tree Tree) {
	var hl = 1
	if tree.root.left != nil {
		hl = heightMax(tree.root.left, hl)
	}
	var hr = 1
	if tree.root.right != nil {
		hr = heightMax(tree.root.left, hr)
	}
	fmt.Println(hl, hr)
	fmt.Println("Tree height is ", int(math.Max(float64(hl), float64(hr))))
}

func heightMax(node *Node, h int) int {
	var hL = h
	var hR = h
	if node.left == nil && node.right == nil {
		fmt.Println(node)
		return h
	}
	if node.left != nil {
		h++
		hL = heightMax(node.left, h)
	}
	if node.right != nil {
		h++
		hR = heightMax(node.right, h)
	}
	return int(math.Max(float64(hL), float64(hR)))
}

//中序遍历
func LDR(tree *Node) []int {
	var sjd []int
	if nil == tree {
		return sjd
	} else {
		sjd1 := LDR(tree.left)
		fmt.Println(sjd1)
		sjd2 := LDR(tree.right)
		//fmt.Println(sjd2)
		sjd = append(append(sjd1, tree.value), sjd2...)
	}
	//fmt.Println(sjd)
	return sjd
}

func insertNode(tree Tree, insertValue int) Tree {
	var currentNode *Node
	var tmp *Node
	i := 0
	if tree.length == 0 {
		currentNode = new(Node)
		currentNode.value = insertValue
		tree.root = currentNode
		return tree
	} else {
		currentNode = tree.root
	}
	for {
		//fmt.Println(currentNode)
		if currentNode.value < insertValue {
			//判断是否有右节点
			if currentNode.right == nil {
				tmp = new(Node)
				tmp.value = insertValue
				currentNode.right = tmp
				break
			} else {
				currentNode = currentNode.right
				continue
			}
		} else {
			if currentNode.left == nil {
				tmp = new(Node)
				tmp.value = insertValue
				currentNode.left = tmp
				break
			} else {
				currentNode = currentNode.left
				continue
			}
		}
		i++
	}
	return tree
}

type TreeNode struct {
	Val   int
	Left  *TreeNode
	Right *TreeNode
}

//中序遍历
func LDR1(tree *Node) {
	if nil == tree {
		return
	} else {
		LDR1(tree.left)
		LDR1(tree.right)
	}
}

type graph struct {
	vertex int
	list   map[int][]int
	found  bool
}

func f1() (r int) {
	defer func() {
		r++
	}()
	return 0
}

func f2() (r int) {
	t := 5
	defer func() {
		t = t + 5
		fmt.Println("ssssssssss", t)
	}()
	return t
}

func f3() (r int) {
	defer func(r int) {
		r = r + 5
	}(r)
	return 1
}
func main() {
	request, _ := http.NewRequest("POST", "http://127.0.0.1:9002/xiaomipush", nil)
	request.Header.Set("Content-type", "application/json")
	client := http.Client{}
	response, err := client.Do(request)
	if nil != err {
		fmt.Println(err)
		return
	}
	if response.StatusCode == 200 {
		body, _ := ioutil.ReadAll(response.Body)
		fmt.Println(string(body))
	}
	select {}

}

//创建图
func NewGraph(v int) *graph {
	g := new(graph)
	g.vertex = v
	g.found = false
	g.list = map[int][]int{}
	i := 0
	for i < v {
		g.list[i] = make([]int, 0)
		i++
	}
	return g
}

//添加边
func (g *graph) addVertex(t int, s int) {
	g.list[t] = push(g.list[t], s)
	g.list[s] = push(g.list[s], t)
}

//取出切片第一个
func pop(list []int) (int, []int) {
	if len(list) > 0 {
		a := list[0]
		b := list[1:]
		return a, b
	} else {
		return -1, list
	}
}

//推入切片
func push(list []int, value int) []int {
	result := append(list, value)
	return result
}

//广度优先搜索
func (g *graph) bfs(s int, t int) {
	if s == t {
		return
	}
	visited := make([]bool, g.vertex+1)
	var queue []int
	queue = append(queue, s)
	prev := make([]int, g.vertex+1)
	i := 0
	for i < len(prev) {
		prev[i] = -1
		i++
	}
	for len(queue) != 0 {
		var w int
		w, queue = pop(queue)
		for j := 0; j < len(g.list[w]); j++ {
			q := g.list[w][j]
			fmt.Println(q)
			if !visited[q] {
				prev[q] = w
				if q == t {
					fmt.Println(prev)
					printPath(prev, s, t)
					return
				}
				visited[q] = true
				queue = append(queue, q)
			}
		}
	}
}

//深度优先搜索
func (g *graph) dfs(s int, t int) {
	prev := make([]int, g.vertex+1)
	for i := 0; i < len(prev); i++ {
		prev[i] = -1
	}
	visit := make([]bool, g.vertex+1)
	g.recurDsf(s, t, prev, visit)
	fmt.Println(prev)
	printPath(prev, s, t)
}

func (g *graph) recurDsf(w int, t int, prev []int, visited []bool) {
	if g.found {
		return
	}
	if w == t {
		g.found = true
		return
	}
	visited[w] = true
	for i := 0; i < len(g.list[w]); i++ {
		q := g.list[w][i]
		if !visited[q] {
			prev[q] = w
			g.recurDsf(g.list[w][i], t, prev, visited)
		}
	}
}

//深度优先搜做
func printPath(prev []int, s int, t int) {
	if prev[t] != -1 && s != t {
		printPath(prev, s, prev[t])
	}
	fmt.Println(t, "  ")
}
