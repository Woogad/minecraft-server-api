import (
    "context"
    "fmt"
    "time"

    "github.com/mcstatus-io/mcutil/v4/status"
)

func main() {
    ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)

    defer cancel()

    response, err := status.Modern(ctx, "demo.mcstatus.io")

    if err != nil {
        panic(err)
    }

    fmt.Println(response)
}