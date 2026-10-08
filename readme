# URL Health Checker — Debugging Project

## What this is

A small real-world CLI tool: it reads a list of URLs from `urls.txt`,
checks each one **concurrently** (goroutines), retries failures, and
prints a report with status code + response time.

This is the kind of tool that actually gets used (uptime monitors, CI
health checks, smoke tests). It's also full of **bugs on purpose** —
the kind beginners write when they first touch goroutines, channels,
and shared state. Your job is to find and fix them, one at a time.

## How to run it

```
cd urlchecker
go run main.go
```

⚠️ Warning: in its current state, it may **crash** or **hang forever**
(no output, no error — just stuck). That's expected. Ctrl+C to stop it.
Don't skip ahead — fix bugs in order, the later ones are hidden behind
the earlier ones.

## Project layout

```
urlchecker/
  main.go              - entry point, reads urls.txt, runs checker
  checker/checker.go   - all the concurrency logic (this is where the bugs live)
  urls.txt             - list of URLs to check
```

## Rules

- Don't rewrite the whole file. Find the *specific* broken lines and fix
  only what's needed.
- After each fix, run it again and see what changes.
- Use `go vet ./...` — it will catch at least one of these for free.
- Use `go run -race main.go` once you get past bug #2 — it will catch
  another one for you.
- Read actual Go error messages / panics carefully. They tell you the
  file and line.

---

## Task 1 (easy) — Cosmetic bug

Run the program. Look closely at the printed table header.
Something is misspelled. Find it and fix it.

**Hint:** it's in `PrintReport`.

---

## Task 2 (easy–medium) — Crash on a bad URL

Run the program (if you haven't fixed anything else yet, it will likely
crash before printing anything useful). You'll see a panic like:

```
panic: runtime error: invalid memory address or nil pointer dereference
```

Something is being used without checking whether it's safe to use.
Find where an error from a network call is being thrown away, and fix
the logic so a failed request doesn't crash the program.

**Hint:** look at `checkOne`. What does `client.Get(url)` return, and
what happens if the second value isn't `nil`?

---

## Task 3 (medium) — Every result looks the same / wrong URL

Once Task 2 is fixed, run it again. Look closely at the report: do the
URLs in the output actually match what each check was supposed to be
checking? (This bug can be subtle — on newer Go versions it may not
even show up, which is itself worth understanding.)

**Hint:** look at the `for _, url := range urls` loop in `CheckAll`
and how `url` is used *inside* the goroutine. What does each goroutine
actually capture?

---

## Task 4 (medium–hard) — The program sometimes reports fewer results than URLs

Run the program a few times in a row (fixed URLs, fixed list). Does it
always report the same number of results? A `sync.WaitGroup` is
supposed to make the main goroutine wait for *all* workers to finish —
but it's being used incorrectly here, creating a race between
"register that a goroutine exists" and "goroutine finishes."

**Hint:** `wg.Add(1)` should always happen *before* `go func(){...}()`
starts, not inside the goroutine itself.

---

## Task 5 (hard) — `go run -race` complains

Once Task 4 is fixed, run:

```
go run -race main.go
```

You should see a warning like `DATA RACE` mentioning `resultsMap`.
Multiple goroutines are writing to the same map at the same time with
no protection. There's already an unused `sync.Mutex` field sitting in
the `Checker` struct for a reason.

**Hint:** lock before writing to `c.resultsMap`, unlock after. Also
ask yourself: is `resultsMap` even necessary if `CheckAll` already
returns a slice? (Bonus: you could remove it entirely.)

---

## Task 6 (hard) — The program hangs forever, no error at all

This is the nastiest one. After fixing the above, run it again. For
some lists of URLs it may just hang with no output and no panic —
nothing in the terminal, no CPU spinning, nothing. This is a
**deadlock**, and Go can sometimes detect it for you
(`fatal error: all goroutines are asleep - deadlock!`) — but not always,
depending on timing.

The problem: the `results` channel is unbuffered, and the order of
"wait for goroutines" vs "read from the channel" is wrong — workers
are blocked trying to *send* on the channel, while the main goroutine
is blocked in `wg.Wait()` waiting for those same workers to finish
*before* it starts reading.

**Hint:** either buffer the channel to the number of URLs, or start
draining the channel in a separate goroutine *before* calling
`wg.Wait()`, then wait and close afterward.

---

## Task 7 (hardest) — Retrying a failing URL hangs forever

Point the checker at a URL that's always down (there's already one
like that in `urls.txt`: the fake domain). Even after fixing Task 6,
checking that one URL alone may never finish.

Somewhere in the retry loop, a counter that's supposed to increase on
every failed attempt... never actually increases. Find it.

**Hint:** look at the `for attempt < c.MaxRetries` loop in `checkOne`.
Where should `attempt` be incremented, and is it?

---

## After all 7 are fixed

The tool should: check every URL concurrently, retry failed ones up to
`MaxRetries` times, never hang, never crash, never race, and print a
clean report with correct per-URL status codes and timings.

## Stretch goals (make it genuinely more valuable)

Once it's bug-free, these are good next real-world features to add
yourself (no bugs provided — this is where you practice *writing*
concurrent Go instead of just fixing it):

1. Add a `-concurrency` flag to cap how many checks run at once (use a
   buffered channel or semaphore — right now it fires all requests at
   once, which real tools never do).
2. Add exponential backoff between retries instead of retrying
   instantly.
3. Write the report to a `report.json` file as well as stdout.
4. Add a unit test for `checkOne` using `httptest.Server` so you don't
   depend on the real internet.
5. Turn it into a long-running monitor: re-check all URLs every N
   seconds and only print a line when a URL's status *changes*.
