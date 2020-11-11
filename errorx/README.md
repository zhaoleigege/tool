### golang错误处理包的实现
用法
```go
import( 
    "fmt"
    "github.com/zhaoleigege/tool/errorx"
)

func main(){
    err := errorx.New(nil, "自定义错误")
    fmt.Printf("出现错误: %+v\n", err)
}

// output
//出现错误: 
           /Users/klook/Desktop/tool/errorx/e/e.go:9 main.main 自定义错误
```