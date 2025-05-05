# Features

## Key aspects of the graceful shutdown implementation:

File: `cmd/server/main.go`

1. **Signal capture**: The code listens for system signals (SIGINT, SIGTERM, etc.) that indicate the server should shut down.
2. **Timeout context**: A timeout context is created to allow existing connections to complete within a grace period (30 seconds in this example).
3. **Resource cleanup**: After calling `server.Shutdown()`, the code can perform additional cleanup like closing database connections.
4. **Error handling**: The implementation handles errors during shutdown and forces termination if graceful shutdown exceeds the timeout.

### Signal capture:

Lines 106-107: sigCh := make(chan os.Signal, 1) and signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM) explicitly set up a channel to listen for SIGINT and SIGTERM signals.
Line 110: sig := <-sigCh blocks execution until one of these signals is received.

### Timeout context:

Line 114: ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second) creates a context with a 5-second timeout specifically for the shutdown process.
Line 119: This ctx is passed to server.Shutdown(ctx).

### Resource cleanup:

The call to server.Shutdown(ctx) happens on line 119.
Execution continues after this call completes (either successfully or with an error) on line 124 (end of the main function). Any necessary cleanup code (e.g., closing database connections, flushing logs) could be placed here, after the if/else block (lines 119-123).

### Error handling:

Lines 119-123: The code explicitly checks the error returned by server.Shutdown(ctx).
If an error occurs (e.g., the timeout is exceeded), it's logged using log.Printf("Server forced to shutdown: %v", err).
If shutdown completes within the timeout, a success message is logged.

## Environment-based server port configuration

- .env file support and a configuration package (`internal/config`) to manage application settings. The server ow reads the PORT variable from the environment (or .env file) to determine which port to listen on, defaulting to 8080.

