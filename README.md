# 🍿 Binding

```go
go get github.com/ziqorg/binding
```

# 🤩 Feel The Magic

```go
x := binding.NewBinding(`{"key1": {"nested_1": [1, 2, true]}}`)
x.Get("key1.nested_1") // ([1, 2, true], nil)

x.Set("key1.nested_1", true) // dynamic type changer
x.Set("key2.keyinner1.keyinner3", "ggwp") // optional chaining upsertion
x.Set("key2.keyinner1.keyinner4", []int{1, 2})

x.Evaluate("!key1.nested_1 && key2.keyinner1.keyinner3 == 'ggwp'") // false
x.Evaluate("max(key2.keyinner1.keyinner4[0], key2.keyinner1.keyinner4[1]) == 2") // true
```
