# learn-doc

## get duplicated index of one list
```python
>>> l = [1,2,3,1,1]
>>> [index for index, value in enumerate(l) if value == 1]
[0, 3, 4]
```
