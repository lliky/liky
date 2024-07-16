# 图







## 题目

### 图的遍历

#### 宽度优先遍历

* 利用队列实现
* 从源节点开始依次按照宽度进队列，然后弹出
* 每弹出一个点，把该节点所有没有进过队列的邻接点放入队列
* 直到队列为空



#### 深度优先遍历

* 利用栈实现
* 从源节点开始把节点按照深度放入栈，然后弹出
* 每弹出一个点，把该节点下一个没有进入过栈的邻接点放入栈
* 直到栈为空



### 拓扑排序算法

适用范围：要求有向图，且有入度为 0 的节点，且没有环

leetcode 113



### kruskal 算法

适用范围：要求无向图



### prim 算法

适用范围：要求无向图