# MongoDB

---

## Insert:

### 1. Batch

```go
var batch []interface{}
for i := range generator.Generate(1000) {
    batch = append(batch, i)
    if len(batch) == 100 {
        collection.InsertMany(ctx, batch)
        batch = batch[:0]
    }
}
```
### 2. Concurrency
```go
jobs := make(chan interface{}, 100)
wg := sync.WaitGroup{}

for w := 0; w < 8; w++ {  
    wg.Add(1)
    go func() {
        defer wg.Done()
        for doc := range jobs {
            collection.InsertOne(ctx, doc)
        }
    }()
}

// feed workers
for doc := range generator.Generate(1000) {
    jobs <- doc
}
close(jobs)
wg.Wait()
```

