# uikitdemo

Minimal Ebiten UI kit demo with consistent proportions.

## Run
```bash
go run .
```

## Design rules
- All widgets share the same control height derived from font metrics.
- External layout can only control X/Y and Width. Height is fixed by the theme.
- No "magic numbers": paddings, radius, border, etc. are derived from the control height.
