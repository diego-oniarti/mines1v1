package shared

import (
    "fmt"
    "math/rand"
    "strings"
    "time"
)

func RandomString(length int, pre string) string {
    const dict = "qwertyuiopasdfghjklzxcvbnmQWERTYUIOPASDFGHJKLZXCVBNM-_0123456789"
    var b strings.Builder
    fmt.Fprint(&b, pre)
    rng := rand.New(rand.NewSource(time.Now().UnixNano()))
    for i := 0; i < length; i++ {
        fmt.Fprintf(&b, "%c", dict[rng.Int()%len(dict)])
    }
    return b.String()
}