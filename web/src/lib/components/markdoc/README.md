<!--- START SELFDOC --->
## SelfDoc
_Auto-generated code documentation to make the repository easier to navigate and contribute to._

_Last Updated: 2023-05-15_

The `markdoc` directory is for exporting Svelte components and configuration objects used in rendering markdown as HTML using the `@markdoc/markdoc` library. It includes components for rendering styled boxes, email labels, and tags, as well as a configuration file for generating a configuration object based on an abstract syntax tree (AST).

### Files
#### Callout.svelte
This file exports a Svelte component called `Callout` that renders a styled box with a title and content. The component accepts a `title` prop and a slot for content. The box includes an icon and is styled differently based on whether the page is in light or dark mode.

#### EmailLabel.svelte
This file exports an EmailLabel Svelte component that takes in a `title` and `color` prop. The `color` prop defaults to `'blue'` and can be set to `'green'` or `'red'`. The component renders a span element with the `title` as its content and a background color determined by the `color` prop.

#### Markdoc.svelte
This file exports a Svelte component called `Markdoc` which takes in a `source` prop containing raw markdown and renders it as HTML using the `@markdoc/markdoc` library. It also imports a `config` object, `Tag` component, `Callout` component, and `EmailLabel` component to be used in the markdown rendering.

#### Tag.svelte
This file exports a Svelte component called `Tag` that renders a tree of nodes. It checks if the content is a tag or an array of tags and recursively renders them. It also checks if the tag is a Svelte component and renders it accordingly.

#### config.ts
This file exports a function that returns a configuration object for the `Markdoc` component. The configuration object includes information about the nodes and tags used in the component, as well as any frontmatter variables. The `config` function takes an abstract syntax tree (AST) as an argument and uses it to generate the configuration object.

#### nodes.ts
This file exports a `heading` object that defines the structure of a heading node in the Markdoc syntax. It also includes a function `headingID` that generates a unique ID for the heading based on its content. The `transform` function of the `heading` object uses this ID to create an HTML heading element with the appropriate level and attributes.

<!--- END SELFDOC --->