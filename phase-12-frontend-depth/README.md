# Phase 12 — Frontend Depth

> Most "full-stack" engineers who lean backend can build a React form. Far fewer can debug a janky animation, optimize a 5MB bundle, audit accessibility, or explain why React re-rendered something. This phase takes you from "I can ship CRUD UIs" to "I can lead the frontend conversation in a design review." That capability matters because at Datadog you'll be writing the React + TypeScript that data engineers use every day in notebooks, and at any senior interview, "tell me about React" comes up.

**Time:** 4-6 weeks alongside everything else.

**You'll know you're done when:** you can explain React's reconciliation algorithm, optimize a slow component, audit a page for a11y, and ship a frontend with measured Core Web Vitals improvements.

---

## Table of contents

1. [Why frontend deserves serious treatment](#why)
2. [Module 12.1 — TypeScript, beyond basic types](#module-121--ts)
3. [Module 12.2 — React mental model](#module-122--react)
4. [Module 12.3 — Reconciliation and the Virtual DOM](#module-123--reconciliation)
5. [Module 12.4 — Hooks deep dive](#module-124--hooks)
6. [Module 12.5 — State management — the real options](#module-125--state)
7. [Module 12.6 — React performance: when re-renders kill you](#module-126--perf)
8. [Module 12.7 — Async, suspense, server components](#module-127--async)
9. [Module 12.8 — Accessibility (a11y), seriously](#module-128--a11y)
10. [Module 12.9 — Performance: bundle, runtime, network](#module-129--web-perf)
11. [Module 12.10 — Build tooling: Vite, esbuild, swc](#module-1210--build)
12. [Module 12.11 — Testing frontends](#module-1211--testing)
13. [Module 12.12 — Micro-frontends (briefly)](#module-1212--microfe)
14. [🛠️ Project: Notebook frontend](#project-frontend)
15. [Interview question bank](#interview-questions)
16. [What you should now know](#what-you-should-now-know)

---

<a name="why"></a>
## Why frontend deserves serious treatment

Two reasons:

1. **Datadog ADP-Notebooks ships a React + TypeScript frontend.** You'll write it. Bugs are visible. Latency is felt. Accessibility matters because some Datadog customers are required by law to have accessible tools.

2. **Most senior backend engineers are weak on frontend.** Being good at both is rare and high-leverage. When the FE team is overloaded, you can help. When a feature crosses the stack, you can build it end-to-end. This is *career-shaping*.

---

<a name="module-121--ts"></a>
## Module 12.1 — TypeScript, beyond basic types

(We covered TS basics in Phase 5. Here we go further.)

### Generics

```ts
// A function that returns whatever it gets
function identity<T>(x: T): T { return x }

// Constraint — T must extend something
function logIfHasName<T extends { name: string }>(x: T): T {
    console.log(x.name)
    return x
}

// Multiple type params
function pair<A, B>(a: A, b: B): [A, B] { return [a, b] }
const p = pair("kvng", 21)  // type: [string, number]
```

### Conditional types

```ts
type IsString<T> = T extends string ? true : false
type A = IsString<"hello">  // true
type B = IsString<42>       // false

// More useful: extract the array element type
type Element<T> = T extends (infer U)[] ? U : never
type Nums = Element<number[]>     // number
type Mixed = Element<(string | boolean)[]>  // string | boolean
```

### Mapped types

```ts
// Make every property optional
type Partial<T> = { [K in keyof T]?: T[K] }

// Make every property readonly
type Readonly<T> = { readonly [K in keyof T]: T[K] }

// Pick specific properties
type Pick<T, K extends keyof T> = { [P in K]: T[P] }

// Practical: derive a "create" type from a "stored" type
interface User { id: number; name: string; createdAt: Date }
type CreateUser = Omit<User, "id" | "createdAt">
//   = { name: string }
```

These ship in TS's standard lib. Internalize them — they're how you express "this function works with any object that has a name" or "make all fields optional."

### Discriminated unions (the most useful TS pattern)

```ts
type RemoteData<T, E> =
    | { state: "idle" }
    | { state: "loading" }
    | { state: "success"; data: T }
    | { state: "error"; error: E }

function render<T>(d: RemoteData<T, Error>) {
    switch (d.state) {
        case "idle":    return "Not started"
        case "loading": return "Loading..."
        case "success": return `Got: ${JSON.stringify(d.data)}`
        case "error":   return `Failed: ${d.error.message}`
    }
}
```

The compiler tells you if you forgot a case. Use this everywhere you have "X is in one of N states."

### `unknown` vs `any`

Both accept anything. The difference: `unknown` requires you to narrow before using, `any` doesn't.

```ts
function safe(input: unknown) {
    if (typeof input === "string") {
        // here input is narrowed to string
        console.log(input.toUpperCase())
    }
}
function unsafe(input: any) {
    console.log(input.toUpperCase())  // crashes at runtime if not string
}
```

**Default to `unknown`. Use `any` only when escape-hatching legacy code.**

---

<a name="module-122--react"></a>
## Module 12.2 — React mental model

The mental model that unlocks React: **a React component is a function from props (and state) to UI.**

```tsx
function Greeting({ name }: { name: string }) {
    return <h1>Hello, {name}</h1>
}
```

When `name` changes, React re-runs the function and updates the DOM where needed.

State adds a wrinkle: a component can hold state across re-renders.

```tsx
function Counter() {
    const [count, setCount] = useState(0)
    return <button onClick={() => setCount(c => c + 1)}>{count}</button>
}
```

The function runs every render. `useState` cheats: React remembers `count` outside the function.

### Rendering rules

A component re-renders when:
- Its own state changes (`setState`)
- Its parent re-renders (passing new props or even just identical props)
- A `useContext` value changes
- An external store it subscribes to changes (Redux, Zustand)

Re-rendering doesn't necessarily mean the DOM updates. React diffs the virtual DOM and only writes changed parts.

---

<a name="module-123--reconciliation"></a>
## Module 12.3 — Reconciliation and the Virtual DOM

> 📖 **Definition — Virtual DOM:** A lightweight in-memory tree of objects describing what UI should look like. React diffs the new vDOM against the previous, computes minimal real DOM mutations, applies them.

### How diffing works (the simplified algorithm)

For each component returning a tree:
1. Compare element type at the same position.
2. If different (e.g., was `<div>`, now `<span>`) → unmount old, mount new.
3. If same → update props, recurse into children.
4. For arrays of children, use `key` prop to match up — without keys, React just diffs by index, which produces wrong updates.

### The `key` prop matters

```tsx
// Bad — without key, when items reorder, React reuses old DOM nodes wrong
{items.map(item => <li>{item.name}</li>)}

// Good
{items.map(item => <li key={item.id}>{item.name}</li>)}
```

Don't use array index as key for reorderable lists — defeats the purpose.

### Fiber: React's incremental renderer

React 16 introduced **Fiber** — a reimplementation that splits rendering into chunks. The renderer can pause, work on something else, resume. This is how React keeps the UI responsive even during big updates.

You don't write Fiber code, but knowing it exists explains things like "why did my console.log run twice in development" (StrictMode double-invokes to catch impurities) and why hooks must be called in the same order every render (Fiber matches them by call order).

---

<a name="module-124--hooks"></a>
## Module 12.4 — Hooks deep dive

### `useState`

```tsx
const [count, setCount] = useState(0)

setCount(5)                  // direct
setCount(c => c + 1)         // functional — use when next state depends on prev
```

Functional updates are critical when batching: multiple synchronous `setCount(count + 1)` calls all see the *stale* count.

### `useEffect`

> 📖 **Definition — Effect:** Side-effects run AFTER render commits to the DOM. Network requests, subscriptions, manual DOM manipulation.

```tsx
useEffect(() => {
    const timer = setInterval(() => tick(), 1000)
    return () => clearInterval(timer)   // cleanup runs before next effect or on unmount
}, [])  // empty deps = run once on mount

useEffect(() => {
    fetch(`/users/${userId}`).then(r => r.json()).then(setUser)
}, [userId])  // re-run when userId changes
```

The dependency array is critical and notoriously misunderstood. **Every value from the component scope used inside the effect must be in the deps.** Lint rule `react-hooks/exhaustive-deps` catches violations.

### `useMemo` and `useCallback`

```tsx
// Memoize expensive computation
const sorted = useMemo(() => bigArray.sort(), [bigArray])

// Memoize a callback so children don't re-render
const handleClick = useCallback(() => {
    doThing(x)
}, [x])
```

**Don't sprinkle these everywhere.** They have overhead. Profile first; memoize only when the cost is real.

### `useRef`

Two uses:
1. Mutable value that doesn't trigger re-render: `const count = useRef(0); count.current++`
2. DOM ref: `const inputRef = useRef<HTMLInputElement>(null); <input ref={inputRef} />`

### `useContext`

```tsx
const ThemeContext = createContext<"light" | "dark">("light")

function App() {
    return <ThemeContext.Provider value="dark"><Page /></ThemeContext.Provider>
}

function Page() {
    const theme = useContext(ThemeContext)
    return <div className={theme}>...</div>
}
```

**Caveat:** every consumer re-renders when context value changes, even if they don't use the changed slice. Avoid putting frequently-changing values directly in context — use a state library (next module) instead.

### Custom hooks

The pattern: extract reusable stateful logic into a `useFoo` function.

```tsx
function useDebounced<T>(value: T, delay: number): T {
    const [debounced, setDebounced] = useState(value)
    useEffect(() => {
        const t = setTimeout(() => setDebounced(value), delay)
        return () => clearTimeout(t)
    }, [value, delay])
    return debounced
}

// In a component:
const query = useDebounced(searchInput, 300)
useEffect(() => fetch(`/search?q=${query}`), [query])
```

Custom hooks compose. They're how React codebases stay DRY.

---

<a name="module-125--state"></a>
## Module 12.5 — State management — the real options

Most apps don't need Redux. But you need *something* once you have shared state across many components.

### The progression

1. **Local state (`useState`)** — start here.
2. **Lift state up + prop drilling** — when a parent owns state used by siblings.
3. **Context** — when prop drilling gets painful and the value doesn't change too often.
4. **Server state libraries (TanStack Query, SWR)** — for data fetched from APIs.
5. **Client state libraries (Zustand, Jotai, Redux Toolkit)** — for complex shared client state.

### TanStack Query (the one to learn)

The 80/20 of frontend state management.

```tsx
const { data, isLoading, error } = useQuery({
    queryKey: ["task", taskId],
    queryFn: () => fetch(`/api/tasks/${taskId}`).then(r => r.json()),
})

if (isLoading) return <Spinner />
if (error) return <Error msg={error.message} />
return <TaskCard task={data} />
```

Handles: caching, retries, background refresh, deduplication, stale-while-revalidate. You almost never write `useEffect + fetch` again.

### Zustand (the simplest client state)

```ts
import { create } from "zustand"

const useStore = create((set) => ({
    count: 0,
    increment: () => set(state => ({ count: state.count + 1 })),
}))

// In a component
const count = useStore(state => state.count)
const increment = useStore(state => state.increment)
```

Tiny API. No reducers, no actions, no middleware unless you want them. Used by lots of modern Datadog-y companies.

### When Redux

When you have a complex, structured state machine that benefits from time-travel debugging, middleware, devtools, and your team is large enough to enforce conventions. Redux Toolkit (RTK) makes it bearable. Most apps don't need it.

---

<a name="module-126--perf"></a>
## Module 12.6 — React performance: when re-renders kill you

Re-rendering itself is fast. What makes it slow:
- Big component trees re-render unnecessarily
- Expensive computations run every render
- Many DOM updates (typing in a slow form)

### Profiling

React DevTools has a Profiler tab. Record an interaction, see which components rendered, how long each took. **Profile before optimizing.** Most "obvious" bottlenecks are wrong.

### The fixes

1. **Memoize components with `React.memo`:**
```tsx
const TaskCard = React.memo(({ task }: { task: Task }) => { ... })
```
Component skips re-render if props are shallow-equal to last time.

2. **Stable callback references with `useCallback`:**
```tsx
const onSelect = useCallback((id) => setSelected(id), [])
// Now <Child onSelect={onSelect} /> doesn't cause Child to re-render
```

3. **Stable object/array references with `useMemo`:**
```tsx
const config = useMemo(() => ({ retries: 3 }), [])
// Without useMemo, a fresh object every render → memoized children re-render
```

4. **Move state down.** State change in a parent re-renders the whole subtree. If only one child cares, push the state into that child.

5. **Virtualize long lists.** Render only visible items. `react-window` or `react-virtual`. A list of 10,000 items becomes a list of ~20 rendered + 9,980 placeholder rows.

6. **Batch updates.** React 18 auto-batches. Pre-18, multiple `setState` outside event handlers caused multiple renders.

### The classic anti-pattern

```tsx
// Bad: new array literal every render
<List items={data.filter(d => d.active)} />

// Good
const active = useMemo(() => data.filter(d => d.active), [data])
<List items={active} />
```

If `<List>` is memoized, the literal version causes it to re-render every time. The useMemo version doesn't.

---

<a name="module-127--async"></a>
## Module 12.7 — Async, suspense, server components

### Suspense

Suspense is React's mechanism for "render this while data is loading":

```tsx
<Suspense fallback={<Spinner />}>
    <UserProfile userId={42} />
</Suspense>
```

`UserProfile` can "suspend" by throwing a Promise (handled by libraries like TanStack Query, Relay, Next.js). React shows the fallback until the Promise resolves.

### React Server Components (RSC)

A 2023+ paradigm: components run on the server, ship rendered HTML + a small client runtime. Used by Next.js App Router.

Why: smaller client bundles, direct DB access from components, automatic code splitting.

Trade-offs: you write code that runs in two places (client and server) with subtle differences. Hooks like `useState` only work in client components. Maturity is improving but still a moving target.

For Datadog notebooks, you'll most likely have a SPA architecture (heavy client) — RSC is more useful for content sites.

---

<a name="module-128--a11y"></a>
## Module 12.8 — Accessibility (a11y), seriously

Most apps fail accessibility audits. At Datadog, customers in regulated industries (healthcare, government) require WCAG compliance — your code must support screen readers, keyboard nav, color-blind users, etc.

The basics:

### Semantic HTML

Use the right element for the job. `<button>` is keyboard-focusable, screen-reader-announced, has built-in click + Enter handling. A `<div onClick>` has none of that.

```tsx
// Bad
<div onClick={onClick} className="btn">Save</div>

// Good
<button onClick={onClick}>Save</button>
```

### ARIA — only when semantic HTML can't

ARIA attributes (`role`, `aria-label`, etc.) are for the cases where you genuinely can't use a semantic element. Don't sprinkle them everywhere; bad ARIA is worse than no ARIA.

### Keyboard navigation

Every interactive element must work with keyboard alone:
- Tab moves between focusable elements
- Enter/Space activate buttons
- Esc closes modals
- Arrow keys for menus, lists

Test by unplugging your mouse and trying to use your app. You'll find issues immediately.

### Focus management

When opening a modal, move focus into it. When closing, return focus to the trigger. When deleting an item, move focus to the next/previous item or a sensible default. Lost focus = lost screen reader users.

### Labels for inputs

```tsx
// Bad
<input placeholder="Email" />

// Good
<label>Email <input type="email" /></label>
// Or
<label htmlFor="email">Email</label>
<input id="email" type="email" />
```

### Color contrast

Minimum 4.5:1 for body text. Use a contrast checker. Don't rely on color alone to convey meaning (red = error needs an icon + text too, for color-blind users).

### Tools

- **axe DevTools** — Chrome extension. Runs automated checks on the current page.
- **VoiceOver** (Mac built-in screen reader) — actually use it. Cmd+F5 to toggle. Painful first time; eye-opening.
- **Storybook a11y addon** — runs axe in your component library.

---

<a name="module-129--web-perf"></a>
## Module 12.9 — Performance: bundle, runtime, network

### Core Web Vitals (Google's metrics)

- **LCP (Largest Contentful Paint):** time until the largest element renders. Target <2.5s.
- **INP (Interaction to Next Paint):** time from user interaction to UI update. Target <200ms.
- **CLS (Cumulative Layout Shift):** how much things jump around as page loads. Target <0.1.

These are what affect SEO and user perception. Measure with Chrome DevTools Lighthouse.

### Bundle size

Every byte of JS is parsed, compiled, executed. A 2MB bundle on a mid-tier phone takes seconds.

Tactics:
- **Code splitting:** load route bundles on demand. `React.lazy(() => import("./Page"))`.
- **Tree shaking:** import only what you use. `import { foo } from "lib"` not `import lib from "lib"`.
- **Analyze bundle:** `vite-plugin-visualizer`, `webpack-bundle-analyzer`. Find the bloat.
- **Replace heavy deps.** Moment.js (200KB) → date-fns (~10KB used). Lodash → individual functions or stdlib.

### Runtime perf

- **Don't update DOM in scroll handlers.** Throttle or `requestAnimationFrame`.
- **CSS animations > JS animations** when possible — they run on the compositor thread.
- **Web Workers** for CPU-heavy work (parsing big JSONs, computing diffs).

### Network

- **HTTP caching headers** — long-lived `Cache-Control: max-age=31536000` for hashed asset URLs.
- **Brotli compression** — better than gzip.
- **HTTP/2 / HTTP/3** — multiplexed; many small requests fine again.
- **CDN** — static assets at the edge.
- **Preload critical resources** — `<link rel="preload" as="font" ...>` etc.

---

<a name="module-1210--build"></a>
## Module 12.10 — Build tooling: Vite, esbuild, swc

The frontend toolchain has been through several generations:

- **2015–2019: Webpack era.** Configurable, slow, complex.
- **2020+: Vite + esbuild + swc era.** Fast (Go/Rust-written compilers), opinionated defaults.

For new projects in 2026, **default to Vite**. It uses esbuild for dev (fast) and Rollup for production builds (better tree-shaking).

```bash
npm create vite@latest myapp -- --template react-ts
cd myapp
npm install
npm run dev
```

Done. No webpack config to write.

For the very-deep-pocket world:
- **Turbopack** (Vercel's Rust-based; powers Next.js)
- **Rspack** (ByteDance's Webpack-compatible Rust replacement)
- **swc** (Rust-based TS/JS compiler used by Next.js)

You don't need to set these up — they're integrated in modern frameworks.

---

<a name="module-1211--testing"></a>
## Module 12.11 — Testing frontends

Three levels:

### Unit (Vitest, Jest)

Test individual functions, hooks. Fast. Run on every save.

```ts
import { describe, it, expect } from "vitest"
import { formatDate } from "./utils"

describe("formatDate", () => {
    it("formats a date", () => {
        expect(formatDate(new Date("2026-04-28"))).toBe("Apr 28, 2026")
    })
})
```

### Component (React Testing Library)

Test components from user perspective.

```tsx
import { render, screen } from "@testing-library/react"
import userEvent from "@testing-library/user-event"
import { TaskInput } from "./TaskInput"

it("submits new task on Enter", async () => {
    const onSubmit = vi.fn()
    render(<TaskInput onSubmit={onSubmit} />)
    const input = screen.getByLabelText("New task")
    await userEvent.type(input, "Buy milk{Enter}")
    expect(onSubmit).toHaveBeenCalledWith("Buy milk")
})
```

Key idea from RTL: **test what users see and do**, not implementation details. Avoid `enzyme`-style "find this internal state" tests — they break on refactors that don't change behavior.

### End-to-end (Playwright, Cypress)

Real browser, real backend (or mocked at network level). Slow, flaky-prone, but catches integration bugs nothing else does.

```ts
test("create a task", async ({ page }) => {
    await page.goto("/")
    await page.getByLabel("New task").fill("Buy milk")
    await page.keyboard.press("Enter")
    await expect(page.getByText("Buy milk")).toBeVisible()
})
```

Pyramid: lots of unit tests, some component tests, few e2e tests. Datadog runs CI Visibility — they care about flaky tests, so write stable ones.

---

<a name="module-1212--microfe"></a>
## Module 12.12 — Micro-frontends (briefly)

> 📖 **Definition — Micro-frontend:** Splitting one frontend into multiple independently-deployable apps that compose into a single user experience.

Why anyone wants this: a 500-engineer org can't all push to one React codebase. Teams want autonomy.

How it works (simplified):
- A "shell" app loads at the URL.
- It dynamically imports child apps as ES modules at runtime.
- Each child app owns its build, deploy, and runtime errors.
- Shared dependencies (React itself) are deduped via "module federation" (Webpack 5+) or import maps.

Frameworks: **Single-spa**, **Module Federation**, **Bit**. Or just Web Components.

For most apps under 50 engineers: **don't.** The complexity is rarely worth it. Most "we should do micro-frontends" instincts are solved by better internal modules.

---

<a name="project-frontend"></a>
## 🛠️ Project: Notebook frontend

A polished React + TypeScript + Vite frontend for the mini-ADP-Notebooks capstone.

**See `projects/notebook-frontend/` for code.**

### Spec

- **Authentication:** login form → JWT, stored in memory (not localStorage — XSS protection).
- **Notebook list page:** grid of notebooks with title, last-edited.
- **Notebook editor:** Monaco editor for cells, Run button, output rendering.
- **WebSocket connection** to backend; cell outputs stream in.
- **State:** TanStack Query for server state; Zustand for editor state.
- **Routing:** React Router.
- **A11y:** keyboard shortcuts (Cmd+Enter to run cell), focus management, ARIA labels.
- **Tests:** Vitest unit, RTL components, Playwright e2e for critical paths.
- **CI:** GitHub Actions runs tests + Lighthouse + axe.

This is the FRONT face of your capstone — same backend as the Phase 9 mini-ADP-Notebooks project.

---

<a name="interview-questions"></a>
## 🎯 Interview question bank

1. **What is the Virtual DOM and why does it exist?**

2. **Walk me through what happens when `setState` is called.**

3. **Why must hooks be called at the top level, in the same order every render?**

4. **What's the difference between `useMemo`, `useCallback`, and `React.memo`?**

5. **How would you optimize a slow list of 10,000 items?**

6. **What's the difference between client-side, server-side, and static rendering?**

7. **Explain Core Web Vitals.**

8. **How do you make a button accessible? A modal?**

9. **TanStack Query vs Redux vs Context — when each?**

10. **What's a closure stale capture, and how does it bite in `useEffect`?**

11. **Strict mode — what does it do, why does it double-render in dev?**

12. **How would you design a frontend for a real-time dashboard with hundreds of charts?**

---

<a name="what-you-should-now-know"></a>
## ✅ What you should now know

- [ ] TypeScript: generics, conditional types, mapped types, discriminated unions
- [ ] React rendering model + reconciliation
- [ ] All major hooks + custom hooks
- [ ] State management: lift up → context → server state lib → client state lib
- [ ] React performance: memo, useMemo, useCallback, virtualization
- [ ] Suspense, server components conceptually
- [ ] Accessibility fundamentals
- [ ] Core Web Vitals + how to optimize them
- [ ] Modern build tooling
- [ ] Three levels of frontend testing
- [ ] When NOT to use micro-frontends

---

🎉 **You've reached the end of the curriculum.** Time to build, ship, and interview.

Go back to the [main README](../README.md) for the bird's-eye view, or jump to the [Capstones](../phase-09-capstones/README.md) and start shipping.
