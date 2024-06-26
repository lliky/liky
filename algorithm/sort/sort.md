## 选择排序

## 冒泡排序

## 插入排序

## 归并排序

## 快速排序

## 堆排序

## 桶排序

*  计数排序
* 基数排序



## 几种比较排序

|                  | 时间复杂度 | 空间复杂度 | 稳定性 |
| ---------------- | ---------- | ---------- | ------ |
| 选择排序         | O(N^2)     | O(1)       | 不稳定 |
| 冒泡排序         | O(N^2)     | O(1)       | 稳定   |
| 插入排序         | O(N^2)     | O(1)       | 稳定   |
| 归并排序         | O(N*logN)  | O(N)       | 稳定   |
| 快速排序（随机） | O(N*logN)  | O(logN)    | 不稳定 |
| 堆排序           | O(N*logN)  | O(1)       | 不稳定 |

一般排序用快速排序，常量比较低

工程上，一般都是几个排序的集合，比如 快速排序 + 插入排序



## 异或运算

也是无进位相加

```go
0 ^ 0 = 0
0 ^ 1 = 1
1 ^ 0 = 1
1 ^ 1 = 0
```

### 运算定律

* 一个值和自身运算

  ```go
  x ^ x = 0
  ```

* 一个值和 0 运算

  ```go
  x ^ 0 = x
  ```

* 可交换性

  ```go
  x ^ y = y ^ x
  ```

* 结合律

  ```
  x ^ (y ^ z) = (x ^ y) ^ z
  ```