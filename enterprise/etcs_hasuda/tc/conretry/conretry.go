package conretry

import (
    "time"
    "context"
)

func Retry(ctx context.Context, attempts uint, interval time.Duration, fn func() error) error {
    var err error
    for attempts > 0 {
        attempts--
        if err = fn(); err == nil {
            break
        }
        select {
        case <-ctx.Done():
            return ctx.Err()
        case <-time.After(interval):
        }
    }
    return err
}


