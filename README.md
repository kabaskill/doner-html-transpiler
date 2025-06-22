# D.Ö.N.E.R - Deutsche Öffnung zur Normalisierten ERkennung

*A German HTML transpiler born from curiosity and a love for döner kebab* 🥙

## The Story Behind the Project

I'd been hearing about Go everywhere but needed a concrete project to actually learn it. Living in Germany, you can't just unnotice how big localization is around here. That made me think "How about I localize HTML to German?" At this point I need to admit that the idea of "German HTML" seemed 'cursed' enough to be actually funny. However, I didn't knew the rabbit hole that I just threw myself in.

## Why This Project Exists

Because why not? 

**Learning Go**: I've always been a 'learn by doing it' kind of guy. So, while I was learning things like how GO's syntax works, what is it good for and so on, I was also thinking about project ideas. Turns out, everybody who uses GO, loves it. Working with pointers and weird syntax bits here and there was a struggle for me as a Typescript developer. At least there is no manual garbage collection. As far as I'm aware of course.

**Compiler Fundamentals**: You don't really use things like AST, Tokenization and DSL when you are developing web apps. This isn't the case for everybody of course, but for me these concepts were just theory. So that's another reason. I wanted to explore those concepts in detail.

**Real-World Challenges**: Building a public API forced me to think about what sort of security issues I might have. API throttling and rate limiting were only the start. I explored how to sanitize the user input while not blocking them and still providing what they asked for. 

## What It Does

D.Ö.N.E.R converts German HTML tags to their English equivalents. Currently, 50 tags and 20 attributes are supported. For the full list check the API Reference on this document or the browser version of the app.

```html
<!-- German HTML -->
<döner>
  <kopf>
    <titel>Meine deutsche Webseite</titel>
  </kopf>
  <körper>
    <hauptüberschrift>Willkommen!</hauptüberschrift>
    <absatz>Das ist ein deutscher Absatz.</absatz>
  </körper>
</döner>

<!-- Becomes Standard HTML -->
<html>
  <head>
    <title>Meine deutsche Webseite</title>
  </head>
  <body>
    <h1>Willkommen!</h1>
    <p>Das ist ein deutscher Absatz.</p>
  </body>
</html>
```

## Technical Architecture

### Backend (Go)

GO is responsible for all the heavy lifting for the app.

- **Lexical Analysis**: Tokenizes German HTML into meaningful components
- **Parser**: Builds an Abstract Syntax Tree (AST) from tokens
- **Transpiler**: Transforms the AST into standard HTML
- **Web Server**: RESTful API with CORS, rate limiting, and security headers
- **Security**: Input validation, output sanitization, and XSS protection

### Frontend (React + TypeScript)

I haven't spent too much time on the frontend side of things, since the interesting part was the transpiler. However, that means there is a lot of room to improve the general UX. Here's the tech-stack:
 - **React 19, Tailwind and Vite**: For the general framework for the app all bundled with [blitz-react](https://www.npmjs.com/package/blitz-react?activeTab=readme)
 - **Slate Editor**: It's a great text input component for web applications. Check it out [here](https://github.com/ianstormtaylor/slate).

### Project Structure

```
doner/
├── backend/             # Go backend
│   ├── main.go          # HTTP server & CLI entry point
│   ├── transpiler.go    # Core transpilation engine
│   ├── parser.go        # HTML parser with error handling
│   ├── lexer.go         # Token-based lexical analyzer
│   ├── dictionary.go    # German→English mappings
│   └── ast.go           # Abstract syntax tree definitions
├── frontend/            # React frontend
│   ├── src/
│   │   ├── App.tsx      # Main application
│   │   ├── components/  # Modular UI components
│   │   └── utils/       # helper functions
│   └── ...other React files
└── README.md           # You are here!
```

## Getting Started

### Quick Start

You need GO and Node installed if you want to run it on local.

1. **Clone and run the backend:**
   ```bash
   cd backend
   go run main.go
   # Server starts on http://localhost:8080
   ```

2. **Start the frontend:**
   ```bash
   cd frontend
   npm install && npm run dev
   # Frontend starts on http://localhost:5173
   ```

3. **Try it out!** Open your browser and start writing German HTML.

### CLI Usage

You can also use D.Ö.N.E.R as a command-line tool:

```bash
cd backend
echo '<döner><kopf><titel>Test</titel></kopf></döner>' > test.doner
go run . test.doner
# Outputs: <html><head><title>Test</title></head></html>
```

## API Reference

### `POST /transpile`
Converts German HTML to standard HTML.

**Request:**
```json
{
  "content": "<döner><kopf><titel>Meine Seite</titel></kopf></döner>"
}
```

**Response:**
```json
{
  "result": "<html><head><title>Meine Seite</title></head></html>",
  "warnings": ["No security issues detected"]
}
```

### `GET /dictionary`
Returns all German→English tag mappings.

**Response:**
```json
{
  "döner": "html",
  "kopf": "head",
  "körper": "body",
  "titel": "title"
}
```

### `GET /health`
Health check endpoint.

**Response:**
```json
{
  "status": "ok",
  "version": "1.0.0"
}
```

## Security

There are some security features in place to keep the app and my sanity intact. However, the web app still gives you the text output. It just won't render the HTML generated if security checks fail. 

Some of those are:
 - 100 KB size limit for requests
 - 1000 Tokens per request
 - HTML tags like script, iframe and event listeners are not supported for rendering.


## What I Learned Building This

**Go vs TypeScript**: Go felt awkward at first. But the simplicity grew on me – you only import external libraries, there is no weird scope issues you might experience on TS, but you just worry about writing the most efficient code.

**Lexer Reality Check**: Building a tokenizer sounds straight forward until you hit edge cases. Right now, the lexer breaks on `!` characters because I haven't handled all possible HTML content properly. DOCTYPE declarations, comments, and CDATA sections all need special treatment that I initially glossed over.

**Parser Complexity**: Converting tokens to an AST while maintaining proper HTML structure and handling malformed input turned out to be trickier than expected. Self-closing tags, nested elements and attribute parsing each required their own logic.

**Output Formatting**: Getting the correnct indentation in the transpiled HTML was one of the challenges. Keeping track of where you are in the AST is definitely more challenging than it sounds.


## Future Ideas

Some features I'd love to add if I found the time:

- **More Languages**: Why stop at German? Why not do a Turkish HTML while we are at it?
- **CSS Extension**: The dictionary currently has 70 entries (tags and attributes). Extending it with CSS properties sounds cool but it is at least 10 times the work.
- **VS Code Extension**: It would be a casual, funny extension to f*ck around
- **Syntax Highlighting**: Make the editor even more beautiful with additional features like highlighting tags and parsing errors.
- **Community Dictionary**: It would be a great addition to have a UI for users who want to add new tags in the dictionary.

## Contributing

Found a bug? Want to add more German tags? I'd love your help!

1. Fork the repository
2. Create a feature branch (`git checkout -b feature/amazing-feature`)
3. Make your changes
4. Add tests if applicable
5. Submit a pull request

Check out the current dictionary in `backend/dictionary.go` – it's easy to add new translations!

## License

This project is open source and available under the [MIT License](LICENSE). Feel free to use it, modify it, or learn from it.

---

*Built with ❤️ in Germany, inspired by döner kebab and the beauty of language diversity.*

**Live Demo**: [Try D.Ö.N.E.R online](https://your-deployment-url.com) (Coming soon!)
