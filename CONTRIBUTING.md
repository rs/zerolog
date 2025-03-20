# Contributing to Zerolog

Thank you for your interest in contributing to **Zerolog**!

Zerolog is a **feature-complete**, high-performance logging library designed to be **lean** and **non-bloated**. The focus of ongoing development is on **bug fixes**, **performance improvements**, and **modernization efforts** (such as keeping up with Go best practices and compatibility with newer Go versions).

## What We're Looking For

We welcome contributions in the following areas:

- **Bug Fixes**: If you find an issue or unexpected behavior, please open an issue and/or submit a fix.
- **Performance Optimizations**: Improvements that reduce memory usage, allocation count, or CPU cycles without introducing complexity are appreciated.
- **Modernization**: Compatibility updates for newer Go versions or idiomatic improvements that do not increase library size or complexity.
- **Documentation Enhancements**: Corrections, clarifications, and improvements to documentation or code comments.

## What We're *Not* Looking For

Zerolog is intended to remain **minimalistic and efficient**. Therefore, we are **not accepting**:

- New features that add optional behaviors or extend API surface area.
- Built-in support for frameworks or external systems (e.g., bindings, integrations).
- General-purpose abstractions or configuration helpers.

If you're unsure whether a change aligns with the project's philosophy, feel free to open an issue for discussion before submitting a PR.

## Contributing Guidelines

1. **Fork the repository**
2. **Create a branch** for your fix or improvement
3. **Write tests** to cover your changes
4. Ensure `go test ./...` passes
5. Run `go fmt` and `go vet` to ensure code consistency
6. **Submit a pull request** with a clear explanation of the motivation and impact

## Code Style

- Keep the code simple, efficient, and idiomatic.
- Avoid introducing new dependencies.
- Preserve backwards compatibility unless explicitly discussed.

---

We appreciate your effort in helping us keep Zerolog fast, minimal, and reliable!
