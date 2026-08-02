### 1. Debug with index

Condition:

```go
i == 5000
```

Kỳ vọng:

```text
i = 5000
user.ID = 5001
```

### 2. Debug with ID

Condition:

```go
user.ID == 7359
```

Watch:

```text
i
user
user.Balance
```

### 3. Debug with invalid item

Condition:

```go
invalid
```

use F5 to go through each invalid item.

### 4. Debug with range

Condition:

```go
i >= 4995 && i <= 5005
```

### 5. Debug each 1000 item
Condition:

```go
isCheckpoint
```

Or:

```go
i%1000 == 0
```

### 6. Debug when error is occured

Condition:

```go
err != nil
```

Watch:

```text
i
user
err
```

### 7. Debug 2 for loop

Se breakpoint at line `orderTotal := ...` or `fmt.Printf` inside order loop.

Condition:

```go
user.ID == 7359 && order.ID == 102
```

### 8. Debug before and after update

set breakpoint at center, condition:

```go
user.ID == 7359
```

Watch:

```text
oldBalance
newBalance
user.Active
```

### 9. Debug map

Set breakpoint at:

```go
mapValue := user.Balance * 2
```

Condition:

```go
id == 7359
```

Dont use hit count cuz map is unordered.

### 10. Debug by param

Condition:

```go
isTarget
```

Or: 

```go
hasInvalidEmail
hasNegativeBalance
isUnderAge
```

## Row with purpose error

| Index | ID | Error |
|---:|---:|---|
| 2499 | 2500 | Name is emtpy |
| 4999 | 5000 | Email is wrong |
| 7358 | 7359 | Balance is negative |
| 8999 | 9000 | Age is under 18 |
