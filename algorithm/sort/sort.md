## 选择排序

## 冒泡排序

## 插入排序







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