# Go Coding Official Standards & Best Practices

## 1. Always Format Code with `gofmt`
Use `gofmt` to automatically format your Go code. Consistent formatting improves readability and is the standard across all Go projects.

## 2. Use `goimports`
`goimports` formats your code and automatically adds or removes import statements, keeping imports clean and organized.

## 3. Use Meaningful Package Names
Choose short, lowercase, and descriptive package names that clearly represent their purpose. Avoid generic names like `utils` or `common`.

## 4. Follow Go Naming Conventions
Use PascalCase for exported identifiers and camelCase for unexported ones. Keep names simple, descriptive, and consistent.

## 5. Keep Packages Focused
A package should have a single responsibility. Avoid combining unrelated functionality into one package.

## 6. Group Imports Properly
Organize imports into three groups: standard library, third-party packages, and internal packages.

## 7. Always Handle Errors
Never ignore returned errors. Check every error and handle it appropriately to prevent unexpected behavior.

## 8. Wrap Errors with Context
When returning errors, add context using `fmt.Errorf` with `%w` to make debugging easier.

## 9. Prefer Early Returns
Return early when encountering errors instead of nesting multiple `if` statements. This keeps code cleaner and easier to read.

## 10. Keep Functions Small
Each function should perform one specific task. Small, focused functions are easier to test and maintain.

## 11. Limit Function Parameters
If a function requires many parameters, group them into a struct to improve readability and flexibility.

## 12. Use Short Receiver Names
Method receivers should use short names like `u`, `c`, or `s` instead of repeating the full type name.

## 13. Choose Pointer or Value Receivers Correctly
Use pointer receivers when modifying data or working with large structs. Use value receivers for immutable behavior.

## 14. Write Useful Comments
Document exported types, functions, and methods. Comments should explain why something exists rather than what the code obviously does.

## 15. Use Context for Cancellation
Pass `context.Context` to operations that may be canceled or have timeouts, such as HTTP requests or database queries.

## 16. Prevent Goroutine Leaks
Every goroutine should have a clear exit condition. Avoid creating goroutines that may run forever unintentionally.

## 17. Close Channels Correctly
The sender is responsible for closing channels. Receivers should never close channels they do not own.

## 18. Write Table-Driven Tests
Use table-driven tests to cover multiple scenarios with less code and better maintainability.

## 19. Run Tests Frequently
Run your test suite regularly to catch issues early and ensure code quality.

## 20. Detect Race Conditions
Use Go's race detector during development to identify concurrent access issues.

## 21. Preallocate Slices
When the size is known, preallocate slices using `make` to reduce memory allocations and improve performance.

## 22. Use `strings.Builder` for String Concatenation
For building large strings efficiently, prefer `strings.Builder` over repeated string concatenation.

## 23. Optimize Only After Profiling
Measure performance with profiling tools before making optimizations. Avoid premature optimization.

## 24. Run Static Analysis
Use tools like `go vet` to detect common mistakes and improve code quality.

## 25. Check for Security Vulnerabilities
Run `govulncheck` regularly to identify known vulnerabilities in your dependencies.

## 26. Use Go Modules
Manage dependencies with Go Modules and keep `go.mod` and `go.sum` committed to version control.

## 27. Keep Dependencies Updated
Regularly update third-party packages to receive bug fixes, performance improvements, and security patches.

## 28. Use Structured Logging
Prefer structured logging libraries such as `slog`, `zap`, or `zerolog` for better log analysis and monitoring.

## 29. Avoid Unnecessary Panic
Use `panic` only for unrecoverable situations. Return errors whenever recovery is possible.

## 30. Define Interfaces Where They Are Used
Create interfaces in the consuming package rather than the implementation package to reduce coupling.

## 31. Prefer Composition Over Inheritance
Build reusable components by composing smaller types instead of creating complex inheritance hierarchies.

## 32. Keep Project Structure Organized
Separate commands, internal packages, configuration, documentation, and reusable components into a clear directory structure.

## 33. Write Readable Code
Prioritize readability over clever or overly complex implementations. Code is read more often than it is written.

## 34. Keep Code Explicit
Prefer clear and straightforward logic instead of hiding behavior behind unnecessary abstractions.

## 35. Follow the Go Philosophy
Write simple, maintainable, idiomatic Go code that is easy for other Go developers to understand and contribute to.