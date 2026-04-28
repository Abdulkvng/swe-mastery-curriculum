# Phase 2 — OOP & Design Patterns in Go

> Object-Oriented Programming is a way of organizing code where data and the operations on that data live together. Most CS curricula teach OOP through Java's "everything is a class with inheritance trees" lens. Go takes a different lens — *composition over inheritance* — and the result is that learning OOP through Go forces you to actually understand the *principles* rather than memorize the syntax.
>
> By the end of this phase you'll know the four OOP pillars deeply, the major design patterns, why "private variables" are a thing, and how to build a real connection pool from scratch.

**Time:** 5–7 days.

**You'll know you're done when:** you can talk confidently about encapsulation, inheritance, polymorphism, and abstraction; you've implemented Singleton (thread-safely), Factory, Strategy, Observer, and Adapter; and you've built a production-grade Postgres connection pool in Go that handles concurrent requests safely.

---

## Table of contents

1. [What does this even mean? — "OOP"](#what-oop-means)
2. [Module 2.1 — Go in 60 minutes (just enough)](#module-21--go-crash-course)
3. [Module 2.2 — The four OOP pillars](#module-22--four-pillars)
4. [Module 2.3 — Composition over inheritance — Go's way](#module-23--composition)
5. [Module 2.4 — SOLID principles](#module-24--solid)
6. [Module 2.5 — Private variables & encapsulation](#module-25--private-variables)
7. [Module 2.6 — Design patterns: the ones that matter](#module-26--patterns)
8. [Module 2.7 — Singleton, properly (thread-safe in Go)](#module-27--singleton)
9. [Module 2.8 — Factory, Strategy, Observer, Adapter, Decorator](#module-28--patterns-deep)
10. [🛠️ Project: Connection Pool from scratch](#project-connection-pool)
11. [Exercises](#exercises)
12. [Interview question bank](#interview-questions)
13. [What you should now know](#what-you-should-now-know)

---

<a name="what-oop-means"></a>
## 🧠 What does this even mean? — "OOP"

Forget the buzzwords. Object-Oriented Programming (OOP) is one answer to a real engineering question: *how do we organize a codebase so it doesn't become spaghetti?*

The answer OOP gives: **bundle data and the operations on that data into the same unit (an "object"), and have these objects talk to each other through well-defined interfaces.**

Compare two ways to model a bank account:

**Procedural style:**
```go
balance := 100.0
deposit(&balance, 50.0)
withdraw(&balance, 30.0)
```

**OOP style:**
```go
account := NewAccount(100.0)
account.Deposit(50.0)
account.Withdraw(30.0)
```

Both work. The OOP version's advantage shows up at scale: when you have 50 different account types, when you need to log every withdrawal, when you need to swap a checking account for a savings account in some code path. The bundling of data + behavior + a clear interface is what makes that scale.

Now — why teach OOP through Go and not Java? Because Java's OOP comes pre-packaged with a lot of *dogma* (deep inheritance hierarchies, abstract factory factories, etc.) that experienced engineers spend years unlearning. Go gives you the *useful parts* of OOP (encapsulation, polymorphism via interfaces, composition) and forces you to skip the parts that turned out to cause more harm than good (deep inheritance).

You'll come out of this phase able to write idiomatic OOP in *any* language, because you'll understand the underlying ideas.

---

<a name="module-21--go-crash-course"></a>
## Module 2.1 — Go in 60 minutes (just enough)

If you already know Go, skim. If not, this is enough to follow the rest.

### Hello world

```go
// hello.go
package main

import "fmt"

func main() {
    fmt.Println("Hello, Kvng!")
}
```

Run: `go run hello.go`. Build: `go build hello.go && ./hello`.

### Variables, types, constants

```go
// Explicit type
var name string = "Kvng"

// Type inferred
age := 21

// Constants
const Pi = 3.14159

// Multiple assignment
x, y := 1, 2
x, y = y, x  // swap
```

### Basic types

```go
var i int = 42           // int (size depends on platform: 64-bit usually)
var i32 int32 = 42       // explicit width
var u uint = 42
var f float64 = 3.14
var b bool = true
var s string = "hello"
var bs []byte = []byte("hello")  // string ↔ []byte
```

### Functions

```go
func add(a, b int) int {
    return a + b
}

// Multiple return values — Go's idiom for errors
func divide(a, b float64) (float64, error) {
    if b == 0 {
        return 0, fmt.Errorf("divide by zero")
    }
    return a / b, nil
}

// Caller pattern
result, err := divide(10, 2)
if err != nil {
    log.Fatal(err)
}
fmt.Println(result)
```

### Slices, maps, ranges

```go
// Slice (dynamic array)
nums := []int{1, 2, 3, 4, 5}
nums = append(nums, 6)
fmt.Println(nums[0], len(nums))

// Map (hash table)
ages := map[string]int{
    "Alice": 30,
    "Bob":   25,
}
ages["Carol"] = 28
v, ok := ages["Dave"]   // ok = false if not present
if !ok { fmt.Println("not found") }

// Range loops
for i, n := range nums {
    fmt.Printf("%d: %d\n", i, n)
}
for k, v := range ages {
    fmt.Printf("%s -> %d\n", k, v)
}
```

### Structs (the heart of Go OOP)

```go
type Person struct {
    Name string
    Age  int
}

p := Person{Name: "Kvng", Age: 21}
p.Age = 22
fmt.Println(p.Name)
```

### Methods (functions attached to a type)

```go
// Receiver — like `this` or `self` in other languages.
func (p Person) Greet() string {
    return "Hi, I'm " + p.Name
}

// Pointer receiver — needed if you want to mutate the struct.
func (p *Person) Birthday() {
    p.Age++
}

p := Person{Name: "Kvng", Age: 21}
fmt.Println(p.Greet())
p.Birthday()
fmt.Println(p.Age) // 22
```

### Interfaces (Go's polymorphism)

```go
// Interface = set of method signatures.
type Greeter interface {
    Greet() string
}

// ANY type that has a Greet() string method satisfies Greeter.
// No "implements" keyword. Implicit. Powerful.
type Dog struct{ Name string }
func (d Dog) Greet() string { return "Woof! I'm " + d.Name }

func sayHi(g Greeter) {
    fmt.Println(g.Greet())
}

sayHi(Person{Name: "Kvng"})
sayHi(Dog{Name: "Rex"})
```

### Goroutines and channels (just a peek; deep in Phase 6)

```go
// A goroutine is a lightweight thread managed by Go's runtime.
go func() {
    fmt.Println("running concurrently")
}()

// A channel is a typed pipe between goroutines.
ch := make(chan int, 10)  // buffered channel, capacity 10
ch <- 42                  // send
v := <-ch                 // receive
close(ch)
```

### Error handling

```go
// Errors are values. No exceptions.
f, err := os.Open("file.txt")
if err != nil {
    return fmt.Errorf("open: %w", err)  // %w wraps, preserving the chain
}
defer f.Close()  // run when function returns
```

### Modules

```bash
go mod init github.com/kvng/myapp
go get github.com/lib/pq
go build ./...
go test ./...
```

That's enough Go to follow the rest. We'll add depth as we go.

---

<a name="module-22--four-pillars"></a>
## Module 2.2 — The four OOP pillars

Every CS textbook lists these. Most don't explain *why* they matter. Let's fix that.

### 1. Encapsulation

> 📖 **Definition — Encapsulation:** Bundling data and the methods that operate on it into a single unit, and *hiding the internals* from the outside world. Outsiders can only interact through a public interface.

Why: if internal details leak, every change ripples through the codebase. With encapsulation, you can refactor internals freely as long as the public interface stays the same.

```go
type BankAccount struct {
    balance float64  // lowercase = unexported = "private"
}

// Public methods (uppercase) form the interface.
func (a *BankAccount) Deposit(amount float64) error {
    if amount <= 0 {
        return errors.New("amount must be positive")
    }
    a.balance += amount
    return nil
}

func (a *BankAccount) Balance() float64 {
    return a.balance
}
```

Outside this package, no one can read or write `balance` directly. They have to go through `Deposit()` / `Balance()`. That gives us validation, logging, audit trails — all without callers caring.

### 2. Inheritance (and why Go skips it)

> 📖 **Definition — Inheritance:** A child class automatically gets the data and behavior of a parent class, and can override or extend it.

Java/Python:
```python
class Animal:
    def __init__(self, name): self.name = name
    def speak(self): return "..."

class Dog(Animal):           # Dog INHERITS from Animal
    def speak(self): return "Woof"
```

Sounds nice. But after 30+ years of OOP, the industry concluded: **inheritance is overused.** The problems:

- **Fragile base class:** changing a parent breaks all its children.
- **Tight coupling:** a child class is permanently glued to its parent.
- **Diamond problem:** what if a class inherits from two parents that both define `foo()`?
- **Forces taxonomy thinking** when most relationships are "has-a," not "is-a."

Go takes a clear stance: **no inheritance**. You compose types instead (next module).

### 3. Polymorphism

> 📖 **Definition — Polymorphism:** "Many forms" — same interface, different implementations. You can treat objects of different concrete types as if they were the same type, as long as they share a common interface.

This IS still in Go, but via interfaces, not class hierarchies:

```go
type Shape interface {
    Area() float64
}

type Circle struct{ R float64 }
func (c Circle) Area() float64 { return math.Pi * c.R * c.R }

type Square struct{ Side float64 }
func (s Square) Area() float64 { return s.Side * s.Side }

func TotalArea(shapes []Shape) float64 {
    total := 0.0
    for _, s := range shapes {
        total += s.Area()  // polymorphic: works on any Shape
    }
    return total
}

shapes := []Shape{Circle{R: 5}, Square{Side: 3}}
fmt.Println(TotalArea(shapes))  // 87.539...
```

`TotalArea` doesn't know if it has Circles, Squares, or new shapes invented tomorrow. As long as they implement `Area() float64`, it works.

### 4. Abstraction

> 📖 **Definition — Abstraction:** Exposing what an object *does* and hiding *how* it does it. The interface is the contract; everything behind it is an implementation detail.

You drive a car without knowing how an engine works. The steering wheel + pedals are the interface. The engine is the implementation. You can swap an electric motor for a gas engine without learning a new way to drive.

In code:

```go
type UserStore interface {
    GetByID(id int64) (User, error)
    Save(u User) error
}

// Implementation 1: Postgres
type PostgresStore struct { db *sql.DB }
func (p PostgresStore) GetByID(id int64) (User, error) { /* SQL query */ }
func (p PostgresStore) Save(u User) error              { /* INSERT/UPDATE */ }

// Implementation 2: in-memory (for tests)
type MemoryStore struct { users map[int64]User }
func (m MemoryStore) GetByID(id int64) (User, error) { /* map lookup */ }
func (m MemoryStore) Save(u User) error              { /* map write */ }

// Caller depends only on the interface.
func RegisterHandler(store UserStore) http.HandlerFunc { ... }
```

In tests, inject `MemoryStore`. In production, inject `PostgresStore`. The handler doesn't change.

### Quick review

| Pillar | One-line | Why it matters |
|---|---|---|
| Encapsulation | Bundle data + methods, hide internals | Refactor freely without breaking callers |
| Inheritance | Child gets parent's stuff (Go skips this) | Code reuse — but use composition instead |
| Polymorphism | One interface, many implementations | Generic algorithms, swap implementations |
| Abstraction | What it does, not how | Decouple modules; swap impls; testability |

---

<a name="module-23--composition"></a>
## Module 2.3 — Composition over inheritance — Go's way

Go has no `extends`. Instead, it has **embedding**: stick one struct inside another, and the outer struct gets the inner's methods.

```go
type Engine struct {
    Horsepower int
}

func (e Engine) Start() {
    fmt.Println("Vroom!")
}

// Car *embeds* Engine.
type Car struct {
    Engine        // ← embedded (no field name)
    Brand string
}

// Now any *Car has Start() automatically.
c := Car{Engine: Engine{Horsepower: 300}, Brand: "Mercedes"}
c.Start()                     // "Vroom!"
fmt.Println(c.Horsepower)     // 300 — promoted field
```

The crucial difference from inheritance: **a Car HAS-A Engine.** It doesn't IS-A Engine. You can swap the engine, have multiple, etc.

```go
type ElectricEngine struct {
    BatteryKWh float64
}
func (e ElectricEngine) Start() {
    fmt.Println("Hum...")
}

type Tesla struct {
    ElectricEngine  // embed ElectricEngine instead
    Model string
}

t := Tesla{ElectricEngine: ElectricEngine{BatteryKWh: 75}, Model: "S"}
t.Start()  // "Hum..."
```

`Car` and `Tesla` aren't related by inheritance. They both happen to have an engine that satisfies a `Starter` interface. That's it. Loose coupling, easy to test, easy to swap.

### "Composition over inheritance" in practice

The principle: **prefer giving an object a field that does the work, over deriving the object from a class that does the work.**

If you find yourself writing `class Whatever extends Something`, ask: "could I instead have a `something` field that I call?" 9 times out of 10, yes, and the resulting code is simpler.

---

<a name="module-24--solid"></a>
## Module 2.4 — SOLID principles

Five principles for maintainable OO code. Coined by Robert Martin (Uncle Bob).

### S — Single Responsibility

> A class/module should have **one reason to change.**

Bad:
```go
type User struct {
    Name string
    Email string
}

func (u *User) Save() error { /* writes to DB */ }
func (u *User) SendWelcomeEmail() error { /* sends email */ }
func (u *User) GeneratePDFReport() error { /* PDF stuff */ }
```

Three reasons to change: DB schema, email template, PDF format. Split:

```go
type User struct { Name, Email string }
type UserRepo struct { db *sql.DB }
func (r UserRepo) Save(u User) error { ... }
type Mailer struct { ... }
func (m Mailer) SendWelcome(u User) error { ... }
type ReportGen struct { ... }
func (r ReportGen) UserPDF(u User) error { ... }
```

### O — Open/Closed

> Open for extension, closed for modification. Add new behavior without editing existing code.

Bad — every new shape requires editing this function:
```go
func Area(shape interface{}) float64 {
    switch s := shape.(type) {
    case Circle: return math.Pi * s.R * s.R
    case Square: return s.Side * s.Side
    case Triangle: return 0.5 * s.Base * s.Height
    // adding Pentagon? edit here.
    }
}
```

Good — define a `Shape` interface. New shapes implement it. `Area` doesn't change:
```go
type Shape interface { Area() float64 }
func TotalArea(shapes []Shape) float64 {
    total := 0.0
    for _, s := range shapes { total += s.Area() }
    return total
}
```

### L — Liskov Substitution

> Subtypes must be substitutable for their base types without breaking correctness.

If `Bird.Fly()` is in your interface and `Penguin` is a Bird, you've violated Liskov — Penguins don't fly. The fix: `Flyer` interface separate from `Bird`. Go avoids this naturally because there's no inheritance to misuse.

### I — Interface Segregation

> Many small interfaces are better than one fat one. Don't force clients to depend on methods they don't use.

Bad:
```go
type Worker interface {
    Work()
    Eat()
    Sleep()
}
```

A `Robot` is a Worker but doesn't eat or sleep. Forcing it to implement those is wrong. Split:

```go
type Worker interface { Work() }
type Eater  interface { Eat() }
type Sleeper interface { Sleep() }
```

Go's standard library is famous for this — `io.Reader` and `io.Writer` are 1-method interfaces. They compose into `io.ReadWriter`.

### D — Dependency Inversion

> High-level modules shouldn't depend on low-level modules. Both should depend on abstractions (interfaces).

Bad:
```go
type OrderService struct {
    db *PostgresDB  // hard-wired to Postgres
}
```

Good:
```go
type OrderRepo interface { Save(o Order) error }
type OrderService struct { repo OrderRepo }  // depends on interface
```

Now `OrderService` works with any repo: Postgres, in-memory test double, mock for unit tests, MySQL on Tuesday.

---

<a name="module-25--private-variables"></a>
## Module 2.5 — Private variables & encapsulation

You asked specifically about "private variables." Let's go deep.

### Why have private variables at all?

Imagine your `BankAccount.balance` is public:

```go
account.Balance = -99999  // valid Go!
```

Now your "money" is gone, no validation triggered. With privacy:

```go
account.balance = -99999  // compile error: cannot refer to unexported field
account.Withdraw(99999)   // returns error: "insufficient funds"
```

Privacy enforces invariants. The class can promise things ("balance is always ≥ 0") because no outsider can break them.

### How different languages do privacy

| Language | Mechanism | Granularity |
|---|---|---|
| Java, C# | `private`/`protected`/`public` keywords | Per-member |
| Python | Convention: `_name` (private), `__name` (name-mangled) | Convention only — language doesn't enforce |
| TypeScript | `private`/`protected`/`public` (compile-time only); `#name` (runtime-enforced) | Per-member |
| C++ | `private:`/`protected:`/`public:` sections | Per-member |
| **Go** | **Capitalization**: `Uppercase` = exported, `lowercase` = unexported | Per-identifier, package-scoped |
| Rust | `pub` keyword; default is module-private | Per-item |

### Go's specific approach

In Go, **the package is the unit of encapsulation**, not the struct:

```go
// file: bank/account.go
package bank

type Account struct {
    balance float64  // lowercase: invisible OUTSIDE the bank package
                     // but VISIBLE inside, even from other Account methods
                     // or other types in this package
}

func (a *Account) Deposit(amount float64) {
    a.balance += amount
}
```

```go
// file: main.go
package main

import "myapp/bank"

func main() {
    a := bank.Account{}
    a.balance = 100  // ❌ COMPILE ERROR: cannot refer to unexported field
    a.Deposit(100)   // ✅ works
}
```

Code inside the `bank` package can read/write `balance` freely. Code outside cannot. This is package-level encapsulation. It encourages package boundaries to be meaningful — each package is an "object" with its own privacy.

### Constructors

Go has no `new Account()` syntax. The convention is a `NewX` function:

```go
package bank

type Account struct {
    balance float64
    id      string
}

// Constructor: validates inputs, returns a properly-built Account.
func NewAccount(initialDeposit float64) (*Account, error) {
    if initialDeposit < 0 {
        return nil, errors.New("initial deposit cannot be negative")
    }
    return &Account{
        balance: initialDeposit,
        id:      uuid.New().String(),
    }, nil
}
```

Why a `*Account` (pointer)? So the receiver methods that mutate `balance` modify the same instance.

### Getters vs. exported fields

```go
// Option A: exported field
type Config struct {
    DatabaseURL string
}

// Option B: private field + getter
type Config struct {
    databaseURL string
}
func (c Config) DatabaseURL() string { return c.databaseURL }
```

Option B is correct when you want to:
- Validate on read (e.g., return defaults)
- Compute the value lazily
- Add logging on access
- Make it read-only (no public setter)

For pure data structs that are just bags of values, option A is fine in Go. Don't getter-spam like Java. Idiomatic Go uses exported fields when there's no invariant to protect.

---

<a name="module-26--patterns"></a>
## Module 2.6 — Design patterns: the ones that matter

The "Gang of Four" book listed 23 patterns. You don't need all 23. The ones that come up constantly in real engineering and interviews:

1. **Singleton** — exactly one instance globally
2. **Factory** — encapsulate the construction of objects
3. **Strategy** — swap algorithms behind a common interface
4. **Observer** — notify many parties of an event
5. **Adapter** — translate one interface to another
6. **Decorator** — wrap an object to add behavior
7. **Builder** — construct complex objects step by step

We'll cover all of these. The next two modules dive into Singleton (because you specifically asked) and the remaining critical ones.

---

<a name="module-27--singleton"></a>
## Module 2.7 — Singleton, properly (thread-safe in Go)

### The pattern

> 📖 **Definition — Singleton:** A class for which exactly *one* instance exists globally, with a single point of access.

When you'd want it: a connection pool, a logger, a config object — things where having two would cause bugs.

### Naive (broken) version

```go
package singleton

type Logger struct { /* ... */ }

var instance *Logger

func GetLogger() *Logger {
    if instance == nil {
        instance = &Logger{}  // RACE CONDITION
    }
    return instance
}
```

Two goroutines calling `GetLogger` simultaneously can both see `instance == nil` and both create one. Now you have two "singletons." Welcome to threading bugs.

### Locked version (works, but slower than necessary)

```go
import "sync"

var (
    instance *Logger
    mu       sync.Mutex
)

func GetLogger() *Logger {
    mu.Lock()
    defer mu.Unlock()
    if instance == nil {
        instance = &Logger{}
    }
    return instance
}
```

Now correct. But every single call locks the mutex. After the first call, that's wasted time.

### Double-checked locking (clever, error-prone)

```go
func GetLogger() *Logger {
    if instance == nil {           // fast path: no lock
        mu.Lock()
        defer mu.Unlock()
        if instance == nil {       // re-check inside lock
            instance = &Logger{}
        }
    }
    return instance
}
```

Subtle: needs memory barriers in some languages (Java had a famous broken version of this for years). In Go, technically a data race on `instance` even with this pattern. Don't write this — there's a better way.

### The Go idiom: `sync.Once`

```go
import "sync"

var (
    instance *Logger
    once     sync.Once
)

func GetLogger() *Logger {
    once.Do(func() {
        instance = &Logger{}
    })
    return instance
}
```

`sync.Once.Do` runs the function exactly once across all goroutines, with proper memory ordering. Concise, correct, fast. **This is THE Go pattern.**

### Even better: avoid Singletons when you can

Singletons are a controversial pattern. Critics argue they're glorified globals — they make testing hard (you can't swap them), they couple code together, and they often signal you should use dependency injection instead.

When in doubt: **pass dependencies explicitly.** A function that takes a `*Logger` parameter is easier to test than one that calls `singleton.GetLogger()`.

But for a connection pool or app-wide config? Singleton is fine.

### Full example

See `code/singleton/` in this folder.

---

<a name="module-28--patterns-deep"></a>
## Module 2.8 — Factory, Strategy, Observer, Adapter, Decorator

### Factory

> 📖 **Definition — Factory:** A function/method that returns instances of a type, often hiding the actual concrete type behind an interface.

```go
type PaymentProcessor interface {
    Charge(amount float64) error
}

type stripeProcessor struct{ apiKey string }
func (s stripeProcessor) Charge(amount float64) error { /* call Stripe */ return nil }

type paypalProcessor struct{ clientID string }
func (p paypalProcessor) Charge(amount float64) error { /* call PayPal */ return nil }

// Factory: caller specifies what they want, gets back the interface.
func NewPaymentProcessor(kind, secret string) (PaymentProcessor, error) {
    switch kind {
    case "stripe":
        return stripeProcessor{apiKey: secret}, nil
    case "paypal":
        return paypalProcessor{clientID: secret}, nil
    default:
        return nil, fmt.Errorf("unknown processor: %s", kind)
    }
}

// Caller doesn't care which one they get.
proc, _ := NewPaymentProcessor("stripe", os.Getenv("STRIPE_KEY"))
proc.Charge(99.99)
```

Why: keeps construction logic in one place. If Stripe later requires extra setup (auth, retries), you change one function, not 50 call sites.

### Strategy

> 📖 **Definition — Strategy:** Define a family of interchangeable algorithms behind a common interface; let the caller pick at runtime.

```go
type CompressionStrategy interface {
    Compress(data []byte) []byte
}

type GzipStrategy struct{}
func (GzipStrategy) Compress(data []byte) []byte { /* gzip */ return data }

type BrotliStrategy struct{}
func (BrotliStrategy) Compress(data []byte) []byte { /* brotli */ return data }

type FileSaver struct {
    Strategy CompressionStrategy
}

func (f FileSaver) Save(data []byte, path string) error {
    compressed := f.Strategy.Compress(data)
    return os.WriteFile(path, compressed, 0644)
}

// Now I can switch strategies based on file size, file type, etc.
saver := FileSaver{Strategy: GzipStrategy{}}
saver.Save([]byte("hello"), "out.gz")
```

Strategy = "I want to swap out this one piece of behavior independently of the rest."

### Observer

> 📖 **Definition — Observer:** One object (the *subject*) maintains a list of dependents (the *observers*) and notifies them when something changes.

Pub/sub, event handlers, change notifications — all observer pattern.

```go
type Event struct {
    Name string
    Data any
}

type Observer interface {
    OnEvent(e Event)
}

type EventBus struct {
    mu        sync.RWMutex
    observers []Observer
}

func (b *EventBus) Subscribe(o Observer) {
    b.mu.Lock()
    defer b.mu.Unlock()
    b.observers = append(b.observers, o)
}

func (b *EventBus) Publish(e Event) {
    b.mu.RLock()
    snapshot := append([]Observer(nil), b.observers...)
    b.mu.RUnlock()
    for _, o := range snapshot {
        go o.OnEvent(e)  // notify in parallel
    }
}

type EmailNotifier struct{}
func (EmailNotifier) OnEvent(e Event) { fmt.Println("Email:", e.Name) }

type AnalyticsLogger struct{}
func (AnalyticsLogger) OnEvent(e Event) { fmt.Println("Logged:", e.Name) }

bus := &EventBus{}
bus.Subscribe(EmailNotifier{})
bus.Subscribe(AnalyticsLogger{})
bus.Publish(Event{Name: "user_signup"})
```

In Go, you'll often skip explicit observer types and use channels directly — same idea, language-native.

### Adapter

> 📖 **Definition — Adapter:** Translates an existing interface (often a third-party one) into the interface your code expects.

```go
// You're using interface UserStore.
type UserStore interface {
    Get(id string) (*User, error)
    Put(u *User) error
}

// But the legacy module has a different shape:
type LegacyDB struct{}
func (l LegacyDB) FindUser(uid string) (LegacyUser, error) { ... }
func (l LegacyDB) WriteRecord(uid string, data []byte) error { ... }

// Adapter glues them.
type LegacyAdapter struct {
    db LegacyDB
}

func (a LegacyAdapter) Get(id string) (*User, error) {
    legacy, err := a.db.FindUser(id)
    if err != nil { return nil, err }
    return &User{ID: legacy.UID, Name: legacy.FullName}, nil
}

func (a LegacyAdapter) Put(u *User) error {
    return a.db.WriteRecord(u.ID, serialize(u))
}

// Now LegacyAdapter satisfies UserStore. Your app code never knows.
```

Adapter pattern is your best friend when integrating with code you don't own.

### Decorator

> 📖 **Definition — Decorator:** Wrap an object so that the wrapper has the same interface but adds behavior before/after delegating.

In Go, this is *everywhere* via `http.Handler` middleware:

```go
type Handler interface {
    ServeHTTP(w ResponseWriter, r *Request)
}

func LoggingMiddleware(next Handler) Handler {
    return HandlerFunc(func(w ResponseWriter, r *Request) {
        start := time.Now()
        next.ServeHTTP(w, r)
        log.Printf("%s %s took %v", r.Method, r.URL.Path, time.Since(start))
    })
}

func AuthMiddleware(next Handler) Handler {
    return HandlerFunc(func(w ResponseWriter, r *Request) {
        if r.Header.Get("Authorization") == "" {
            http.Error(w, "unauthorized", 401)
            return
        }
        next.ServeHTTP(w, r)
    })
}

// Compose: outer wraps inner.
handler := LoggingMiddleware(AuthMiddleware(myAppHandler))
```

The outer decorator runs, then calls the inner. Logging wraps Auth wraps the actual handler. Each layer adds one concern. That's the whole magic of HTTP middleware.

---

<a name="project-connection-pool"></a>
## 🛠️ Project: Connection Pool from scratch

This is the project you specifically asked for. It's also a frequent interview topic ("how would you build a connection pool?") and the kind of thing you'll see in real codebases at Datadog.

**See `projects/connection-pool/` for full code.**

### What you'll build

A generic connection pool for `*sql.DB` (Postgres) connections, with:

- **Configurable size**: `min`, `max` connections
- **Acquire/Release semantics**: `pool.Acquire(ctx)` returns a connection; `pool.Release(conn)` returns it
- **Wait queue with timeout**: if no conn available, wait up to N seconds
- **Health checks**: ping connections periodically; replace dead ones
- **Connection lifetime**: max age and max idle time, then recycle
- **Concurrent-safe**: many goroutines calling Acquire/Release
- **Thread-safe Singleton helper**: optional `pool.Default()` using `sync.Once`
- **Graceful shutdown**: drain in-flight, close all
- **Metrics**: how many in use, idle, total acquired ever

### The interfaces (sketch)

```go
type Pool interface {
    Acquire(ctx context.Context) (*Conn, error)
    Release(c *Conn) error
    Close() error
    Stats() Stats
}

type Conn struct {
    db        *sql.DB
    createdAt time.Time
    lastUsed  time.Time
    pool      *pool
    inUse     bool
}

type Config struct {
    DSN             string
    MinConns        int
    MaxConns        int
    AcquireTimeout  time.Duration
    MaxLifetime     time.Duration
    MaxIdleTime     time.Duration
    HealthInterval  time.Duration
}
```

### Why this hits multiple OOP concepts at once

- **Encapsulation**: `Pool` exposes Acquire/Release; internal `idle` channel, `inUse` map are private.
- **Singleton**: optional `Default()` using `sync.Once`.
- **Strategy**: pluggable `HealthCheck func(*sql.Conn) error`.
- **Observer**: optional event hooks for `OnAcquire`, `OnRelease`, `OnHealthFail`.
- **SOLID**: single-responsibility (just connection management); dependency-inverted (caller passes a `Driver` interface, not hard-wired to Postgres).
- **Concurrency** (preview of Phase 6): the heart of the project is a goroutine-safe channel-based queue.

The full project is in `projects/connection-pool/`. Read the code, understand every line, modify it, write tests for it.

---

<a name="exercises"></a>
## Exercises

1. **Encapsulation drill.** Write a `Stack[T]` in Go with `Push`, `Pop`, `Peek`, `Len`, `IsEmpty`. Internal storage must be private. Add a test that confirms an outside package can't mess with the slice directly.

2. **Composition drill.** Build a `Vehicle` with a `Wheels` count, then a `Car` and `Truck` that embed it. Add `Cargo` to `Truck` only. Demonstrate a single function `Describe(v Describer)` that works on both via interface.

3. **Polymorphism drill.** Implement `Animal` interface with `Speak() string`. Make a `[]Animal` containing Dogs, Cats, Cows. Loop and print each Speak.

4. **Singleton drill.** Implement an in-memory `Cache` singleton using `sync.Once`. Bonus: also implement a non-singleton `NewCache` constructor and write a test that proves they're independent.

5. **Strategy drill.** Implement a `Sorter` interface with `Sort([]int) []int`. Provide `BubbleSort`, `QuickSort`, and `MergeSort` strategies. Benchmark them.

6. **Observer drill.** Build an `OrderEventBus` where multiple handlers can subscribe (email, slack, audit log). Publish "OrderCreated" and confirm all handlers run concurrently.

7. **Adapter drill.** You have a third-party `MetricsClient.Send(name string, value int)`. Adapt it to your codebase's `Metrics` interface that uses `Increment(name string)` and `Gauge(name string, value float64)`.

8. **Decorator drill.** Wrap an HTTP handler with three middlewares: logging, request-id injection, and a panic recoverer. Show the order of operations.

9. **Connection pool challenge.** Extend the pool with: (a) per-acquire timeouts; (b) a `Stats()` endpoint exposed over HTTP; (c) graceful drain on SIGTERM.

---

<a name="interview-questions"></a>
## 🎯 Interview question bank

1. **Explain the four pillars of OOP.**

2. **What's the difference between encapsulation and abstraction?**
   *(Encapsulation = bundling data+methods, hiding internals. Abstraction = exposing what an object does, hiding how. They overlap, but encapsulation is about ONE object; abstraction is about a contract between objects.)*

3. **What's "composition over inheritance"? Why do we say it?**

4. **Implement Singleton in your favorite language. Make it thread-safe.**

5. **What does SOLID stand for? Pick one and give a real example.**

6. **What's the difference between `private` and `protected`?**
   *(Private = same class only. Protected = same class + subclasses. Go has neither — uses package-level lowercase.)*

7. **When would you use Strategy vs Factory vs Decorator?**
   - Factory: encapsulate which concrete type to make.
   - Strategy: swap one piece of behavior at runtime.
   - Decorator: wrap an object to add behavior without touching it.

8. **Build a thread-safe LRU cache.** *(Comes up at Datadog and Apple.)*

9. **Build a connection pool.** *(Yes, this is a real interview question.)*

10. **Why are interfaces in Go better than abstract base classes in Java?**
    *(Implicit satisfaction → loose coupling. Smaller, more focused interfaces. No "implements" boilerplate. You can satisfy interfaces you don't own — useful for testing and adapting.)*

11. **Open/Closed Principle — give a code example.**

12. **What's a "code smell" that suggests you should refactor toward composition?**
    *(Deep inheritance hierarchy, child classes overriding many parent methods, parent class with optional features each child uses differently.)*

---

<a name="what-you-should-now-know"></a>
## ✅ What you should now know

- [ ] What OOP is *for*, not just what it is
- [ ] The four pillars and concrete code examples of each
- [ ] Why Go has no inheritance and what composition replaces it with
- [ ] All five SOLID principles with examples
- [ ] How privacy works in Go (package-level capitalization)
- [ ] The seven major design patterns
- [ ] Singleton with `sync.Once` (the Go idiom)
- [ ] When NOT to use Singleton (favoring DI)
- [ ] You've built a real connection pool

---

**Next:** [Phase 3 — Data Structures & Algorithms](../phase-03-dsa/README.md)
