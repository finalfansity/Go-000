package main

import (
    "context"
    "fmt"
    "golang.org/x/sync/errgroup"
    "net/http"
    "os"
    "os/signal"
    "syscall"
    "time"
)

type httpServe struct {
    serve http.Server
}

func NewServe(s string) *httpServe{
    return &httpServe{http.Server{Addr: s}}
}

func (h *httpServe) start() error{
    return h.serve.ListenAndServe()
}

func (h *httpServe) shutdown(ctx context.Context) error{
    return h.serve.Shutdown(ctx)
}


func main(){
    ctx, cancel := context.WithCancel(context.Background())
    lag, _ := errgroup.WithContext(ctx)
    hserve1 := NewServe(":8111")
    hserve2 := NewServe(":8112")
    lag.Go(func() error {
        if err := hserve1.start(); err != nil{
            cancel()
            return err
        }
        return nil
    })
    lag.Go(func() error {
        if err := hserve2.start(); err != nil{
            cancel()
            return err
        }
        return nil
    })


    ch := make(chan os.Signal, 1)
    signal.Notify(ch, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM,syscall.SIGINT)
    go func() {
        for {
            select {
            case sig:= <- ch:
                switch sig {
                case  syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
                    cancel()
                }
            }
        }
    }()

    go func() {
        select {
        case <- ctx.Done():
            ctx, cancel := context.WithTimeout(context.Background(), 5 * time.Second)
            defer cancel()
            go func() {
                if err := hserve1.shutdown(ctx); err != nil{
                    fmt.Println("http server 01 shutdown error ", err)
                }
            }()
            go func() {
                if err := hserve2.shutdown(ctx); err != nil{
                    fmt.Println("http server 02 shutdown error ", err)
                }
            }()
        }
    }()

    if err := lag.Wait(); err != nil{
        fmt.Println("exit finished ", err)
    }
}
