package main

import (
    "errors"
    "testing"

    "simple-docker-healthcheck/healthchecks"
)

type stubHealthCheck struct {
    ok  bool
    err error
}

func (s *stubHealthCheck) Execute() (bool, error) {
    return s.ok, s.err
}

type testExit struct {
    code int
}

func TestExecuteHealthCheck(t *testing.T) {
    originalExit := osExit
    defer func() { osExit = originalExit }()

    cases := []struct {
        name     string
        health   healthchecks.HealthCheck
        wantCode int
    }{
        {
            name:     "healthcheck passed",
            health:   &stubHealthCheck{ok: true, err: nil},
            wantCode: 0,
        },
        {
            name:     "healthcheck failed",
            health:   &stubHealthCheck{ok: false, err: nil},
            wantCode: 1,
        },
        {
            name:     "healthcheck error",
            health:   &stubHealthCheck{ok: false, err: errors.New("execute failed")},
            wantCode: 1,
        },
    }

    for _, tc := range cases {
        tc := tc
        t.Run(tc.name, func(t *testing.T) {
            osExit = func(code int) {
                panic(testExit{code: code})
            }

            defer func() {
                if r := recover(); r != nil {
                    exit, ok := r.(testExit)
                    if !ok {
                        t.Fatalf("unexpected panic: %v", r)
                    }
                    if exit.code != tc.wantCode {
                        t.Fatalf("expected exit code %d, got %d", tc.wantCode, exit.code)
                    }
                } else {
                    t.Fatal("expected osExit to be called")
                }
            }()

            executeHealthCheck(tc.health)
        })
    }
}
