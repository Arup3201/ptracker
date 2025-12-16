# Design System — _Serious, Focused, Dense PM Tool_

**(Dark-only • WCAG AA • Hybrid UI)**

Think **GitHub Projects + Linear + Jira**, but cleaner and calmer.

## 1️ Color System (Dark-Only, Tokenized)

### Core Palette

We avoid flashy colors. Everything is low-saturation, high-contrast.

```ts
// semantic meaning > visual meaning
--bg-root:        #0d1117; // app background (GitHub dark)
--bg-surface:     #161b22; // cards, panels
--bg-elevated:    #1f2630; // modals, popovers

--border-default: #30363d;
--border-muted:   #21262d;

--text-primary:   #e6edf3;
--text-secondary: #9ba3b4;
--text-muted:     #6e7681;

--primary:        #2f81f7; // actions, links
--primary-hover:  #1f6feb;

--success:        #3fb950;
--warning:        #d29922;
--danger:         #f85149;
```

✅ **WCAG AA compliant**
✅ Neutral-first (content > chrome)
✅ Semantic tokens (no “blue-500” usage in components)

## 2️ Typography System (Serious, Compact)

### Font Choice

```txt
Primary: Inter (UI, body, labels)
Fallback: system-ui, sans-serif
```

Why:

- Excellent readability at small sizes
- Works well in dense UIs
- Neutral (doesn’t draw attention)

### Type Scale (Compact)

| Usage          | Size | Weight | Line-height |
| -------------- | ---- | ------ | ----------- |
| Page title     | 20px | 600    | 1.3         |
| Section header | 14px | 600    | 1.4         |
| Body text      | 13px | 400    | 1.5         |
| Secondary text | 12px | 400    | 1.4         |
| Labels / Meta  | 11px | 500    | 1.3         |

**Rule**:

> Never exceed 20px for headings — this keeps the app focused and dense.

## 3️ Spacing & Density Rules

### Spacing Scale (Strict)

```txt
xs  = 4px
sm  = 8px
md  = 12px
lg  = 16px
xl  = 20px
```

**Rules**

- Cards: `p-4` (16px)
- Forms: vertical gap `gap-3` (12px)
- Tables/Lists row height ≈ **36–40px**
- Never use random spacing values

## 4️ Border Radius System (Hybrid)

```txt
xs: 4px   → inputs, buttons
sm: 6px   → cards
md: 8px   → modals, popovers
```

**Rules**

- No pill buttons
- No 12px+ radius anywhere
- Sharp but not harsh

## 5️ Shadow & Elevation (Very Subtle)

Dark UIs break easily with heavy shadows — we use **micro-elevation**.

```txt
Level 0: none
Level 1: 0 1px 2px rgba(0,0,0,0.4)
Level 2: 0 4px 12px rgba(0,0,0,0.6)
```

Usage:

- Cards → Level 1
- Dropdowns / Modals → Level 2
- Never stack shadows

## 6️ Component Rules (Non-Negotiable)

### Buttons

**Primary**

- bg: `primary`
- text: `text-primary`
- height: 32px
- radius: 4px

**Secondary**

- bg: transparent
- border: `border-default`
- hover: `bg-surface`

**Danger**

- bg: transparent
- text: `danger`
- hover: `bg-danger/10`

No gradients. No glowing effects.

### Inputs

- Height: 32px
- bg: `bg-surface`
- border: `border-default`
- focus:
  - border: `primary`
  - ring: none (border-only focus)

Error:

- border: `danger`
- helper text below input (11px)

### Cards

- bg: `bg-surface`
- border: `border-default`
- radius: 6px
- shadow: level 1

Used for:

- Project tiles
- Task containers
- Panels

### Modals / Drawers

- bg: `bg-elevated`
- radius: 8px
- shadow: level 2
- width: controlled (never full-screen unless task editor)

## 7️ Layout Philosophy

### Global Layout

```
┌──────── Sidebar (icons + labels)
│
│   ┌──── Top bar (breadcrumbs, actions)
│   │
│   │   Main content (scroll only here)
```

Rules:

- Sidebar width: ~240px
- No horizontal scrolling
- Actions live in top-right, not inside content

## 8️ Tailwind Setup Direction (High Level)

You’ll define:

- CSS variables for colors
- Tailwind tokens mapping to variables
- No raw hex usage in JSX

Later I can give you:

- `tailwind.config.ts`
- `globals.css`
- Button/Input/Card components

## 9️ OAuth/Auth Page Styling Rule

Auth pages should be:

- Minimal
- Centered card
- No marketing fluff

**Keycloak redirect ≠ design heavy**

One card:

- Logo
- “Continue with …”
- Subtext (11px)

## Mental Model (Very Important)

> **If it draws attention, it’s wrong.**
> The content (tasks, status, ownership) should dominate — not the UI.
