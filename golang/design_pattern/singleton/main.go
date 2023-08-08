package main

import (
	"fmt"
	"math/rand"
	"sync"
	"time"
)

// 单例对象
type Singleton struct {
	Name string
	Now  int64
}

func (s *Singleton) String() string {
	return fmt.Sprintf("%s:(地址:%p) %d", s.Name, s, s.Now)
}

const TIMES = 5

func main() {

	wg := sync.WaitGroup{}
	for i := 0; i < TIMES; i++ {
		wg.Add(1)
		go func() {
			fmt.Println(SingleSingleton()) // 可以看到，这个初始化每个对象都不一样
			wg.Done()
		}()
	}
	wg.Wait()
	for i := 0; i < TIMES; i++ {
		wg.Add(1)
		go func() {
			fmt.Println(OnceSingleton()) // 可以看到，这个初始化每个对象都一样
			wg.Done()
		}()
	}
	wg.Wait()
	for i := 0; i < TIMES; i++ {
		wg.Add(1)
		go func() {
			fmt.Println(OnceSingletonHunger()) // 可以看到，这个初始化每个对象都一样
			wg.Done()
		}()
	}
	wg.Wait()

}

var instance *Singleton

// 简单实现
func SingleSingleton() *Singleton {
	if instance == nil {
		time.Sleep(time.Microsecond * time.Duration(rand.Intn(100))) // 模拟初始化时不稳定延迟
		instance = &Singleton{
			Name: "并发不安全实例",
			Now:  time.Now().UnixMicro(),
			// now,
		} // 非并发安全，多个线程可能会多次实例化，
	}

	return instance
}

var once sync.Once
var instance2 *Singleton

// 并发安全实现
func OnceSingleton() *Singleton {
	if instance2 == nil { // 这里其实是懒汉模式，只有在首次尝试使用的时候才初始化
		// 问题：既然once.Do有一个无锁的原子操作来判断实例，那么这个地方还需要使用install是否为nil来做判断吗
		// 这个问题的核心是instance==nil判断和atomic.LoadUint32他们之间的耗时和性能如何。
		// 我看instance==nil更轻量一些，因为它没有函数栈，而且也不需要调CPU的相关指令。
		// chatgpt告诉我说比较操作需要读取变量当前值比较和判断，而CAS操作能在单个原子指令中完成，那到底谁性能更好呢，不知道
		// 写一个实例看看呗

		// 多线程执行once.Do,那么他们是不是都会阻塞在这里，而不是立即返回呢，显然它是阻塞的，这样就保证了
		time.Sleep(time.Millisecond * time.Duration(rand.Intn(100)))
		once.Do(func() {
			time.Sleep(time.Millisecond * time.Duration(rand.Intn(100))) // 模拟初始化时不稳定延迟
			instance2 = &Singleton{
				Name: "sync.Once初始化实例",
				Now:  time.Now().UnixMicro(),
			}
		})
		// once.Do看起来内部实现也是个双检锁，虽然叫双检锁，但第一层并不是锁，而是无锁的原子操作
		// 它会先通过原子操作来检查是否被初始化过，如果没有，那么进入doSlow，使用CPU原子操作的好处是不需要加锁
		// 进入doSlow后，会有一个mutex锁，如果获锁了，还需要再看一次是否被初始化过，如果未被初始化，那么才会再初始化
	}
	return instance2
}

var instance3 = &Singleton{
	Name: "饿汉实例",
	Now:  time.Now().UnixMicro(),
}

// 饿汉模式实现
func OnceSingletonHunger() *Singleton {
	if instance3 == nil {
		panic("实例未初始化") // 什么情况会出现这个呢，这里为什么要报错，比如设想一种场景，实例的初始化需要加载大内存数据，而实例需要对外提供网络服务，假设加载大内存数据需要10s，网络服务超时只有2s，那么使用懒汉模式就会一直超时
	}
	return instance3
}
