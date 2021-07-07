# learn-doc

## get duplicated index of one list
```python
>>> l = [1,2,3,1,1]
>>> [index for index, value in enumerate(l) if value == 1]
[0, 3, 4]
```

## python中函数式编程支持:
### filter 函数的功能相当于过滤器。调用一个布尔函数bool_func来迭代遍历每个seq中的元素；返回一个使bool_seq返回值为true的元素的序列。
```python
a = [1,2,3,4,5,6,7]
b = filter(lambda x: x > 5, a)
print b
[6,7]
```
### map函数是对一个序列的每个项依次执行函数，下面是对一个序列每个项都乘以2：
```python
>>> a = map(lambda x:x*2,[1,2,3])
>>> list(a)
[2, 4, 6]
```
### reduce函数是对一个序列的每个项迭代调用函数，下面是求3的阶乘：
```python
>>> reduce(lambda x,y:x*y,range(1,4))
6
```

## python中字典排序
```python
dict1 = {'3': 'a', '4': 'b', '2': 'z', '5': 'x'}
sorted(dict1.items(), key=lambda x:x[0]) # sort by key
sorted(dict1.items(), key=lambda x:x[1]) # sort by value
```

## Python中单下划线和双下划线
__foo__:一种约定,Python内部的名字,用来区别其他用户自定义的命名,以防冲突，就是例如__init__(),__del__(),__call__()这些特殊方法

_foo:一种约定,用来指定变量私有.程序员用来指定私有变量的一种方式.不能用from module import * 导入，其他方面和公有一样访问；

__foo:这个有真正的意义:解析器用_classname__foo来代替这个名字,以区别和其他类相同的命名,它无法直接像公有成员一样随便访问,通过对象名._类名__xxx这样的方式可以访问.
