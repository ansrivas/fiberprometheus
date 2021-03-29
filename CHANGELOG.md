## ChangeLog
---
## [2021-03-29] - v2.1.2
### Bug Fix:
- Fixes #39, thanks @sunhailin-Leo

## [2021-02-08] - v2.1.1
### Enhancements:
- Fix the LICENSE headers and introduce MIT License

## [2021-01-18] - v2.1.0
### Enhancements:
- New method `NewWithLabels` now accepts a `map[string]string` so that users can create custom labels easily.
- Bumped gofiber to v2.3.3

## [2020-11-27] - v2.0.1
### Enhancements:
- Bug Fix: RequestInFlight won't decrease if ctx.Next() return error
- Bumped gofiber to v2.2.1
- Use go 1.15

## [2020-09-15] - v2.0.0
### Enhancements:
- Support gofiber-v2
- New import path would be github.com/ansrivas/fiberprometheus/v2


## [2020-07-08] - 0.3.2
### Enhancements:
- Upgrade gofiber to 1.14.4

## [2020-07-08] - 0.3.0
### Enhancements:
- Support a new method to provide a namespace and a subsystem for the service

