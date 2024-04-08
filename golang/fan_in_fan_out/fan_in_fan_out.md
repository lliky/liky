# FAN IN FAN OUT

## FAN IN 模式
一个 goroutine 从多个通道读取数据，知道这些通道关闭，IN 是一种收敛模式，被称为扇入，可以用来收集数据。

## FAN OUT 模式
多个 goroutine c从同一个通道读取数据，知道该通道关闭。OUT 是一种张开模式，被称为扇出，可以用来分发任务。