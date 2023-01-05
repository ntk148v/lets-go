# Context

## Article 1

<https://medium.com/@cep21/how-to-correctly-use-context-context-in-go-1-7-8f2c0fafdf39>

### Integrate Context into API

- Context is intended to be request scope.
- There are currently two ways to integrate Context objects into API:
  - The first parameter of a function call.
  - Optional config on a request structure.

### Context should flow through your program

- **Don't store it** somewhere like in a struct.
- Context should be an interface that is passed from function to function down your callback. Ideally, a Context object is created with each request and expires when the request is over.
- The one exception to not storing a context is when you need to put it in a struct that is used purely as a message that is passed across a channel.

### All blocking/long operations should be cancelable

### Context.Value and request-scoped values (a warning)

- A request-scoped value is one **derived** from data in the **incoming request** and **goes away** when the request is over.
- Obvious request scoped data could be _who_ is making the request (user ID), _how_ they are making it (internal or external), from _where_ they are making it (user IP address), and _how_ important this request should be.

### Context.Value obscures your program's flow

- Obscures expected input and output of a function or library.

### Context.Value should inform, not control

- **Inform, not control**.
- **The content of context.Value is for maintainers not users**. It should dnever be required input for documented or expected results.

### Try not to use context.Value

## Article 2

<http://p.agnihotry.com/post/understanding_the_context_package_in_golang/>

**Best practices**:

- context.Background should be used only at the highest level, as the root of all derived contexts.
- context.TODO should be used where not sure what to use or if the current function will be updated to use context in future.
- context cancelations are advisory, the functions may take time to clean up and exit.
- context.Value should be used very rarely, it should never be used to pass in optional parameters. This makes the API implicit and can introduce bugs. Instead, such values should be passed in as arguments.
- Don't store context in struct, pass them explicitly in function, perferably, as the 1st argument.
- Never pass nil context, instead, use a TODO if you are not sure what to use.
- The context struct does not have a cancel method because only the function that derives the context should cancel it.
