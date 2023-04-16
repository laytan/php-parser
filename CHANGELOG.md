# Changelog

All notable changes to this project will be documented in this file, in reverse chronological order by release.

## v0.10.0 2023-17-04

### Breaking Changes

- Changed VKCOM to laytan in imports and `go.mod`

### Added

- `ast.Type` enum and `ast.Vertex` `GetType() ast.Type` function for easily checking and storing vertex types
- Expose the internal lexer using `lexer.New(src []byte, config conf.Config) (Lexer, error)`

### Fixes

- Added some boundary checks in the formatter to prevent some panics
- A newline will now be added after a doc comment in the printer
- Some fixes for start and end column addition in previous release

## v0.9.0 2023-26-03

### Breaking Changes

- Dropped PHP 5 support to increase velocity/productivity

### Added

- Nodes and tokens now also contain their StartCol and EndCol
- 2 new interfaces for traversers (backwards compatible):
    - `traversers.CheckEntrance` interface: provides an `EnterNode(ast.Vertex) bool` that, when implemented, is called before entering a node, allowing short circuiting and more control over the traversal
    - `traversers.NotifyLeave` interface: provides an `LeaveNode(ast.Vertex)` that, when implemented, is called when a node has fully been traversed, again allowing for more control

### Fixes

- Fixed a bug where the parser would use token.DefaultBlockSize instead of position.DefaultBlockSize when instantiating the parser

### Internals

- Rely on code generations for several repeating/boilerplate tasks
- Refactors and deduplication after PHP 5 removal
- Add golangci-lint configuration and an editorconfig file

## v0.8.3 2022-09-09

### Changed

- reduce memory usage by allocating smaller position blocks

## `v0.8.2` 2022-26-06

### Added

- [#28](https://github.com/VKCOM/php-parser/pull/28): `php8.2`: added readonly classes support
- [#29](https://github.com/VKCOM/php-parser/pull/29): `php8.1`: added intersection types support

## `v0.8.1-rc.1` 2021-09-08

### Added

- [#6](https://github.com/VKCOM/php-parser/pull/6): `php8.1`: added `readonly` modifier
- [#8](https://github.com/VKCOM/php-parser/pull/8): `php8.1`: added `never` type
- [#10](https://github.com/VKCOM/php-parser/pull/10): `php8.1`: added new octal numbers syntax
- [#12](https://github.com/VKCOM/php-parser/pull/12): `php8.1`: added enums
- [#15](https://github.com/VKCOM/php-parser/pull/15): `php8.1`: added `final` modifier for constants in class
- [#18](https://github.com/VKCOM/php-parser/pull/18): `php8.1`: added first class callable syntax

### Changed

- [`4cd50d`](https://github.com/VKCOM/php-parser/commit/85b5d3ef36c9b12923404caf1c57497aa84cd50d): `cmd`: added file path output before errors

### Fixed

- [#22](https://github.com/VKCOM/php-parser/pull/22): fixed bug with `#` comments

## `v0.8.0-rc.2` 2021-30-07

### Added

- [#10](https://github.com/i582/php-parser/pull/10): `php8`: nullsafe operator (`?->`)
- [#13](https://github.com/i582/php-parser/pull/13): `php8`: named arguments 
- [#19](https://github.com/i582/php-parser/pull/19): `php8`: `match` expression 
- [#21](https://github.com/i582/php-parser/pull/21): `php8`: union types in type hints and `static` typehint 
- [#23](https://github.com/i582/php-parser/pull/23): `php8`: block `catch` without variable 
- [#25](https://github.com/i582/php-parser/pull/25): `php8`: trailing comma in parameter lists 
- [#27](https://github.com/i582/php-parser/pull/27): `php8`: `throw` can be used as an expression 
- [#32](https://github.com/i582/php-parser/pull/32): `php8`: declaring properties in the constructor 
- [#34](https://github.com/i582/php-parser/pull/34): `php8`: attributes 
- [#38](https://github.com/i582/php-parser/pull/38): `php8`: trailing comma in closure use list 

### Changed

- [#30](https://github.com/i582/php-parser/pull/30): `php8`: concatenation precedence 
- [#36](https://github.com/i582/php-parser/pull/36): `php8`: names in the namespace are treated as a single token 
- [#42](https://github.com/i582/php-parser/pull/42): `php8`: deferencable changes and arbitrary expressions in `new`/`instanceof` 

### Removed

- [#11](https://github.com/i582/php-parser/pull/11): `php8`: removed `(real)` cast 
- [#15](https://github.com/i582/php-parser/pull/15): `php8`: removed `(unset)` cast 
- [#17](https://github.com/i582/php-parser/pull/17): `php8`: removed `{}` access 

---

Versions prior to 0.8.0 were not included in this changelog.
