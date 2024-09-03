# Changelog

## [1.0.0] - 2024-09-03
### Added
- Support for pre- and post-log hooks in the Logger struct.
    - Added methods `WithPreHook` and `WithPostHook` to add hooks to the logger.
    - Modified `LogE` method to execute hooks before and after logging.
- Context support in the Logger struct.
    - Added method `WithContext` to set the context for the logger.
    - Modified `getData` function to include context information in the log data.

## [0.1.0] - 2024-09-03
### Added
- Initial release of Loggo
- Configurable log levels (Debug, Info, Warn, Error, Fatal)
- Customizable output destinations (e.g., `os.Stdout`, `os.Stderr`, files)
- Flexible message templates
- Custom time providers for log timestamps
- Thread-safe logging with `sync.Mutex`
