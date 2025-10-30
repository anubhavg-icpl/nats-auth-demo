# Contributing to NATS Auth Demo

Thank you for your interest in contributing! This project aims to provide clear, educational examples of NATS authorization and multi-tenancy features.

## How to Contribute

### Reporting Issues
- Use GitHub Issues to report bugs or suggest features
- Provide clear descriptions and reproduction steps
- Include relevant configuration files and error messages

### Submitting Changes
1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Test your changes thoroughly
5. Commit with clear messages (`git commit -m 'Add amazing feature'`)
6. Push to your branch (`git push origin feature/amazing-feature`)
7. Open a Pull Request

### Code Style
- Follow Go best practices and conventions
- Use `go fmt` to format code
- Add comments for complex logic
- Keep examples simple and educational

### Adding New Examples
If adding new authorization examples:
1. Create a new configuration file in `config/`
2. Add corresponding Go example in `examples/`
3. Update the main menu in `cmd/main.go`
4. Document the example in README.md
5. Test thoroughly with both Docker and Podman

### Testing
- Ensure all examples run successfully
- Test with both `docker-compose` and `podman-compose`
- Verify authorization rules work as expected
- Check that denied operations are properly rejected

## Code of Conduct
- Be respectful and constructive
- Focus on education and clarity
- Help others learn NATS

## Questions?
Feel free to open an issue for questions or discussions.

Thank you for helping make this project better!
